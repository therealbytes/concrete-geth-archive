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

package codec

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
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

func DecodeBool(_ int, data []byte) bool {
	return data[0]&1 == byte(0x01)
}

func EncodeBytes(size int, b []byte) []byte {
	return common.LeftPadBytes(b, size)
}

func DecodeBytes(size int, data []byte) []byte {
	return common.LeftPadBytes(data, size)
}

func EncodeInt(size int, i *big.Int) []byte {
	buf := make([]byte, size)
	if i.Sign() == -1 {
		for j := 0; j < size; j++ {
			buf[j] = 0xFF
		}
	}
	iBytes := math.U256Bytes(i)
	copy(buf, iBytes[len(iBytes)-size:])
	return buf
}

func DecodeInt(size int, data []byte) *big.Int {
	b := new(big.Int).SetBytes(data)
	if data[0]&0x80 != 0 {
		for i := len(data); i < size; i++ {
			b.Or(b, new(big.Int).Lsh(big.NewInt(0xFF), uint(8*i)))
		}
	}
	return b
}

func EncodeUint(size int, i *big.Int) []byte {
	buf := make([]byte, size)
	iBytes := math.U256Bytes(i)
	copy(buf, iBytes[len(iBytes)-size:])
	return buf
}

func DecodeUint(size int, data []byte) *big.Int {
	return new(big.Int).SetBytes(data)
}
