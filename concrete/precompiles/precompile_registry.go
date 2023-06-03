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

package precompiles

import (
	_ "embed"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/concrete/api"
	"github.com/ethereum/go-ethereum/concrete/lib"
)

//go:embed sol/abi/PrecompileRegistry.abi
var precompileRegistryABI string

var PrecompileRegistryMetadata = PrecompileMetadata{
	Name:        "PrecompileRegistry",
	Version:     "0.1.0",
	Author:      "The concrete-geth Authors",
	Description: "A registry of precompiles indexed by address and name.",
	Source:      "https://github.com/therealbytes/concrete-geth/tree/concrete/concrete/precompiles/precompile_registry.go",
}

var (
	GetFrameworkGas            = uint64(10)
	GetPrecompileGas           = uint64(10)
	GetPrecompileByNameGas     = uint64(10)
	GetPrecompiledAddressesGas = uint64(10)
	GetPrecompilesGas          = uint64(10)
)

func init() {
	abiReader := strings.NewReader(precompileRegistryABI)
	ABI, err := abi.JSON(abiReader)
	if err != nil {
		panic(err)
	}
	precompileRegistry := lib.NewPrecompileWithABI(ABI, map[string]lib.MethodPrecompile{
		"getFramework":            &getFramework{},
		"getPrecompile":           &getPrecompile{},
		"getPrecompileByName":     &getPrecompileByName{},
		"getPrecompiledAddresses": &getPrecompiledAddresses{},
		"getPrecompiles":          &getPrecompiles{},
	})
	AddPrecompile(api.PrecompileRegistryAddress, precompileRegistry, PrecompileRegistryMetadata)
}

type getFramework struct {
	lib.BlankMethodPrecompile
}

func (p *getFramework) RequiredGas(input []byte) uint64 {
	return GetFrameworkGas
}

func (p *getFramework) Run(concrete api.API, input []byte) ([]byte, error) {
	return p.CallRunWithArgs(func(concrete api.API, args []interface{}) ([]interface{}, error) {
		metadata := struct {
			Name    string `json:"name"`
			Version string `json:"version"`
			Source  string `json:"source"`
		}{
			"Concrete",
			"0.1.0",
			"https://github.com/therealbytes/concrete-geth",
		}
		return []interface{}{metadata}, nil
	}, concrete, input)
}

type getPrecompile struct {
	lib.BlankMethodPrecompile
}

func (p *getPrecompile) RequiredGas(input []byte) uint64 {
	return GetPrecompileGas
}

func (p *getPrecompile) Run(concrete api.API, input []byte) ([]byte, error) {
	return p.CallRunWithArgs(func(concrete api.API, args []interface{}) ([]interface{}, error) {
		address := common.Address(args[0].(common.Address))
		if metadata, ok := metadataByAddress[address]; ok {
			return []interface{}{&metadata}, nil
		}
		return []interface{}{PrecompileMetadata{}}, nil
	}, concrete, input)
}

type getPrecompileByName struct {
	lib.BlankMethodPrecompile
}

func (p *getPrecompileByName) RequiredGas(input []byte) uint64 {
	return GetPrecompileByNameGas
}

func (p *getPrecompileByName) Run(concrete api.API, input []byte) ([]byte, error) {
	return p.CallRunWithArgs(func(concrete api.API, args []interface{}) ([]interface{}, error) {
		name := args[0].(string)
		if metadata, ok := metadataByName[name]; ok {
			return []interface{}{metadata.Addr}, nil
		}
		return []interface{}{common.Address{}}, nil
	}, concrete, input)
}

type getPrecompiledAddresses struct {
	lib.BlankMethodPrecompile
}

func (p *getPrecompiledAddresses) RequiredGas(input []byte) uint64 {
	return GetPrecompiledAddressesGas
}

func (p *getPrecompiledAddresses) Run(concrete api.API, input []byte) ([]byte, error) {
	return p.CallRunWithArgs(func(concrete api.API, args []interface{}) ([]interface{}, error) {
		return []interface{}{precompiledAddresses}, nil
	}, concrete, input)
}

type getPrecompiles struct {
	lib.BlankMethodPrecompile
}

func (p *getPrecompiles) RequiredGas(input []byte) uint64 {
	return GetPrecompilesGas
}

func (p *getPrecompiles) Run(concrete api.API, input []byte) ([]byte, error) {
	return p.CallRunWithArgs(func(concrete api.API, args []interface{}) ([]interface{}, error) {
		return []interface{}{precompileMetadata}, nil
	}, concrete, input)
}