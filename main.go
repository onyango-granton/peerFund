package main

import (
    "fmt"
    "log"
    "net/http"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "go-blockchain/badgerDb"
    "go-blockchain/blockchain"
)

func main() {
    // Initialize BadgerDB
    err := badgerDb.OpenDB("/tmp/badger")
    if err != nil {
        log.Fatal(err)
    }
    defer badgerDb.CloseDB()

    // Initialize blockchain with the genesis block
    blockchain.InitBlockchain()

    // Set up Gin web framework
    router := gin.Default()
    router.LoadHTMLGlob("templates/*")

    // Serve the form page at the root URL
    router.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", nil)
    })

    // Handle form submission
    router.POST("/submit", func(c *gin.Context) {
        // Parse form data for borrower
        borrower := blockchain.Borrower{
            Name:           c.PostForm("borrower_name"),
            IDNumber:       c.PostForm("borrower_id_number"),
            BusinessPermit: c.PostForm("borrower_business_permit"),
            BusinessName:   c.PostForm("borrower_business_name"),
            EmailAddress:   c.PostForm("borrower_email_address"),
            PhoneNumber:    c.PostForm("borrower_phone_number"),
            Location:       c.PostForm("borrower_location"),
        }

        // Parse form data for lender
        lender := blockchain.Lender{
            Name:         c.PostForm("lender_name"),
            IDNumber:     c.PostForm("lender_id_number"),
            EmailAddress: c.PostForm("lender_email_address"),
            PhoneNumber:  c.PostForm("lender_phone_number"),
        }

        // Parse and validate transaction amount
        amount, err := strconv.ParseFloat(c.PostForm("transaction_amount"), 64)
        if err != nil {
            c.String(http.StatusBadRequest, "Invalid amount")
            return
        }

        // Create a new transaction
        transaction := blockchain.Transaction{
            BorrowerID:      borrower.IDNumber,
            LenderID:        lender.IDNumber,
            TransactionDate: time.Now(),
            Amount:          amount,
        }

        // Generate a new block and add it to the blockchain
        newBlock := blockchain.GenerateBlock(blockchain.Bc.Blocks[len(blockchain.Bc.Blocks)-1], transaction)
        if blockchain.IsBlockValid(newBlock, blockchain.Bc.Blocks[len(blockchain.Bc.Blocks)-1]) {
            blockchain.Bc.Blocks = append(blockchain.Bc.Blocks, newBlock)

            // Store the data in BadgerDB
            err := badgerDb.StoreData("borrower_"+borrower.IDNumber, fmt.Sprintf("%+v", borrower))
            if err != nil {
                c.String(http.StatusInternalServerError, err.Error())
                return
            }

            err = badgerDb.StoreData("lender_"+lender.IDNumber, fmt.Sprintf("%+v", lender))
            if err != nil {
                c.String(http.StatusInternalServerError, err.Error())
                return
            }

            err = badgerDb.StoreData("transaction_"+borrower.IDNumber+"_"+lender.IDNumber, fmt.Sprintf("%+v", transaction))
            if err != nil {
                c.String(http.StatusInternalServerError, err.Error())
                return
            }

            err = badgerDb.StoreData("blockchain", fmt.Sprintf("%+v", blockchain.Bc))
            if err != nil {
                c.String(http.StatusInternalServerError, err.Error())
                return
            }

            // Respond to the client
            c.String(http.StatusOK, "Data stored successfully")
        } else {
            c.String(http.StatusInternalServerError, "Failed to validate the new block")
        }
    })

    // Start the web server on port 8080
    router.Run(":8080")
}
