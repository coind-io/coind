package chain

import (
	"errors"

	"github.com/boltdb/bolt"

	"github.com/coind-io/coind/lib/block"
	"github.com/coind-io/coind/lib/encoding"
	"github.com/coind-io/coind/lib/hash"
)

type BlockSchema struct {
	dbtx *bolt.Tx
}

func NewBlockSchema(dbtx *bolt.Tx) *BlockSchema {
	bs := new(BlockSchema)
	bs.dbtx = dbtx
	return bs
}

func (bs *BlockSchema) MakeBucket() error {
	_, err := bs.dbtx.CreateBucketIfNotExists([]byte{BUCKET_BLOCK_HASH_PREFIX})
	if err != nil {
		return err
	}
	_, err = bs.dbtx.CreateBucketIfNotExists([]byte{BUCKET_BLOCK_HEIGHT_PREFIX})
	if err != nil {
		return err
	}
	return nil
}

func (bs *BlockSchema) Blocks(height uint64, limit uint16) ([]*block.Block, error) {
	cursor := bs.dbtx.Bucket([]byte{BUCKET_BLOCK_HEIGHT_PREFIX}).Cursor()
	dst := make([]*block.Block, 0, limit)
	count := uint16(0)
	if limit == 0 {
		limit = 20
	}
	var ikey, hkey []byte
	if height == 0 {
		ikey, hkey = cursor.Last()
	} else {
		ikey, hkey = cursor.Seek(encoding.I2b(height))
	}
	for ; ikey != nil; ikey, hkey = cursor.Prev() {
		if count == limit {
			break
		}
		count = count + 1
		bhash := hash.NewHash256()
		err := bhash.Update(hkey)
		if err != nil {
			return nil, err
		}
		b, err := bs.Block(bhash)
		if err != nil {
			return nil, err
		}
		dst = append(dst, b)
	}
	return dst, nil
}

func (bs *BlockSchema) Block(bhash *hash.Hash) (*block.Block, error) {
	raw := bs.dbtx.Bucket([]byte{BUCKET_BLOCK_HASH_PREFIX}).Get(bhash.Bytes())
	if raw == nil {
		return nil, errors.New("block does not exist")
	}
	b := block.NewBlock(0)
	err := encoding.Unmarshal(raw, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (bs *BlockSchema) BlockIndex(height uint64) (*hash.Hash, error) {
	ikey := encoding.I2b(height)
	raw := bs.dbtx.Bucket([]byte{BUCKET_BLOCK_HEIGHT_PREFIX}).Get(ikey)
	if raw == nil {
		return nil, errors.New("block hash does not exist")
	}
	bhash := hash.NewHash256()
	err := bhash.Update(raw)
	if err != nil {
		return nil, err
	}
	return bhash, nil
}

func (bs *BlockSchema) Best() *block.Block {
	cursor := bs.dbtx.Bucket([]byte{BUCKET_BLOCK_HEIGHT_PREFIX}).Cursor()
	ikey, raw := cursor.Last()
	if ikey == nil {
		return nil
	}
	bhash := hash.NewHash256()
	err := bhash.Update(raw)
	if err != nil {
		panic(err)
	}
	b, _ := bs.Block(bhash)
	return b
}

func (bs *BlockSchema) Genesis() *block.Block {
	cursor := bs.dbtx.Bucket([]byte{BUCKET_BLOCK_HEIGHT_PREFIX}).Cursor()
	ikey, raw := cursor.First()
	if ikey == nil {
		return nil
	}
	bhash := hash.NewHash256()
	err := bhash.Update(raw)
	if err != nil {
		panic(err)
	}
	b, _ := bs.Block(bhash)
	return b
}

func (bs *BlockSchema) WriteBlock(b *block.Block) error {
	raw, err := encoding.Marshal(b)
	if err != nil {
		return err
	}
	bhash, err := b.Hash()
	if err != nil {
		return err
	}
	err = bs.dbtx.Bucket([]byte{BUCKET_BLOCK_HASH_PREFIX}).Put(bhash.Bytes(), raw)
	if err != nil {
		return err
	}
	// 建立高度索引
	ikey := encoding.I2b(b.Height())
	err = bs.dbtx.Bucket([]byte{BUCKET_BLOCK_HEIGHT_PREFIX}).Put(ikey, bhash.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (bs *BlockSchema) ExecuteBlock(b *block.Block) error {
	//保存区块
	err := bs.WriteBlock(b)
	if err != nil {
		return err
	}
	// 执行交易
	for _, cointx := range b.TxList {
		err := NewTransferSchema(bs.dbtx).ExecuteTx(cointx)
		if err != nil {
			return err
		}
	}
	return nil
}
