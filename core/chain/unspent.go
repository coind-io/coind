package chain

import (
	"bytes"
	"errors"

	"github.com/boltdb/bolt"

	"github.com/coind-io/coind/lib/encoding"
	"github.com/coind-io/coind/lib/tx"
)

type UnspentSchema struct {
	dbtx *bolt.Tx
}

func NewUnspentSchema(dbtx *bolt.Tx) *UnspentSchema {
	us := new(UnspentSchema)
	us.dbtx = dbtx
	return us
}

func (us *UnspentSchema) MakeBucket() error {
	_, err := us.dbtx.CreateBucketIfNotExists([]byte{BUCKET_UNSPENT_HASH_PREFIX})
	if err != nil {
		return err
	}
	_, err = us.dbtx.CreateBucketIfNotExists([]byte{BUCKET_UNSPENT_HAVE_PREFIX})
	if err != nil {
		return err
	}
	return nil
}

func (us *UnspentSchema) FetchByKey(key []byte) (*tx.Unspent, error) {
	uraw := us.dbtx.Bucket([]byte{BUCKET_UNSPENT_HASH_PREFIX}).Get(key)
	if uraw == nil {
		return nil, errors.New("unspent does not exist")
	}
	unspent := tx.NewUnspent()
	err := encoding.Unmarshal(uraw, unspent)
	if err != nil {
		return nil, err
	}
	return unspent, nil
}

func (us *UnspentSchema) FetchByOwner(owner *tx.Address) ([]*tx.Unspent, error) {
	cursor := us.dbtx.Bucket([]byte{BUCKET_UNSPENT_HAVE_PREFIX}).Cursor()
	dst := make([]*tx.Unspent, 0, 20)
	for ikey, key := cursor.Seek(owner.ToBytes()); ikey != nil; ikey, key = cursor.Next() {
		if bytes.HasPrefix(ikey, owner.ToBytes()) == false {
			break
		}
		unspent, err := us.FetchByKey(key)
		if err != nil {
			return nil, err
		}
		dst = append(dst, unspent)
	}
	return dst, nil
}

func (us *UnspentSchema) Burning(key []byte) error {
	// 获取必要的信息
	unspent, err := us.FetchByKey(key)
	if err != nil {
		return err
	}
	owner, err := unspent.Export.GetAddress()
	if err != nil {
		return err
	}
	// 扣除余额
	err = NewBalanceSchema(us.dbtx).Decrease(owner, unspent.Export.GetAmount())
	if err != nil {
		return err
	}
	// 删除索引
	ikey := encoding.BConnect(owner.ToBytes(), key)
	err = us.dbtx.Bucket([]byte{BUCKET_UNSPENT_HAVE_PREFIX}).Delete(ikey)
	if err != nil {
		return err
	}
	// 删除数据
	err = us.dbtx.Bucket([]byte{BUCKET_UNSPENT_HASH_PREFIX}).Delete(key)
	if err != nil {
		return err
	}
	return nil
}

func (us *UnspentSchema) WriteUnspent(unspent *tx.Unspent) error {
	// 存入数据库
	uraw, err := encoding.Marshal(unspent)
	if err != nil {
		return err
	}
	err = us.dbtx.Bucket([]byte{BUCKET_UNSPENT_HASH_PREFIX}).Put(unspent.Key(), uraw)
	if err != nil {
		return err
	}
	// 建立未花费与地址的关联
	owner, err := unspent.Export.GetAddress()
	if err != nil {
		return err
	}
	ikey := encoding.BConnect(owner.ToBytes(), unspent.Key())
	err = us.dbtx.Bucket([]byte{BUCKET_UNSPENT_HAVE_PREFIX}).Put(ikey, unspent.Key())
	if err != nil {
		return err
	}
	// 添加余额
	err = NewBalanceSchema(us.dbtx).Increase(owner, unspent.Export.GetAmount())
	if err != nil {
		return err
	}
	return nil
}
