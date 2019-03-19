package httpapi

import (
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/coind-io/coind/lib/encoding"
	"github.com/coind-io/coind/lib/tx"
)

type SendModule struct {
	deps *Deps
}

func NewSendModule(deps *Deps) *SendModule {
	sm := new(SendModule)
	sm.deps = deps
	return sm
}

func (sm *SendModule) RawTx(ctx *gin.Context) {
	// 解析参数
	params := struct {
		Hex string `json:"hex"`
	}{}
	err := json.NewDecoder(ctx.Request.Body).Decode(&params)
	if err != nil {
		NewErrResp(err).Encode(ctx.Writer)
		return
	}
	if params.Hex == "" {
		NewErrResp(errors.New("illegal params")).Encode(ctx.Writer)
		return
	}
	rawtx, err := hex.DecodeString(params.Hex)
	if err != nil {
		NewErrResp(err).Encode(ctx.Writer)
		return
	}
	// 解码交易
	cointx := tx.NewTransfer()
	err = encoding.Unmarshal(rawtx, cointx)
	if err != nil {
		NewErrResp(err).Encode(ctx.Writer)
		return
	}
	txhash, err := cointx.Hash()
	if err != nil {
		NewErrResp(err).Encode(ctx.Writer)
		return
	}
	// 执行交易
	err = sm.deps.memtx.ExecuteTx(cointx)
	if err != nil {
		NewErrResp(err).Encode(ctx.Writer)
		return
	}
	NewResp().Put("txhash", txhash).Encode(ctx.Writer)
	return
}
