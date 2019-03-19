package tx

import (
	"math/big"
)

func NewGenesisTx(address *Address, amount *big.Int) (*Transfer, error) {
	builder := NewBuilder()
	builder.SetFee(big.NewInt(0))
	builder.SetRefund(address)
	builder.AddExport(address, amount)
	return builder.BuildGenesis()
}
