package blockchain

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "time"
)

type Borrower struct {
    Name           string
    IDNumber       string
    BusinessPermit string
    BusinessName   string
    EmailAddress   string
    PhoneNumber    string
    Location       string
}

type Lender struct {
    Name         string
    IDNumber     string
    EmailAddress string
    PhoneNumber  string
}

type Transaction struct {
    BorrowerID      string
    LenderID        string
    TransactionDate time.Time
    Amount          float64
}

type Block struct {
    Index        int
    Timestamp    time.Time
    Transaction  Transaction
    PrevHash     string
    Hash         string
}

type Blockchain struct {
    Blocks []Block
}

var Bc Blockchain

// InitBlockchain initializes the blockchain with the genesis block
func InitBlockchain() {
    genesisBlock := Block{}
    genesisBlock = Block{0, time.Now(), Transaction{}, "", calculateHash(genesisBlock)}
    Bc = Blockchain{[]Block{genesisBlock}}
}

// calculateHash computes the SHA-256 hash of a block
func calculateHash(block Block) string {
    record := string(block.Index) + block.Timestamp.String() + fmt.Sprintf("%v", block.Transaction) + block.PrevHash
    hash := sha256.New()
    hash.Write([]byte(record))
    hashed := hash.Sum(nil)
    return hex.EncodeToString(hashed)
}

// GenerateBlock creates a new block using the previous block and transaction data
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

// IsBlockValid validates a new block by checking the index and hashes
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
