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
)

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
