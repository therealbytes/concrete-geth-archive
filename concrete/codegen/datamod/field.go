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

package datamod

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func EncodeAddress(_ int, address common.Address) []byte {
	return address.Bytes()
}

func DecodeAddress(_ int, data []byte) common.Address {
	return common.BytesToAddress(data)
}

func EncodeBool(_ int, b bool) []byte {
	if b {
		return []byte{1}
	}
	return []byte{0}
}

// TODO: bool decoding consistency
func DecodeBool(_ int, data []byte) bool {
	return data[0] != 0
}

func EncodeBytes(size int, b []byte) []byte {
	return common.LeftPadBytes(b, size)
}

func DecodeBytes(size int, data []byte) []byte {
	return common.LeftPadBytes(data, size)
}

// TODO: this
func EncodeInt(size int, i *big.Int) []byte {
	return nil
}

func DecodeInt(size int, data []byte) *big.Int {
	return nil
}

func EncodeUint(size int, i *big.Int) []byte {
	return nil
}

func DecodeUint(size int, data []byte) *big.Int {
	return nil
}

type FieldType struct {
	Name       string
	Size       int
	GoType     string
	EncodeFunc string
	DecodeFunc string
}

var NameToFieldType map[string]FieldType

func init() {
	NameToFieldType = make(map[string]FieldType)
	NameToFieldType["address"] = FieldType{
		Name:       "address",
		Size:       20,
		GoType:     "common.Address",
		EncodeFunc: "EncodeAddress",
		DecodeFunc: "DecodeAddress",
	}
	NameToFieldType["bool"] = FieldType{
		Name:       "bool",
		Size:       1,
		GoType:     "bool",
		EncodeFunc: "EncodeBool",
		DecodeFunc: "DecodeBool",
	}
	for ii := 1; ii <= 32; ii++ {
		name := "bytes" + fmt.Sprint(ii)
		NameToFieldType[name] = FieldType{
			Name:       name,
			Size:       ii,
			GoType:     "[]byte",
			EncodeFunc: "EncodeBytes",
			DecodeFunc: "DecodeBytes",
		}
	}
	for ii := 1; ii <= 32; ii++ {
		name := "int" + fmt.Sprint(ii*8)
		NameToFieldType[name] = FieldType{
			Name:       name,
			Size:       ii,
			GoType:     "*big.Int",
			EncodeFunc: "EncodeInt",
			DecodeFunc: "DecodeInt",
		}
	}
	for ii := 1; ii <= 32; ii++ {
		name := "uint" + fmt.Sprint(ii*8)
		NameToFieldType[name] = FieldType{
			Name:       name,
			Size:       ii,
			GoType:     "*big.Int",
			EncodeFunc: "EncodeUint",
			DecodeFunc: "DecodeUint",
		}
	}

	NameToFieldType["int"] = NameToFieldType["int256"]
	NameToFieldType["uint"] = NameToFieldType["uint256"]
}
