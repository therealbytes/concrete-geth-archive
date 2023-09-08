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

package lib

import (
	"bytes"
)

type DatastoreStruct struct {
	store   DatastoreSlot
	arr     SlotArray
	offsets []int
	sizes   []int
	cache   [][]byte
}

func NewDatastoreStruct(store DatastoreSlot, sizes []int) *DatastoreStruct {
	var (
		offset  = 0
		offsets = make([]int, len(sizes))
		size    = sizes[0]
	)
	for ii := 1; ii < len(sizes); ii++ {
		size = sizes[ii]
		if offset/32 < (offset+size)/32 {
			offset = (offset/32 + 1) * 32
		}
		offset += size
		offsets[ii] = offset
	}
	nSlots := (offset + size + 31) / 32

	return &DatastoreStruct{
		store:   store,
		arr:     store.SlotArray([]int{nSlots}),
		offsets: offsets,
		sizes:   sizes,
		cache:   make([][]byte, len(sizes)),
	}
}

func (s *DatastoreStruct) GetField(index int) []byte {
	fieldSize := s.sizes[index]
	if fieldSize == 0 {
		return nil
	}

	if len(s.cache[index]) == 0 {
		s.cache[index] = make([]byte, fieldSize)
	} else {
		data := make([]byte, fieldSize)
		copy(data, s.cache[index])
		return data
	}

	absOffset := s.offsets[index]
	slotIndex, slotOffset := absOffset/32, absOffset%32
	slotData := s.arr.Value(slotIndex).Bytes32()
	return slotData[slotOffset : slotOffset+fieldSize]
}

func (s *DatastoreStruct) SetField(index int, data []byte) {
	fieldSize := s.sizes[index]
	if fieldSize == 0 {
		return
	}

	if len(s.cache[index]) == 0 {
		s.cache[index] = make([]byte, fieldSize)
	} else if bytes.Equal(s.cache[index], data) {
		return
	}

	absOffset := s.offsets[index]
	slotIndex, slotOffset := absOffset/32, absOffset%32
	slotRef := s.arr.Value(slotIndex)
	slotData := slotRef.Bytes32()

	copy(slotData[slotOffset:slotOffset+fieldSize], data)
	copy(s.cache[index], data)
	slotRef.SetBytes32(slotData)
}

func (s *DatastoreStruct) GetField_bytes(index int) []byte {
	absOffset := s.offsets[index]
	slotIndex := absOffset / 32
	slotData := s.arr.Value(slotIndex).Bytes()
	return slotData
}

func (s *DatastoreStruct) SetField_bytes(index int, data []byte) {
	if bytes.Equal(s.cache[index], data) {
		return
	}

	absOffset := s.offsets[index]
	slotIndex := absOffset / 32
	slotRef := s.arr.Value(slotIndex)

	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)
	s.cache[index] = dataCopy
	slotRef.SetBytes(data)
}
