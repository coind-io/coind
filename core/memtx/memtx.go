package memtx

import (
	"math/big"
	"sync"

	"github.com/coind-io/coind/lib/tx"
)

type MemTx struct {
	deps   *Deps
	store  *Store
	txlist []*tx.Transfer
	mu     sync.Mutex
}

func NewMemTx(deps *Deps) (*MemTx, error) {
	mtx := new(MemTx)
	mtx.deps = deps
	mtx.store = newStore(deps)
	mtx.txlist = make([]*tx.Transfer, 0, 128)
	err := mtx.deps.Verify()
	if err != nil {
		return nil, err
	}
	return mtx, nil
}

func (mtx *MemTx) ExecuteTx(cointx *tx.Transfer) error {
	mtx.mu.Lock()
	defer mtx.mu.Unlock()
	store := mtx.store.Frok()
	err := store.ExecuteTx(cointx)
	if err != nil {
		return err
	}
	mtx.store = store.Merge()
	mtx.txlist = append(mtx.txlist, cointx)
	return nil
}

func (mtx *MemTx) FetchCoinsByOwner(owner *tx.Address) ([]*tx.Coin, error) {
	return mtx.store.FetchCoinsByOwner(owner)
}

func (mtx *MemTx) FetchBalanceByOwner(owner *tx.Address) *big.Int {
	return mtx.store.FetchBalanceByOwner(owner)
}

func (mtx *MemTx) TxList() []*tx.Transfer {
	return mtx.txlist
}

func (mtx *MemTx) Reset() {
	mtx.mu.Lock()
	defer mtx.mu.Unlock()
	mtx.store = newStore(mtx.deps)
	mtx.txlist = make([]*tx.Transfer, 0, 128)
	return
}
