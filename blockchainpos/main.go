package main

import (
	"sync"
	"crypto/sha256"
	"encoding/hex"
	"time"
	"net"
	"io"
	"bufio"
	"strconv"
	"log"
	"fmt"
	"encoding/json"
)

type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
	Validator string
}

var Blockchain []Block
var tempBlocks []Block

var candidateBlocks = make(chan Block)
var announcements = make(chan string)

var mutex = &sync.Mutex{}

var validators = make(map[string]int)

func calculateHash(s string ) string {
	h:=sha256.New()
	h.Write([]byte(s))
	hashed:=h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func calculateBlockHash(block Block) string{
	record:=string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash
	return calculateHash(record)
}

func generateBlock(oldBlock Block, BPM int, address string) (Block, error) {
	var newBlock Block
	t:=time.Now()
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateBlockHash(newBlock)
	newBlock.Validator = address
	return newBlock, nil
}

func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index + 1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash!=newBlock.PrevHash{
		return false
	}
	if calculateBlockHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	go func() {
		for {
			msg:=<-announcements
			io.WriteString(conn, msg)
		}
	}()
	var address string
	io.WriteString(conn, "Enter token balance:")
	scanBalance:=bufio.NewScanner(conn)
	for scanBalance.Scan() {
		balance, err := strconv.Atoi(scanBalance.Text())
		if err!=nil {
			log.Printf("%v not a number: %v", scanBalance.Text(), err)
			return
		}
		t:=time.Now()
		address = calculateHash(t.String())
		validators[address] = balance
		fmt.Println("validators: ", validators)
		break
	}
	io.WriteString(conn, "\nEnter a new BPM:")

	scanBPM:=bufio.NewScanner(conn)

	go func() {
		for {
			for scanBPM.Scan() {
				bpm, err:=strconv.Atoi(scanBPM.Text())
				if err!=nil {
					log.Printf("%v not a number: %v", scanBPM.Text(), err)
					delete(validators, address)
					conn.Close()
				}
				mutex.Lock()
				oldLastIndex:=Blockchain[len(Blockchain) - 1]
				mutex.Unlock()

				newBlock, err:=generateBlock(oldLastIndex, bpm, address)
				if err!=nil {
					log.Println(err)
					continue
				}
				if isBlockValid(newBlock, oldLastIndex) {
					candidateBlocks <- newBlock
				}
				io.WriteString(conn, "\nEnter a new BPM:")
			}
		}
	}()

	for {
		time.Sleep(time.Minute)
		mutex.Lock()
		output,err:=json.Marshal(Blockchain)
		mutex.Unlock()
		if err!=nil{
			log.Fatal(err)
		}
		io.WriteString(conn, string(output) + "\n")
	}
}

func pickWinner(){
	time.Sleep(30 * time.Second)
	mutex.Lock()
	temp:=tempBlocks
	mutex.Unlock()
	lotteryPool:=[]string{}
	if len(temp) > 0 {

	}
}

func main() {
	
}
