package memtx

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/coind-io/coind/lib/hash"
	"github.com/coind-io/coind/lib/tx"
)

type Transfers struct {
	dels map[string]bool
	mmap map[string]*tx.Transfer
}

type Unspents struct {
	dels map[string]bool
	mmap map[string]*tx.Unspent
}

type Balances struct {
	dels map[string]bool
	mmap map[string]*big.Int
}

type Owners struct {
	dels map[string]bool
	mmap map[string]map[string]*tx.Coin
}

type Store struct {
	deps     *Deps
	parent   *Store
	txs      *Transfers
	utxos    *Unspents
	balances *Balances
}

func newStore(deps *Deps) *Store {
	store := new(Store)
	store.deps = deps
	store.txs = &Transfers{
		dels: make(map[string]bool),
		mmap: make(map[string]*tx.Transfer),
	}
	store.utxos = &Unspents{
		dels: make(map[string]bool),
		mmap: make(map[string]*tx.Unspent),
	}
	store.balances = &Balances{
		dels: make(map[string]bool),
		mmap: make(map[string]*big.Int),
	}
	return store
}

func (ms *Store) Frok() *Store {
	dst := newStore(ms.deps)
	dst.parent = ms
	return dst
}

func (ms *Store) Merge() *Store {
	if ms.parent == nil {
		return ms
	}
	for k, v := range ms.txs.dels {
		ms.parent.txs.dels[k] = v
		delete(ms.parent.txs.mmap, k)
	}
	for k, v := range ms.txs.mmap {
		ms.parent.txs.mmap[k] = v
	}
	for k, v := range ms.utxos.dels {
		ms.parent.utxos.dels[k] = v
		delete(ms.parent.utxos.mmap, k)
	}
	for k, v := range ms.utxos.mmap {
		ms.parent.utxos.mmap[k] = v
	}
	for k, v := range ms.balances.dels {
		ms.parent.balances.dels[k] = v
		delete(ms.parent.balances.mmap, k)
	}
	for k, v := range ms.balances.mmap {
		ms.parent.balances.mmap[k] = v
	}
	return ms.parent
}

func (ms *Store) Tx(txhash *hash.Hash) (*tx.Transfer, error) {
	var dst *tx.Transfer
	for m := ms; m != nil; m = m.parent {
		if m.txs.dels[txhash.String()] == true {
			return nil, errors.New("transfer does not exist")
		}
		if m.txs.mmap[txhash.String()] == nil {
			continue
		}
		dst = m.txs.mmap[txhash.String()]
		break
	}
	if dst != nil {
		return dst, nil
	}
	return ms.deps.chain.Tx(txhash)
}

func (ms *Store) WriteTx(cointx *tx.Transfer) (*hash.Hash, error) {
	txhash, err := cointx.Hash()
	if err != nil {
		return nil, err
	}
	ctx, _ := ms.Tx(txhash)
	if ctx != nil {
		return nil, errors.New("transfer already exists")
	}
	ms.txs.mmap[txhash.String()] = cointx
	delete(ms.txs.dels, txhash.String())
	return txhash, nil
}

func (ms *Store) FetchBalanceByOwner(owner *tx.Address) *big.Int {
	var dst *big.Int
	bkey := owner.String()
	for m := ms; m != nil; m = m.parent {
		if m.balances.dels[bkey] == true {
			return big.NewInt(0)
		}
		if m.balances.mmap[bkey] == nil {
			continue
		}
		dst = m.balances.mmap[bkey]
		break
	}
	if dst != nil {
		return dst
	}
	return ms.deps.chain.FetchBalanceByOwner(owner)
}

func (ms *Store) IncreaseBalance(owner *tx.Address, amount *big.Int) {
	balance := ms.FetchBalanceByOwner(owner)
	balance = big.NewInt(0).Add(balance, amount)
	ms.balances.mmap[owner.String()] = balance
	delete(ms.balances.dels, owner.String())
	return
}

func (ms *Store) DecreaseBalance(owner *tx.Address, amount *big.Int) {
	balance := ms.FetchBalanceByOwner(owner)
	balance = big.NewInt(0).Sub(balance, amount)
	ms.balances.mmap[owner.String()] = balance
	delete(ms.balances.dels, owner.String())
	return
}

