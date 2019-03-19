package chain

import (
	"math/big"

	"github.com/boltdb/bolt"

	"github.com/coind-io/coind/lib/tx"
)

type BalanceSchema struct {
	dbtx *bolt.Tx
}

func NewBalanceSchema(dbtx *bolt.Tx) *BalanceSchema {
	bs := new(BalanceSchema)
	bs.dbtx = dbtx
	return bs
}

func (bs *BalanceSchema) MakeBucket() error {
	_, err := bs.dbtx.CreateBucketIfNotExists([]byte{BUCKET_BALANCE_HASH_PREFIX})
	if err != nil {
		return err
	}
	return nil
}

func (bs *BalanceSchema) Fetch(owner *tx.Address) *big.Int {
	raw := bs.dbtx.Bucket([]byte{BUCKET_BALANCE_HASH_PREFIX}).Get(owner.ToBytes())
	if raw == nil {
		return big.NewInt(0)
	}
	return big.NewInt(0).SetBytes(raw)
}

func (bs *BalanceSchema) Decrease(owner *tx.Address, amount *big.Int) error {
	value := bs.Fetch(owner)
	value = big.NewInt(0).Sub(value, amount)
	key := owner.ToBytes()
	raw := value.Bytes()
	err := bs.dbtx.Bucket([]byte{BUCKET_BALANCE_HASH_PREFIX}).Put(key, raw)
	if err != nil {
		return err
	}
	return nil
}

func (bs *BalanceSchema) Increase(owner *tx.Address, amount *big.Int) error {
	value := bs.Fetch(owner)
	value = big.NewInt(0).Add(value, amount)
	key := owner.ToBytes()
	raw := value.Bytes()
	err := bs.dbtx.Bucket([]byte{BUCKET_BALANCE_HASH_PREFIX}).Put(key, raw)
	if err != nil {
		return err
	}
	return nil
}
