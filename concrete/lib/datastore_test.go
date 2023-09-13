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

//go:build !tinygo

// This file will ignored when building with tinygo to prevent compatibility
// issues.

package lib

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/concrete/api"
	"github.com/ethereum/go-ethereum/concrete/mock"
	"github.com/stretchr/testify/require"
)

func TestEnvKeyValueStore(t *testing.T) {
	var (
		r        = require.New(t)
		address  = common.HexToAddress("0xc0ffee0001")
		config   = api.EnvConfig{Trusted: true, Ephemeral: true}
		meterGas = false
		gas      = uint64(0)
	)
	tests := []struct {
		name string
		kv   KeyValueStore
	}{
		{
			name: "Persistent",
			kv:   newEnvPersistentKeyValueStore(mock.NewMockEnvironment(address, config, meterGas, gas)),
		},
		{
			name: "Ephemeral",
			kv:   newEnvEphemeralKeyValueStore(mock.NewMockEnvironment(address, config, meterGas, gas)),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				kv  = test.kv
				key = common.Hash{0x01}
			)
			value := kv.Get(key)
			r.Equal(common.Hash{}, value)
			kv.Set(key, common.Hash{0x02})
			value = kv.Get(key)
			r.Equal(common.Hash{0x02}, value)
		})
	}
}

func TestDatastore(t *testing.T) {
	var (
		r        = require.New(t)
		address  = common.HexToAddress("0xc0ffee0001")
		config   = api.EnvConfig{}
		meterGas = false
		gas      = uint64(0)
	)
	env := mock.NewMockEnvironment(address, config, meterGas, gas)
	ds := NewPersistentDatastore(env)
	key := []byte{0x01}
	slot := ds.Get(key)
	r.NotNil(slot)
	r.NotNil(slot.Datastore())
	r.Equal(slot.Slot(), common.BytesToHash(key))
}

func TestDatastoreSlot(t *testing.T) {
	var (
		r        = require.New(t)
		address  = common.HexToAddress("0xc0ffee0001")
		config   = api.EnvConfig{}
		meterGas = false
		gas      = uint64(0)
	)
	env := mock.NewMockEnvironment(address, config, meterGas, gas)
	ds := NewPersistentDatastore(env)
	key := []byte{0x01}
	slot := ds.Get(key)

	r.Equal(common.Hash{}, slot.Bytes32())
	r.Equal(false, slot.Bool())
	r.Equal(common.Address{}, slot.Address())
	r.Equal(int64(0), slot.BigUint().Int64())
	r.Equal(int64(0), slot.BigInt().Int64())
	r.Equal(uint64(0), slot.Uint64())
	r.Equal(int64(0), slot.Int64())
	r.Equal([]byte{}, slot.Bytes())

	slot.SetBytes32(common.Hash{0x01})
	r.Equal(common.Hash{0x01}, slot.Bytes32())

	slot.SetBool(true)
	r.Equal(true, slot.Bool())

	slot.SetAddress(address)
	r.Equal(address, slot.Address())

	slot.SetBigUint(big.NewInt(1))
	r.Equal(int64(1), slot.BigUint().Int64())

	slot.SetBigInt(big.NewInt(-1))
	r.Equal(int64(-1), slot.BigInt().Int64())

	slot.SetUint64(1)
	r.Equal(uint64(1), slot.Uint64())

	slot.SetInt64(-1)
	r.Equal(int64(-1), slot.Int64())

	slot.SetBytes([]byte{0x01, 0x02, 0x03})
	r.Equal([]byte{0x01, 0x02, 0x03}, slot.Bytes())
}
