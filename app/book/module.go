package coinbook

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/coind-io/coind/lib/hash"
)

type BookModule struct {
	deps *Deps
}

func NewBookModule(deps *Deps) *BookModule {
	bm := new(BookModule)
	bm.deps = deps
	return bm
}

func (bm *BookModule) Status(ctx *gin.Context) {
	lastblock := bm.deps.chain.Best()
	// 计算块哈希
	bkhash, err := lastblock.Hash()
	if err != nil {
		ctx.Error(err)
		return
	}
	// 获得块高
	bkheight := lastblock.Height()
	// 渲染页面
	ctx.HTML(200, "status.tpl", gin.H{
		"LastBlock": gin.H{
			"Height": bkheight,
			"Hash":   bkhash.String(),
		},
	})
	return
}

func (bm *BookModule) Block(ctx *gin.Context) {
	// 获得区块哈希
	bkhash, err := hash.NewHash256FromHex(ctx.Param("bkhash"))
	if err != nil {
		ctx.Error(err)
		return
	}
	// 获得原始区块
	rawbk, err := bm.deps.chain.Block(bkhash)
	if err != nil {
		ctx.Error(err)
		return
	}
	// 构造渲染区块
	renbk := new(Block)
	renbk.BkHeader = new(Header)
	renbk.BkPrev = nil
	renbk.BkNext = nil
	renbk.TxList = make([]*Transfer, 0, len(rawbk.TxHash))
	renbk.BkHeader.BkHash = bkhash.String()
	renbk.BkHeader.BkHeight = rawbk.Header.Height
	renbk.BkHeader.BkVersion = rawbk.Header.Version
	renbk.BkHeader.BkMerkleRoot = rawbk.Header.MerkleRoot.String()
	renbk.BkHeader.BkTime = time.Unix(int64(rawbk.Header.Timestamp), 0).Format("2006-01-02 15:04:05")
	renbk.BkHeader.BkSize = rawbk.Size()
	renbk.BkHeader.TxCount = uint64(len(rawbk.TxHash))
	// 计算确认数
	bestbk := bm.deps.chain.Best()
	renbk.BkHeader.Confirm = bestbk.Header.Height - rawbk.Header.Height
	// 附加前驱区块
	if rawbk.Header.PrevBlock.IsZero() == false {
		renbk.BkPrev = new(Header)
		renbk.BkPrev.BkHash = rawbk.Header.PrevBlock.String()
	}
	// 附加后继区块
	nbkhash, _ := bm.deps.chain.BlockIndex(rawbk.Header.Height + 1)
	if nbkhash != nil {
		renbk.BkNext = new(Header)
		renbk.BkNext.BkHash = nbkhash.String()
	}
	// 附加关联交易
	for _, txhash := range rawbk.TxHash {
		rawtx, err := bm.deps.chain.Tx(txhash)
		if err != nil {
			ctx.Error(err)
			return
		}
		// 转换引入
		imps := make([]*Import, 0, len(rawtx.Imports))
		for _, imp := range rawtx.Imports {
			if imp.Genesis == true {
				imps = append(imps, &Import{
					TxGenesis: true,
				})
				break
			}
			imptx, err := bm.deps.chain.Tx(imp.TxHash)
			if err != nil {
				ctx.Error(err)
				return
			}
			impaddr, err := imptx.Exports[int(imp.TxIndex)].GetAddress()
			if err != nil {
				ctx.Error(err)
				return
			}
			imps = append(imps, &Import{
				TxHash:     imp.TxHash.String(),
				ImpAddress: impaddr.String(),
				ImpAmount:  imptx.Exports[int(imp.TxIndex)].GetAmount().Uint64() / 100000000,
			})
		}
		// 转换导出
		exps := make([]*Export, 0, len(rawtx.Exports))
		for _, exp := range rawtx.Exports {
			expaddr, err := exp.GetAddress()
			if err != nil {
				ctx.Error(err)
				return
			}
			expamount := exp.GetAmount().Uint64() / 100000000
			exps = append(exps, &Export{
				ExpAddress: expaddr.String(),
				ExpAmount:  expamount,
			})
		}
		rentx := new(Transfer)
		rentx.TxHash = txhash.String()
		rentx.TxTime = time.Unix(int64(rawtx.Timestamp.Uint64()), 0).Format("2006-01-02 15:04:05")
		rentx.TxImps = imps
		rentx.TxExps = exps
		renbk.TxList = append(renbk.TxList, rentx)
	}
	// 渲染页面
	ctx.HTML(200, "block.tpl", renbk)
	return
}
