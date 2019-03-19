package coinbook

type Import struct {
	TxHash     string
	TxGenesis  bool
	ImpAddress string
	ImpAmount  uint64
}

type Export struct {
	ExpAddress string
	ExpAmount  uint64
}

type Transfer struct {
	TxHash string
	TxTime string
	TxImps []*Import
	TxExps []*Export
}

type Header struct {
	BkHash       string
	BkHeight     uint64
	BkVersion    uint16
	BkMerkleRoot string
	BkTime       string
	BkSize       uint64
	TxCount      uint64
	Confirm      uint64
}

type Block struct {
	BkHeader *Header
	BkPrev   *Header
	BkNext   *Header
	TxList   []*Transfer
}
