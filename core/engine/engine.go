package engine

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"sync"
	"time"

	"github.com/coind-io/coind/lib/crypto"
	"github.com/coind-io/coind/lib/tx"
)

type Engine struct {
	deps     *Deps
	equitypk *crypto.PrivKey
	refundpk *crypto.PrivKey
	mu       sync.Mutex
}

func NewEngine(deps *Deps) (*Engine, error) {
	engine := new(Engine)
	engine.deps = deps
	err := engine.deps.Verify()
	if err != nil {
		return nil, err
	}
	err = engine.makeEquityPK()
	if err != nil {
		return nil, err
	}
	err = engine.makeRefundPK()
	if err != nil {
		return nil, err
	}
	return engine, nil
}

func (eng *Engine) makeEquityPK() error {
	err := os.MkdirAll(eng.deps.datadir, 0777)
	if err != nil {
		return err
	}
	filename := path.Join(eng.deps.datadir, "equity.key")
	exists, err := PathExists(filename)
	if err != nil {
		return err
	}
	if exists == false {
		pk, _, err := crypto.GenerateKeyPair()
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filename, pk.ToBytes(), 0666)
		if err != nil {
			return err
		}
	}
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	pk := crypto.NewPrivKey()
	err = pk.Update(raw)
	if err != nil {
		return err
	}
	eng.equitypk = pk
	// 显示地址
	address, err := tx.NewAddressFromPubKey(pk.PubKey())
	if err != nil {
		return err
	}
	fmt.Println("using miner address: ", address)
	return nil
}

func (eng *Engine) makeRefundPK() error {
	err := os.MkdirAll(eng.deps.datadir, 0777)
	if err != nil {
		return err
	}
	filename := path.Join(eng.deps.datadir, "refund.key")
	exists, err := PathExists(filename)
	if err != nil {
		return err
	}
	if exists == false {
		pk, _, err := crypto.GenerateKeyPair()
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filename, pk.ToBytes(), 0666)
		if err != nil {
			return err
		}
	}
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	pk := crypto.NewPrivKey()
	err = pk.Update(raw)
	if err != nil {
		return err
	}
	eng.refundpk = pk
	// 显示私钥
	fmt.Println(pk)
	// 显示地址
	address, err := tx.NewAddressFromPubKey(pk.PubKey())
	if err != nil {
		return err
	}
	fmt.Println("using refund address: ", address)
	return nil
}

func (eng *Engine) MainLoop() {
	first := eng.deps.chain.Genesis()
	last := eng.deps.chain.Best()
	slots := (uint64(time.Now().Unix()) - first.Timestamp()) / 10
	if slots == last.Height()+1 {
		eng.createFork()
		return
	}
	return
}

func (eng *Engine) createFork() {
	eng.mu.Lock()
	defer eng.mu.Unlock()
	// 构造空区块
	newb, err := eng.deps.chain.Best().Fork()
	if err != nil {
		return
	}
	// 加入奖励性交易
	refund, err := tx.NewAddressFromPubKey(eng.refundpk.PubKey())
	if err != nil {
		return
	}
	rewardtx, err := tx.NewGenesisTx(refund, big.NewInt(100000000))
	if err != nil {
		return
	}
	err = newb.AddTransfer(rewardtx)
	if err != nil {
		return
	}
	// 加入内存交易
	txlist := eng.deps.memtx.TxList()
	for _, cointx := range txlist {
		err := newb.AddTransfer(cointx)
		if err != nil {
			return
		}
	}
	eng.deps.memtx.Reset()
	// 加入矿工信息
	newb.Header.Creator = eng.equitypk.PubKey()
	err = newb.Sign(eng.equitypk)
	if err != nil {
		return
	}
	// 提交区块
	err = eng.deps.chain.ExecuteBlock(newb)
	if err != nil {
		return
	}
	fmt.Println("miner forging new block success !", newb.Height())
}

func (eng *Engine) Close() error {
	return nil
}
