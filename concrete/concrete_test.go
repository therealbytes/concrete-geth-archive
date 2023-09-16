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

package concrete

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/concrete/api"
	"github.com/stretchr/testify/require"
)

type pcSet struct {
	blockNumber uint64
	precompiles map[common.Address]Precompile
}

type pcSingle struct {
	blockNumber uint64
	address     common.Address
	precompile  Precompile
}

type pcBlank struct{}

func (pc *pcBlank) IsStatic(input []byte) bool {
	return true
}

func (pc *pcBlank) Finalise(API api.Environment) error {
	return nil
}

func (pc *pcBlank) Commit(API api.Environment) error {
	return nil
}

func (pc *pcBlank) Run(API api.Environment, input []byte) ([]byte, error) {
	return []byte{}, nil
}

var _ Precompile = &pcBlank{}

var (
	addrIncl1 = common.BytesToAddress([]byte{128})
	addrIncl2 = common.BytesToAddress([]byte{129})
	addrExcl  = common.BytesToAddress([]byte{130})
	// Block numbers are deliberately not in order
	pcSets = []pcSet{
		{
			blockNumber: 0,
			precompiles: map[common.Address]Precompile{
				addrIncl1: &pcBlank{},
				addrIncl2: &pcBlank{},
			},
		},
		{
			blockNumber: 20,
			precompiles: map[common.Address]Precompile{
				addrIncl2: &pcBlank{},
			},
		},
		{
			blockNumber: 10,
			precompiles: map[common.Address]Precompile{
				addrIncl1: &pcBlank{},
			},
		},
		{
			blockNumber: 40,
			precompiles: map[common.Address]Precompile{
				addrIncl1: &pcBlank{},
				addrIncl2: &pcBlank{},
			},
		},
		{
			blockNumber: 30,
			precompiles: map[common.Address]Precompile{},
		},
	}
	pcSingles = []pcSingle{
		{
			blockNumber: 5,
			address:     addrIncl1,
			precompile:  &pcBlank{},
		},
		{
			blockNumber: 15,
			address:     addrIncl2,
			precompile:  &pcBlank{},
		},
	}
)

func verifyPrecompileSet(t *testing.T, config *GenericPrecompileRegistry, num uint64, p pcSet) {
	r := require.New(t)
	// Assert that all the provided addresses have been returned and all the returned
	// addresses were provided
	addresses := config.ActivePrecompiles(num)
	r.Len(addresses, len(p.precompiles))
	for _, address := range addresses {
		_, ok := p.precompiles[address]
		r.True(ok)
	}
	for address := range p.precompiles {
		r.Contains(addresses, address)
	}
	// Assert that all active addresses map to the correct precompile
	for address, setPc := range p.precompiles {
		configPc, ok := config.Precompile(address, num)
		r.True(ok)
		r.Equal(setPc, configPc)
	}
	// Assert that inactive addresses do not map to a precompile
	pc, ok := config.Precompile(addrExcl, num)
	r.Nil(pc)
	r.False(ok)
}

func verifyPrecompileSingle(t *testing.T, config *GenericPrecompileRegistry, num uint64, p pcSingle) {
	r := require.New(t)
	// Assert that all the provided addresses have been returned and all the returned
	// addresses were provided
	addresses := config.ActivePrecompiles(num)
	r.Len(addresses, 1)
	r.Equal(p.address, addresses[0])
	// Assert that all active addresses map to the correct precompile
	configPc, ok := config.Precompile(p.address, num)
	r.True(ok)
	r.Equal(p.precompile, configPc)
	// Assert that inactive addresses do not map to a precompile
	pc, ok := config.Precompile(addrExcl, num)
	r.Nil(pc)
	r.False(ok)
}

func TestConcreteConfig(t *testing.T) {
	t.Run("AddPrecompiles", func(t *testing.T) {
		config := NewRegistry()
		for _, d := range pcSets {
			config.AddPrecompiles(d.blockNumber, d.precompiles)
		}
		for _, d := range pcSets {
			require.Panics(t, func() {
				config.AddPrecompiles(d.blockNumber, d.precompiles)
			})
		}
		for _, d := range pcSets {
			// Check that the precompiles are returned correctly for the first, second and last
			// block in each range
			for _, delta := range []uint64{0, 1, 9} {
				blockNumber := d.blockNumber + delta
				verifyPrecompileSet(t, config, blockNumber, d)
			}
		}
	})
	t.Run("AddPrecompile", func(t *testing.T) {
		t.Run("OnEmpty", func(t *testing.T) {
			config := NewRegistry()
			for _, d := range pcSingles {
				config.AddPrecompile(d.blockNumber, d.address, d.precompile)
			}
			for _, d := range pcSingles {
				require.Panics(t, func() {
					config.AddPrecompile(d.blockNumber, d.address, d.precompile)
				})
			}
			for _, d := range pcSingles {
				// Check that the precompiles are returned correctly for the first, second and last
				// block in each range
				for _, delta := range []uint64{0, 1, 9} {
					blockNumber := d.blockNumber + delta
					verifyPrecompileSingle(t, config, blockNumber, d)
				}
			}
		})
		t.Run("OnExisting", func(t *testing.T) {
			config := NewRegistry()
			for _, d := range pcSets {
				config.AddPrecompiles(d.blockNumber, d.precompiles)
			}
			for _, d := range pcSingles {
				config.AddPrecompile(d.blockNumber, d.address, d.precompile)
			}
			for _, d := range pcSingles {
				require.Panics(t, func() {
					config.AddPrecompile(d.blockNumber, d.address, d.precompile)
				})
			}
			for _, d := range pcSets {
				blockNumber := d.blockNumber
				verifyPrecompileSet(t, config, blockNumber, d)
			}
			for _, d := range pcSingles {
				blockNumber := d.blockNumber
				verifyPrecompileSingle(t, config, blockNumber, d)
			}
		})
	})
}
