package luas

func NewBlock() *Block {
	return &Block{}
}

type Block struct {
}

// 当前矿工的地址
func (b *Block) Coinbase() []byte {
	return nil
}

// 当前区块的高度
func (b *Block) Number() uint64 {
	return 0
}

// 得到区块的hash
func (b *Block) Blockhash(height uint64) []byte {
	return nil
}

// 当前区块时间
func (b *Block) Timestamp() uint64 {
	return 0
}
