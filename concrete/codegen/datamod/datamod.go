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
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"unicode"

	"github.com/ethereum/go-ethereum/concrete/lib"
)

//go:embed struct.tpl
var structTpl string

type StorageStruct struct {
	storage lib.SlotArray
	offsets []int
	sizes   []int
}

func NewStorageStruct(slot lib.StorageSlot, sizes []int) *StorageStruct {
	offsets := make([]int, len(sizes))
	offset := 0
	nSlots := 0
	for i := 1; i < len(sizes); i++ {
		size := sizes[i]
		if offset/32 < (offset+size)/32 {
			offset = (offset/32 + 1) * 32
			nSlots++
		}
		offset += size
		offsets[i] = offset
	}
	return &StorageStruct{
		storage: slot.SlotArray([]int{nSlots}),
		offsets: offsets,
		sizes:   sizes,
	}
}

func (s *StorageStruct) GetField(index int) []byte {
	absOffset := s.offsets[index]
	slotIdx := absOffset / 32
	slotOffset := absOffset % 32

	slotValue := s.storage.Value(slotIdx).Bytes32()
	size := s.sizes[index]
	return slotValue[slotOffset : slotOffset+size]
}

func (s *StorageStruct) SetField(index int, data []byte) {
	absOffset := s.offsets[index]
	slotIdx := absOffset / 32
	slotOffset := absOffset % 32

	slot := s.storage.Value(slotIdx)
	slotValue := slot.Bytes32()
	size := s.sizes[index]
	copy(slotValue[slotOffset:slotOffset+size], data)
	slot.SetBytes32(slotValue)
}

type FieldSchema struct {
	Name  string
	Title string
	Index int
	Type  FieldType
}

type MappingSchema struct {
	Name   string
	Keys   []FieldSchema
	Values []FieldSchema
}

type ModelSchema []MappingSchema

type MappingUnmarshal struct {
	KeySchema map[string]string `json:"keySchema"`
	Schema    map[string]string `json:"schema"`
}

type ModelUnmarshal map[string]MappingUnmarshal

type Config struct {
	JSON    string
	Out     string
	Package string
}

func lowerFirstLetter(str string) string {
	if len(str) == 0 {
		return ""
	}
	runes := []rune(str)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func upperFirstLetter(str string) string {
	if len(str) == 0 {
		return ""
	}
	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func isValidName(name string) bool {
	if len(name) == 0 {
		return false
	}
	re := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	return re.MatchString(name) && len(strings.TrimSpace(name)) == len(name)
}

func newFieldSchema(name string, index int, typeStr string) (FieldSchema, error) {
	if !isValidName(name) {
		return FieldSchema{}, fmt.Errorf("invalid name: %s", name)
	}
	fieldType, ok := NameToFieldType[typeStr]
	if !ok {
		return FieldSchema{}, fmt.Errorf("invalid type: %s", typeStr)
	}
	return FieldSchema{
		Name:  lowerFirstLetter(name),
		Title: upperFirstLetter(name),
		Index: index,
		Type:  fieldType,
	}, nil
}

func GenerateDataModel(config Config) error {
	if !isValidName(config.Package) {
		return fmt.Errorf("invalid package name: %s", config.Package)
	}

	jsonContent, err := os.ReadFile(config.JSON)
	if err != nil {
		return err
	}

	var unmarshaledModel ModelUnmarshal
	err = json.Unmarshal(jsonContent, &unmarshaledModel)
	if err != nil {
		return err
	}

	var model ModelSchema
	for name, mapping := range unmarshaledModel {
		if !isValidName(name) {
			return fmt.Errorf("invalid name: %s", name)
		}
		newMapping := MappingSchema{Name: upperFirstLetter(name)}
		for keyName, keyType := range mapping.KeySchema {
			fieldSchema, err := newFieldSchema(keyName, len(newMapping.Keys), keyType)
			if err != nil {
				return err
			}
			newMapping.Keys = append(newMapping.Keys, fieldSchema)
		}
		for valueName, valueType := range mapping.Schema {
			fieldSchema, err := newFieldSchema(valueName, len(newMapping.Values), valueType)
			if err != nil {
				return err
			}
			newMapping.Values = append(newMapping.Values, fieldSchema)
		}
		model = append(model, newMapping)
	}

	funcMap := template.FuncMap{
		"sub": func(a, b int) int { return a - b },
	}

	tmpl, err := template.New("struct").Funcs(funcMap).Parse(structTpl)
	if err != nil {
		return err
	}

	for _, mapping := range model {
		data := map[string]interface{}{
			"Package": config.Package,
			"Schema":  mapping,
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return err
		}
		outPath := filepath.Join(config.Out, lowerFirstLetter(mapping.Name)+".go")
		err := os.WriteFile(outPath, buf.Bytes(), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
