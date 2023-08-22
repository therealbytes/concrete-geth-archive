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

package api

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestGas(t *testing.T) {
	var (
		r       = require.New(t)
		address = common.HexToAddress("0xc0ffee0001")
		config  = EnvConfig{
			Static:    false,
			Ephemeral: false,
			Preimages: false,
			Trusted:   false,
		}
		meterGas = true
		gas      = uint64(1e6)
	)

	env := newMockEnv(address, config, meterGas, gas)

	r.Equal(gas, env.Gas())

	// GetGasLeft() costs gas, so the cost of that operation must be subtracted
	// from the total gas.
	getGasLeftOpCost := env.table[GetGasLeft_OpCode].constantGas
	gas -= getGasLeftOpCost
	r.Equal(gas, env.GetGasLeft())
	r.Equal(gas, env.Gas())
}

func TestBlockOps_Minimal(t *testing.T) {
	var (
		r       = require.New(t)
		address = common.HexToAddress("0xc0ffee0001")
		config  = EnvConfig{
			Static:    true,
			Ephemeral: false,
			Preimages: false,
			Trusted:   false,
		}
		meterGas = true
		gas      = uint64(1e6)
	)

	env := newMockEnv(address, config, meterGas, gas)

	r.Equal(env.block.GetHash(0), env.GetBlockHash(0))
	r.Equal(env.block.GasLimit(), env.GetBlockGasLimit())
	r.Equal(env.block.BlockNumber(), env.GetBlockNumber())
	r.Equal(env.block.Timestamp(), env.GetBlockTimestamp())
	r.Equal(env.block.Difficulty(), env.GetBlockDifficulty())
	r.Equal(env.block.BaseFee(), env.GetBlockBaseFee())
	r.Equal(env.block.Coinbase(), env.GetBlockCoinbase())
	r.Equal(env.block.Random(), env.GetPrevRandom())
}

func TestCallOps_Minimal(t *testing.T) {
	var (
		r       = require.New(t)
		address = common.HexToAddress("0xc0ffee0001")
		config  = EnvConfig{
			Static:    true,
			Ephemeral: false,
			Preimages: false,
			Trusted:   false,
		}
		meterGas = true
		gas      = uint64(1e6)
	)

	env := newMockEnv(address, config, meterGas, gas)

	r.Equal(env.call.TxGasPrice(), env.GetTxGasPrice())
	r.Equal(env.call.TxOrigin(), env.GetTxOrigin())
	r.Equal(env.call.CallData(), env.GetCallData())
	r.Equal(env.call.CallDataSize(), env.GetCallDataSize())
	r.Equal(env.call.Caller(), env.GetCaller())
	r.Equal(env.call.CallValue(), env.GetCallValue())
}
