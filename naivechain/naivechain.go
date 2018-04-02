package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
	"fmt"
)

type Block struct {
	Index     int
	Timestamp string
	Data      string
	Hash      string
	PrevHash  string
}

var Blockchain []*Block

func calculateHash(block *Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + block.Data + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func generateBlock(oldBlock *Block, data string) *Block {
	newBlock := &Block{}
	t := time.Now()
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Data = data
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)
	return newBlock
}

func main() {
	//block := Block{0, "2018-04-02 15:47:43", "aaaa", "0", "0"}
	//fmt.Println(calculateHash(&block))
	t := time.Now()
	datas := []string{"btc", "etc", "key", "eos", "bihu", "xiequan.info"}
	for k, v := range datas {
		if k>=1{
			newBlock := generateBlock(Blockchain[len(Blockchain)-1], v)
			Blockchain = append(Blockchain, newBlock)
		} else {
			genesisBlock := &Block{0, t.String(), "创世快", calculateHash(&Block{}), ""}
			Blockchain = append(Blockchain, genesisBlock)
		}
	}
	for _,v := range Blockchain {
		fmt.Printf("index: %d\n", v.Index)
		fmt.Printf("data: %s\n", v.Data)
		fmt.Printf("timestamp: %s\n", v.Timestamp)
		fmt.Printf("hash: %s\n", v.Hash)
		fmt.Printf("prevHash: %s\n", v.PrevHash)
	}
}
