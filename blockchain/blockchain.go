package blockchain

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "time"
)

// Borrower struct holds details about the borrower
type Borrower struct {
    Name           string
    IDNumber       string
    BusinessPermit string
    BusinessName   string
    EmailAddress   string
    PhoneNumber    string
    Location       string
}

// Lender struct holds details about the lender
type Lender struct {
    Name         string
    IDNumber     string
    EmailAddress string
    PhoneNumber  string
}

// Transaction struct holds details about the transaction
type Transaction struct {
    BorrowerID      string
    LenderID        string
    TransactionDate time.Time
    Amount          float64
}

// Block struct represents a single block in the blockchain
type Block struct {
    Index        int
    Timestamp    time.Time
    Transaction  Transaction
    PrevHash     string
    Hash         string
}

// Blockchain struct represents the entire blockchain
type Blockchain struct {
    Blocks []Block
}

var Bc Blockchain

func calculateHash(block Block) string {
    record := fmt.Sprintf("%d%s%v%s", block.Index, block.Timestamp.String(), block.Transaction, block.PrevHash)
    hash := sha256.New()
    hash.Write([]byte(record))
    hashed := hash.Sum(nil)
    return hex.EncodeToString(hashed)
}

func GenerateBlock(oldBlock Block, transaction Transaction) Block {
    newBlock := Block{
        Index:       oldBlock.Index + 1,
        Timestamp:   time.Now(),
        Transaction: transaction,
        PrevHash:    oldBlock.Hash,
        Hash:        "",
    }
    newBlock.Hash = calculateHash(newBlock)
    return newBlock
}

func IsBlockValid(newBlock, oldBlock Block) bool {
    if oldBlock.Index+1 != newBlock.Index {
        return false
    }
    if oldBlock.Hash != newBlock.PrevHash {
        return false
    }
    if calculateHash(newBlock) != newBlock.Hash {
        return false
    }
    return true
}

func InitBlockchain() {
    genesisBlock := Block{
        Index:       0,
        Timestamp:   time.Now(),
        Transaction: Transaction{},
        PrevHash:    "",
        Hash:        calculateHash(Block{}),
    }
    Bc = Blockchain{[]Block{genesisBlock}}
}
