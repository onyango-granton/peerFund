package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

type Block struct {
	index        int
	lender       string
	borrower     string
	amount       int
	timestamp    string
	previousHash string
	hash         string
	prev *Block
}

func calculateHash(b Block) string {
	record := strconv.Itoa(b.index) + b.timestamp + strconv.Itoa(b.amount) + b.previousHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func (newBlock *Block) generateBlock(oldBlock *Block,  borrower, lender string, amount int) (*Block, error) {
	//var newBlock Block
	t := time.Now()

	newBlock.index = oldBlock.index + 1
	newBlock.timestamp = t.String()
	newBlock.previousHash = oldBlock.hash
	newBlock.amount = amount
	newBlock.borrower = borrower
	newBlock.lender = lender
	newBlock.prev = oldBlock
	newBlock.hash = calculateHash(*newBlock)
	

	return newBlock, nil
}

func isBlockValid(newBlock, oldBlock *Block) bool {
	if oldBlock.index + 1 != newBlock.index{
		return false
	}
	if oldBlock.hash != newBlock.previousHash{
		return false
	}
	if calculateHash(*newBlock) != newBlock.hash{
		return false
	}
	if newBlock.prev != oldBlock{
		return false
	}
	return true
}



func main(){
	//
	//zeroBlock := &Block{}
	zeroBlock := &Block{index: 0,lender: "hamza",borrower: "nurr", amount: 2500, timestamp: time.Now().String(),previousHash: "", prev: nil}
	
	zeroBlock.hash = calculateHash(*zeroBlock)
	//fmt.Println(zeroBlock)

	//blockOne,_ := generateBlock(*zeroBlock, "sheila","juma",4000)

	newBlock := &Block{}
	newBlock,_ = newBlock.generateBlock(zeroBlock,"sheila","juma",4000)
	//fmt.Println(newBlock.prev)

	if newBlock.prev != zeroBlock{
		fmt.Println("Here")
	}

	if isBlockValid(newBlock,zeroBlock){
		fmt.Println(isBlockValid(newBlock,zeroBlock))
	}

}