func (ms *Store) FetchUnspentByKey(ukey []byte) (*tx.Unspent, error) {
	var dst *tx.Unspent
	for m := ms; m != nil; m = m.parent {
		if m.utxos.dels[string(ukey)] == true {
			return nil, errors.New("unspent does not exist")
		}
		if m.utxos.mmap[string(ukey)] == nil {
			continue
		}
		dst = m.utxos.mmap[string(ukey)]
		break
	}
	if dst != nil {
		return dst, nil
	}
	return ms.deps.chain.FetchUnspentByKey(ukey)
}

func (ms *Store) BurringUnspent(ukey []byte) error {
	// 减少余额
	unspent, err := ms.FetchUnspentByKey(ukey)
	if err != nil {
		return err
	}
	owner, err := unspent.Export.GetAddress()
	if err != nil {
		return err
	}
	ms.DecreaseBalance(owner, unspent.Export.GetAmount())
	// 删除数据
	delete(ms.utxos.mmap, string(ukey))
	ms.utxos.dels[string(ukey)] = true
	return nil
}

func (ms *Store) WriteUnspent(unspent *tx.Unspent) error {
	// 增加余额
	owner, err := unspent.Export.GetAddress()
	if err != nil {
		return err
	}
	ms.IncreaseBalance(owner, unspent.Export.GetAmount())
	// 加入数据
	ukey := unspent.Key()
	ms.utxos.mmap[string(ukey)] = unspent
	delete(ms.utxos.dels, string(ukey))
	return nil
}

func (ms *Store) ExecuteTx(cointx *tx.Transfer) error {
	// 拒绝创世交易
	for _, imp := range cointx.Imports {
		if imp.Genesis == true {
			return errors.New("not accept genesis transfer")
		}
	}
	// 保存交易
	txhash, err := ms.WriteTx(cointx)
	if err != nil {
		return err
	}
	// 转换输入为未花费
	unspents := make([]*tx.Unspent, 0, len(cointx.Imports))
	for _, imp := range cointx.Imports {
		unspent, err := ms.FetchUnspentByKey(imp.Key())
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
	// 计算总输入
	input := big.NewInt(0)
	for _, unspent := range unspents {
		input = big.NewInt(0).Add(input, unspent.Export.GetAmount())
	}
	// 额度验证
	if output.Cmp(input) == 1 {
		return fmt.Errorf("shortage of coin %d < %d", input, output)
	}
	// 解锁输入
	for i, unspent := range unspents {
		err := unspent.Export.UnLock(&tx.Voucher{
			Digest: txhash.Bytes(),
			Redeem: cointx.Imports[i].Redeem,
			Signs:  cointx.Signature,
		})
		if err != nil {
			return err
		}
	}
	// 燃烧已用未花费输出
	for _, unspent := range unspents {
		err := ms.BurringUnspent(unspent.Key())
		if err != nil {
			return err
		}
	}
	// 构造新的未花费
	for i, exp := range cointx.Exports {
		unspent := tx.NewUnspent()
		unspent.Export = exp
		unspent.TxHash = txhash
		unspent.TxIndex = uint16(i)
		err := ms.WriteUnspent(unspent)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ms *Store) FetchCoinsByOwner(owner *tx.Address) ([]*tx.Coin, error) {
	coins, err := ms.deps.chain.FetchCoinsByOwner(owner)
	if err != nil {
		return nil, err
	}
	dst := make([]*tx.Coin, 0, len(coins)+len(ms.utxos.mmap))
	for _, coin := range coins {
		if ms.utxos.dels[string(coin.Key())] == true {
			continue
		}
		dst = append(dst, coin)
	}
	for _, unspents := range ms.utxos.mmap {
		addr, err := unspents.Export.GetAddress()
		if err != nil {
			return nil, err
		}
		if owner.String() != addr.String() {
			continue
		}
		coin := tx.NewCoin()
		coin.TxHash = unspents.TxHash
		coin.TxIndex = unspents.TxIndex
		coin.Amount = unspents.Export.GetAmount()
		dst = append(dst, coin)
	}
	return dst, nil
}
