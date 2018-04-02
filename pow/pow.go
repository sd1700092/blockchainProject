package main

import (
	"math"
	"math/big"
)

type Block struct {
	Timestamp int64
	Data []byte
	PrevBlockHash []byte
	Hash []byte
	Nonce int
}

type ProofOfWork struct {
	block *Block
	target *big.Int
}

var maxNonce = math.MaxInt64

const targetBits = 24

func NewBlock(data string, prevBlockHash []byte) *Block {

}

func main() {
	
}
