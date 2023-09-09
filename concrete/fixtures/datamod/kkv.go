/* Autogenerated file. Do not edit manually. */

package datamod

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/concrete/codegen/datamod/codec"
	"github.com/ethereum/go-ethereum/concrete/crypto"
	"github.com/ethereum/go-ethereum/concrete/lib"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = common.Big1
	_ = codec.EncodeAddress
)

var (
	KkvDefaultKey = crypto.Keccak256([]byte("datamod.v1.Kkv"))
)

type KkvRow struct {
	lib.DatastoreStruct
}

func NewKkvRow(store lib.DatastoreSlot) *KkvRow {
	sizes := []int{32}
	return &KkvRow{*lib.NewDatastoreStruct(store, sizes)}
}

func (v *KkvRow) Get() (
	common.Hash,
) {
	return codec.DecodeHash(32, v.GetField(0))
}

func (v *KkvRow) Set(
	value common.Hash,
) {
	v.SetField(0, codec.EncodeHash(32, value))
}

func (v *KkvRow) GetValue() common.Hash {
	data := v.GetField(0)
	return codec.DecodeHash(32, data)
}

func (v *KkvRow) SetValue(value common.Hash) {
	data := codec.EncodeHash(32, value)
	v.SetField(0, data)
}

type Kkv struct {
	store lib.DatastoreSlot
}

func NewKkv(ds lib.Datastore) *Kkv {
	return &Kkv{ds.Get(KkvDefaultKey)}
}

func NewKkvWithKey(ds lib.Datastore, key []byte) *Kkv {
	return &Kkv{ds.Get(key)}
}

func (m *Kkv) Get(
	key1 common.Hash,
	key2 common.Hash,
) *KkvRow {
	store := m.store.Mapping().GetNested(
		codec.EncodeHash(32, key1),
		codec.EncodeHash(32, key2),
	)
	return NewKkvRow(store)
}