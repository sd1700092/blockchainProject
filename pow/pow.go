package main

import (
	"math"
	"math/big"
	"time"
	"fmt"
	"bytes"
	"encoding/binary"
	"log"
	"crypto/sha256"
	"strconv"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

var maxNonce = math.MaxInt64

const targetBits = 24

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits)) // 左移232位，这样一个256位的数的前24位就都是0了
	pow := &ProofOfWork{b, target}
	return pow
}

// POW验证机制的实现
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	fmt.Printf("Mining the block containing \" %s \"\n", pow.block.Data)
	for nonce < maxNonce {
		data := pow.prepareData(nonce) //1.准备数据
		hash = sha256.Sum256(data)     //2.用SHA-256对数据进行哈希
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:]) //3.将哈希转换成一个大整数
		if hashInt.Cmp(pow.target) == -1 { //4.如果比目标工作量证明的值要小，那么就break掉，代表挖到矿了
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")
	return nonce, hash[:]
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	isValid := hashInt.Cmp(pow.target) == -1
	return isValid
}

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func main() {
	data1 := []byte("I like donuts")
	data2 := []byte("I like donutsca07ca")
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	fmt.Printf("%x\n", sha256.Sum256(data1))
	fmt.Printf("%64x\n", target)
	fmt.Printf("%x\n", sha256.Sum256(data2))

	bc := NewBlockchain()
	bc.AddBlock("Send 1 BTC to 土拨鼠")
	bc.AddBlock("Send 2 more BTC to 土拨鼠")
	for _, block := range bc.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("POW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}

type Blockchain struct {
	blocks []*Block
}

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}
