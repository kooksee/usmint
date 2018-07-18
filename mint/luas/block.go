package luas

func NewBlock() *Block {
	return &Block{}
}

type Block struct {
}

// ﻿当前块的矿工的地址
func (b *Block) Coinbase() []byte {
	return nil
}

// ﻿当前块的数量
func (b *Block) Number() uint64 {
	return 0
}

// ﻿给定的块的hash值, 只有最近工作的256个块的hash值
func (b *Block) Blockhash(height uint64) []byte {
	return nil
}

//﻿当前块的时间戳
func (b *Block) Timestamp() uint64 {
	return 0
}


