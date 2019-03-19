package httpapi

import (
	_ "encoding/hex"
	_ "net/http"
	_ "strconv"

	"github.com/gin-gonic/gin"

	_ "github.com/coind-io/coind/lib/hash"
	"github.com/coind-io/coind/lib/tx"
)

type FetchModule struct {
	deps *Deps
}

func NewFetchModule(deps *Deps) *FetchModule {
	fm := new(FetchModule)
	fm.deps = deps
	return fm
}

func (fm *FetchModule) Coins(ctx *gin.Context) {
	owner, err := tx.NewAddressFromWIF(ctx.Param("owner"))
	if err != nil {
		NewErrResp(err).Encode(ctx.Writer)
		return
	}
	coins, err := fm.deps.memtx.FetchCoinsByOwner(owner)
	if err != nil {
		NewErrResp(err).Encode(ctx.Writer)
		return
	}
	resp := NewResp().Put("coins", coins)
	ctx.JSON(200, resp.dict)
	return
}

func (fm *FetchModule) Balance(ctx *gin.Context) {
	owner, err := tx.NewAddressFromWIF(ctx.Param("owner"))
	if err != nil {
		NewErrResp(err).Encode(ctx.Writer)
		return
	}
	balance := fm.deps.memtx.FetchBalanceByOwner(owner)
	resp := NewResp().Put("balance", balance)
	ctx.JSON(200, resp.dict)
	return
}

/*
func (fm *FetchModule) Blocks(w http.ResponseWriter, r *http.Request) {
	bs, err := fm.deps.chain.Blocks(0, 0)
	if err != nil {
		NewErrResp(err).Encode(w)
		return
	}
	NewResp().Put("blocks", bs).Encode(w)
	return
}

func (fm *FetchModule) BlockIndex(w http.ResponseWriter, r *http.Request) {
	sindex := bone.GetValue(r, "index")
	nindex, err := strconv.ParseUint(sindex, 10, 64)
	if err != nil {
		NewErrResp(err).Encode(w)
		return
	}
	bhash, err := fm.deps.chain.BlockIndex(nindex)
	if err != nil {
		NewErrResp(err).Encode(w)
		return
	}
	NewResp().Put("bhash", bhash).Encode(w)
	return
}

func (fm *FetchModule) Block(w http.ResponseWriter, r *http.Request) {
	rawhash, err := hex.DecodeString(bone.GetValue(r, "hash"))
	if err != nil {
		NewErrResp(err).Encode(w)
		return
	}
	bhash := hash.NewHash256()
	err = bhash.Update(rawhash)
	if err != nil {
		NewErrResp(err).Encode(w)
		return
	}
	b, err := fm.deps.chain.Block(bhash)
	if err != nil {
		NewErrResp(err).Encode(w)
		return
	}
	NewResp().Put("block", b).Encode(w)
	return
}

func (fm *FetchModule) Tx(w http.ResponseWriter, r *http.Request) {
	rawhash, err := hex.DecodeString(bone.GetValue(r, "hash"))
	if err != nil {
		NewErrResp(err).Encode(w)
		return
	}
	txhash := hash.NewHash256()
	err = txhash.Update(rawhash)
	if err != nil {
		NewErrResp(err).Encode(w)
		return
	}
	cointx, err := fm.deps.chain.Tx(txhash)
	if err != nil {
		NewErrResp(err).Encode(w)
		return
	}
	NewResp().Put("tx", cointx).Encode(w)
	return
}

func (fm *FetchModule) Address(w http.ResponseWriter, r *http.Request) {
	address, err := tx.NewAddressFromWIF(bone.GetValue(r, "address"))
	if err != nil {
		NewErrResp(err).Encode(w)
		return
	}
	balance := fm.deps.memtx.FetchBalanceByOwner(address)
	NewResp().Put("balance", balance).Encode(w)
	return
}
*/
