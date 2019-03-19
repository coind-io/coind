package chain

import (
	"math/big"
	"os"
	"path"

	"github.com/boltdb/bolt"

	"github.com/coind-io/coind/lib/block"
	"github.com/coind-io/coind/lib/hash"
	"github.com/coind-io/coind/lib/tx"
)

type Chain struct {
	deps *Deps
	kvdb *bolt.DB
}

func NewChain(deps *Deps) (*Chain, error) {
	chain := new(Chain)
	chain.deps = deps
	err := chain.deps.Verify()
	if err != nil {
		return nil, err
	}
	err = chain.makeDB()
	if err != nil {
		return nil, err
	}
	err = chain.makeBucket()
	if err != nil {
		return nil, err
	}
	err = chain.makeGenesis()
	if err != nil {
		return nil, err
	}
	return chain, nil
}

func (chain *Chain) makeDB() error {
	err := os.MkdirAll(chain.deps.datadir, 0777)
	if err != nil {
		return err
	}
	dbname := path.Join(chain.deps.datadir, "chain.db")
	kvdb, err := bolt.Open(dbname, 0666, nil)
	if err != nil {
		return err
	}
	chain.kvdb = kvdb
	return nil
}

func (chain *Chain) makeBucket() error {
	return chain.kvdb.Update(func(dbtx *bolt.Tx) error {
		err := NewBalanceSchema(dbtx).MakeBucket()
		if err != nil {
			return err
		}
		err = NewUnspentSchema(dbtx).MakeBucket()
		if err != nil {
			return err
		}
		err = NewTransferSchema(dbtx).MakeBucket()
		if err != nil {
			return err
		}
		err = NewBlockSchema(dbtx).MakeBucket()
		if err != nil {
			return err
		}
		return nil
	})
}

func (chain *Chain) makeGenesis() error {
	return chain.kvdb.Update(func(dbtx *bolt.Tx) error {
		genesis := NewBlockSchema(dbtx).Genesis()
		if genesis != nil {
			return nil
		}
		genesis = block.NewBlock(0)
		err := NewBlockSchema(dbtx).WriteBlock(genesis)
		if err != nil {
			return err
		}
		return nil
	})
}

func (chain *Chain) Best() *block.Block {
	var dst *block.Block
	chain.kvdb.View(func(dbtx *bolt.Tx) error {
		dst = NewBlockSchema(dbtx).Best()
		return nil
	})
	return dst
}

func (chain *Chain) Genesis() *block.Block {
	var dst *block.Block
	chain.kvdb.View(func(dbtx *bolt.Tx) error {
		dst = NewBlockSchema(dbtx).Genesis()
		return nil
	})
	return dst
}

func (chain *Chain) Blocks(height uint64, limit uint16) ([]*block.Block, error) {
	var dst []*block.Block
	return dst, chain.kvdb.Update(func(dbtx *bolt.Tx) error {
		bs, err := NewBlockSchema(dbtx).Blocks(height, limit)
		if err != nil {
			return err
		}
		dst = bs
		return nil
	})
}

func (chain *Chain) Block(bhash *hash.Hash) (*block.Block, error) {
	var dst *block.Block
	return dst, chain.kvdb.View(func(dbtx *bolt.Tx) error {
		b, err := NewBlockSchema(dbtx).Block(bhash)
		if err != nil {
			return err
		}
		dst = b
		return nil
	})
}

func (chain *Chain) BlockIndex(height uint64) (*hash.Hash, error) {
	var dst *hash.Hash
	return dst, chain.kvdb.View(func(dbtx *bolt.Tx) error {
		bhash, err := NewBlockSchema(dbtx).BlockIndex(height)
		if err != nil {
			return err
		}
		dst = bhash
		return nil
	})
}

func (chain *Chain) Tx(txhash *hash.Hash) (*tx.Transfer, error) {
	var dst *tx.Transfer
	return dst, chain.kvdb.View(func(dbtx *bolt.Tx) error {
		cointx, err := NewTransferSchema(dbtx).Tx(txhash)
		if err != nil {
			return err
		}
		dst = cointx
		return nil
	})
}

func (chain *Chain) FetchUnspentByKey(key []byte) (*tx.Unspent, error) {
	var dst *tx.Unspent
	return dst, chain.kvdb.View(func(dbtx *bolt.Tx) error {
		unspent, err := NewUnspentSchema(dbtx).FetchByKey(key)
		if err != nil {
			return err
		}
		dst = unspent
		return nil
	})
}

func (chain *Chain) FetchCoinsByOwner(owner *tx.Address) ([]*tx.Coin, error) {
	var dst []*tx.Coin
	return dst, chain.kvdb.View(func(dbtx *bolt.Tx) error {
		unspents, err := NewUnspentSchema(dbtx).FetchByOwner(owner)
		if err != nil {
			return nil
		}
		coins := make([]*tx.Coin, 0, len(unspents))
		for _, unspent := range unspents {
			coin := tx.NewCoin()
			coin.TxHash = unspent.TxHash
			coin.TxIndex = unspent.TxIndex
			coin.Amount = unspent.Export.GetAmount()
			coins = append(coins, coin)
		}
		dst = coins
		return nil
	})
}

func (chain *Chain) FetchBalanceByOwner(owner *tx.Address) *big.Int {
	var dst *big.Int
	chain.kvdb.View(func(dbtx *bolt.Tx) error {
		dst = NewBalanceSchema(dbtx).Fetch(owner)
		return nil
	})
	return dst
}

func (chain *Chain) ExecuteBlock(b *block.Block) error {
	return chain.kvdb.Update(func(dbtx *bolt.Tx) error {
		return NewBlockSchema(dbtx).ExecuteBlock(b)
	})
}

func (chain *Chain) Close() error {
	err := chain.kvdb.Close()
	if err != nil {
		return err
	}
	return nil
}
