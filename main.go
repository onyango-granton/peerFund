package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
	"strconv"
    "github.com/dgraph-io/badger/v3"
    "github.com/gin-gonic/gin"
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

var db *badger.DB

func main() {
    var err error
    db, err = badger.Open(badger.DefaultOptions("/tmp/badger"))
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    router := gin.Default()
    router.LoadHTMLGlob("templates/*")

    router.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", nil)
    })

    router.POST("/submit", func(c *gin.Context) {
        borrower := Borrower{
            Name:           c.PostForm("borrower_name"),
            IDNumber:       c.PostForm("borrower_id_number"),
            BusinessPermit: c.PostForm("borrower_business_permit"),
            BusinessName:   c.PostForm("borrower_business_name"),
            EmailAddress:   c.PostForm("borrower_email_address"),
            PhoneNumber:    c.PostForm("borrower_phone_number"),
            Location:       c.PostForm("borrower_location"),
        }

        lender := Lender{
            Name:         c.PostForm("lender_name"),
            IDNumber:     c.PostForm("lender_id_number"),
            EmailAddress: c.PostForm("lender_email_address"),
            PhoneNumber:  c.PostForm("lender_phone_number"),
        }

        amount, err := strconv.ParseFloat(c.PostForm("transaction_amount"), 64)
        if err != nil {
            c.String(http.StatusBadRequest, "Invalid amount")
            return
        }

        transaction := Transaction{
            BorrowerID:      borrower.IDNumber,
            LenderID:        lender.IDNumber,
            TransactionDate: time.Now(),
            Amount:          amount,
        }

        err = db.Update(func(txn *badger.Txn) error {
            // Store borrower
            err := txn.Set([]byte("borrower_"+borrower.IDNumber), []byte(fmt.Sprintf("%+v", borrower)))
            if err != nil {
                return err
            }

            // Store lender
            err = txn.Set([]byte("lender_"+lender.IDNumber), []byte(fmt.Sprintf("%+v", lender)))
            if err != nil {
                return err
            }

            // Store transaction
            err = txn.Set([]byte("transaction_"+borrower.IDNumber+"_"+lender.IDNumber), []byte(fmt.Sprintf("%+v", transaction)))
            if err != nil {
                return err
            }

            return nil
        })
        if err != nil {
            c.String(http.StatusInternalServerError, err.Error())
            return
        }

        c.String(http.StatusOK, "Data stored successfully")
    })

    router.Run(":8080")
}
