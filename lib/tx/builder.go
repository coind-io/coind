package tx

import (
	"errors"
	"fmt"
	"math/big"
)

type Builder struct {
	err     error
	fee     *big.Int
	coins   []*Coin
	exports []Export
	refund  *Address
}

func NewBuilder() *Builder {
	builder := new(Builder)
	builder.coins = []*Coin{}
	builder.exports = []Export{}
	return builder
}

func (b *Builder) SetFee(fee *big.Int) {
	if b.err != nil {
		return
	}
	b.fee = fee
	return
}

func (b *Builder) SetRefund(refund *Address) {
	if b.err != nil {
		return
	}
	b.refund = refund
	return
}

func (b *Builder) SetCoins(coins []*Coin) {
	if b.err != nil {
		return
	}
	b.coins = coins
	return
}

func (b *Builder) AddExport(address *Address, amount *big.Int) {
	if b.err != nil {
		return
	}
	atype := address.Type()
	etype := ExportType(atype)

	exp, err := NewExport(etype)
	if err != nil {
		b.err = err
		return
	}

	err = exp.SetAddress(address)
	if err != nil {
		b.err = err
		return
	}

	exp.SetAmount(amount)
	b.exports = append(b.exports, exp)
	return
}

func (b *Builder) Build() (*Transfer, error) {
	if b.err != nil {
		return nil, b.err
	}
	// 计算总输出
	out_total := big.NewInt(0).Add(big.NewInt(0), b.fee)
	for _, exp := range b.exports {
		out_total = big.NewInt(0).Add(out_total, exp.GetAmount())
	}
	// 挑选最佳输入
	in_total := big.NewInt(0)
	used := []*Coin{}
	for i := 0; i < len(b.coins) && in_total.Cmp(out_total) == -1; i++ {
		coin := b.coins[i]
		used = append(used, coin)
		in_total = big.NewInt(0).Add(in_total, coin.Amount)
	}
	b.coins = used
	// 配平
	if in_total.Cmp(out_total) == -1 {
		return nil, fmt.Errorf("shortage of coin %d < %d %d", in_total, out_total, len(b.coins))
	}
	// 找零
	balance := big.NewInt(0).Sub(in_total, out_total)
	if b.refund == nil {
		return nil, errors.New("refund address is empty")
	}
	if balance.Cmp(big.NewInt(0)) == 1 {
		b.AddExport(b.refund, balance)
		if b.err != nil {
			return nil, b.err
		}
	}
	// 构建
	tx := NewTransfer()
	tx.Exports = b.exports
	for _, coin := range b.coins {
		imp := NewImport()
		imp.TxHash = coin.TxHash
		imp.TxIndex = coin.TxIndex
		imp.Redeem = coin.Redeem
		tx.AddImport(imp)
	}
	return tx, nil
}

func (b *Builder) BuildGenesis() (*Transfer, error) {
	if b.err != nil {
		return nil, b.err
	}
	// 计算总输出
	out_total := big.NewInt(0)
	for _, exp := range b.exports {
		out_total = big.NewInt(0).Add(out_total, exp.GetAmount())
	}
	// 压缩所有输出
	if b.refund == nil {
		return nil, errors.New("refund address is empty")
	}
	b.exports = []Export{}
	b.AddExport(b.refund, out_total)
	if b.err != nil {
		return nil, b.err
	}
	// 构建
	tx := NewTransfer()
	tx.Exports = b.exports
	imp := NewImport()
	imp.Genesis = true
	tx.AddImport(imp)
	return tx, nil
}
