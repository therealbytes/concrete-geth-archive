// Copyright 2023 The concrete-geth Authors
//
// The concrete-geth library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The concrete library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the concrete library. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/concrete/codegen/solgen"
	"github.com/spf13/cobra"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func isDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if info.IsDir() {
		return true, nil
	} else {
		return false, nil
	}
}

func fileName(path string) string {
	filenameWithExt := filepath.Base(path)
	filename := strings.TrimSuffix(filenameWithExt, filepath.Ext(filenameWithExt))
	return filename
}

func main() {
	var rootCmd = &cobra.Command{Use: "concrete"}

	var cmdCodegen = &cobra.Command{
		Use:   "solgen",
		Short: "generate a solidity precompile caller library from an ABI file",
		Run:   runCodeGen,
	}

	cmdCodegen.Flags().String("abi", "", "path to the ABI file")
	cmdCodegen.Flags().String("out", "", "path to the output file")
	cmdCodegen.Flags().String("solidity", "", "path to the solidity file")
	cmdCodegen.Flags().StringP("name", "n", "", "name for the generated library")
	cmdCodegen.Flags().StringP("address", "a", "", "precompile address")

	rootCmd.AddCommand(cmdCodegen)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runCodeGen(cmd *cobra.Command, args []string) {
	abiPath, err := cmd.Flags().GetString("abi")
	checkErr(err)
	outPath, err := cmd.Flags().GetString("out")
	checkErr(err)
	solPath, err := cmd.Flags().GetString("solidity")
	checkErr(err)
	name, err := cmd.Flags().GetString("name")
	checkErr(err)
	address, err := cmd.Flags().GetString("address")
	checkErr(err)

	if abiPath == "" {
		fmt.Println("Missing ABI file path (--abi)")
		os.Exit(1)
	}
	if outPath == "" {
		fmt.Println("Missing output file path (--out))")
		os.Exit(1)
	}

	if address == "" {
		fmt.Println("Missing precompile address (--address)")
	} else if !common.IsHexAddress(common.HexToAddress(address).Hex()) {
		fmt.Println("Invalid address:", address)
		os.Exit(1)
	} else {
		address = common.HexToAddress(address).Hex()
	}

	abiIsDir, err := isDir(abiPath)
	checkErr(err)
	outIsDir, err := isDir(outPath)
	checkErr(err)

	if abiIsDir {
		fmt.Println("ABI path must be a file")
		os.Exit(1)
	}

	if name == "" {
		name = fileName(abiPath) + "Precompile"
	}

	if outIsDir {
		outPath = filepath.Join(outPath, name+".sol")
	}

	config := solgen.Config{
		Name:    name,
		Address: common.HexToAddress(address),
		ABI:     abiPath,
		Out:     outPath,
		Sol:     solPath,
	}

	fmt.Printf(`Generating solidity library
Name     : %s
Address  : %s
ABI      : %s
Output   : %s
Solidity : %s
`, name, address, abiPath, outPath, solPath)

	err = solgen.GenerateSolidityLib(config)
	checkErr(err)

	fmt.Printf("Library generated successfully.\nLibrary written to: %s\n", outPath)
}
