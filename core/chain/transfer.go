package chain

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/boltdb/bolt"

	"github.com/coind-io/coind/lib/encoding"
	"github.com/coind-io/coind/lib/hash"
	"github.com/coind-io/coind/lib/tx"
)

type TransferSchema struct {
	dbtx *bolt.Tx
}

func NewTransferSchema(dbtx *bolt.Tx) *TransferSchema {
	txs := new(TransferSchema)
	txs.dbtx = dbtx
	return txs
}

func (txs *TransferSchema) MakeBucket() error {
	_, err := txs.dbtx.CreateBucketIfNotExists([]byte{BUCKET_TRANSFER_HASH_PREFIX})
	if err != nil {
		return err
	}
	return nil
}

func (txs *TransferSchema) Tx(txhash *hash.Hash) (*tx.Transfer, error) {
	txraw := txs.dbtx.Bucket([]byte{BUCKET_TRANSFER_HASH_PREFIX}).Get(txhash.Bytes())
	if txraw == nil {
		return nil, errors.New("transfer does not exist")
	}
	cointx := tx.NewTransfer()
	err := encoding.Unmarshal(txraw, cointx)
	if err != nil {
		return nil, err
	}
	return cointx, nil
}

func (txs *TransferSchema) WriteTx(cointx *tx.Transfer) (*hash.Hash, error) {
	// 序列化交易
	txhash, err := cointx.Hash()
	if err != nil {
		return nil, err
	}
	txraw, err := encoding.Marshal(cointx)
	if err != nil {
		return nil, err
	}
	// 写入数据库
	err = txs.dbtx.Bucket([]byte{BUCKET_TRANSFER_HASH_PREFIX}).Put(txhash.Bytes(), txraw)
	if err != nil {
		return nil, err
	}
	return txhash, nil
}

func (txs *TransferSchema) ExecuteTx(cointx *tx.Transfer) error {
	// 保存交易
	txhash, err := txs.WriteTx(cointx)
	if err != nil {
		return err
	}
	// 转换输入为未花费
	unspents := make([]*tx.Unspent, 0, len(cointx.Imports))
	for _, imp := range cointx.Imports {
		if imp.Genesis == true {
			break
		}
		unspent, err := NewUnspentSchema(txs.dbtx).FetchByKey(imp.Key())
		if err != nil {
			return err
		}
		unspents = append(unspents, unspent)
	}
	// 计算总输出
	output := big.NewInt(0)
	for _, exp := range cointx.Exports {
		output = big.NewInt(0).Add(output, exp.GetAmount())
	}
	//计算总输入
	input := big.NewInt(0)
	for i, imp := range cointx.Imports {
		if imp.Genesis == true {
			input = big.NewInt(0).Add(big.NewInt(0), output)
			break
		}
		input = big.NewInt(0).Add(input, unspents[i].Export.GetAmount())
	}
	// 额度验证
	if output.Cmp(input) == 1 {
		return fmt.Errorf("shortage of coin %d < %d", input, output)
	}
	// 解锁输入
	for i, imp := range cointx.Imports {
		if imp.Genesis == true {
			break
		}
		err = unspents[i].Export.UnLock(&tx.Voucher{
			Digest: txhash.Bytes(),
			Redeem: imp.Redeem,
			Signs:  cointx.Signature,
		})
		if err != nil {
			return err
		}
	}
	// 燃烧已用未花费输出
	for _, imp := range cointx.Imports {
		if imp.Genesis == true {
			break
		}
		err = NewUnspentSchema(txs.dbtx).Burning(imp.Key())
		if err != nil {
			return err
		}
	}
	// 构造新的未花费
	for i, e := range cointx.Exports {
		unspent := tx.NewUnspent()
		unspent.Export = e
		unspent.TxIndex = uint16(i)
		unspent.TxHash = txhash
		err := NewUnspentSchema(txs.dbtx).WriteUnspent(unspent)
		if err != nil {
			return err
		}
	}
	return nil
}
