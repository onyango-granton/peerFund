package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

type borrowersList struct {
	user *borrower
}

type loanList struct {
	loans *loan
	index int
}

type lendersList struct {
	user *lender
}

type loan struct {
	index        int
	lender       string
	borrower     string
	amount       int
	timestamp    string
	previousHash string
	hash         string
	prev         *loan
}

type borrower struct {
	fName            string
	lName            string
	IDNo             int
	businessPermitNo string
	phoneNumber      int
	email            string
	uid              string
	prev             *borrower
}

type lender struct {
	fName       string
	lName       string
	IDNo        int
	phoneNumber int
	email       string
	uid         string
	prev        *lender
}

func (b *loan)calculateHash() string {
	record := strconv.Itoa(b.index) + b.timestamp + strconv.Itoa(b.amount) + b.previousHash
	h := sha256.New()
	h.Write([]byte(record))
	//fmt.Println(b.timestamp)
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func (loans *loanList) generateloan(borrower, lender string, amount int) (*loan, error) {
	//var newloan loan
	newloan := &loan{}
	t := time.Now()

	newloan.index = loans.index
	newloan.timestamp = t.String()
	newloan.amount = amount
	newloan.borrower = borrower
	newloan.lender = lender
	newloan.hash = newloan.calculateHash()

	loans.index++

	if loans.loans == nil {
		loans.loans = newloan
	} else {
		current := loans.loans
		newloan.prev = current
		newloan.previousHash = current.hash
		loans.loans = newloan
	}

	return newloan, nil
}

func (lList *lendersList) regLender(fName, lname, email string, idNo, phoneNumber int) (*lender, error) {
	newL := &lender{}
	newL.fName = fName
	newL.lName = lname
	newL.email = email
	newL.phoneNumber = phoneNumber
	newL.IDNo = idNo
	newL.uid = generateUID(nil, newL)

	if lList.user == nil {
		lList.user = newL
	} else {
		current := lList.user
		newL.prev = current
		lList.user = newL
	}

	return newL, nil
}

func generateUID(userB *borrower, userL *lender) string {
	var record string
	if userB != nil {
		record = userB.fName + userB.lName + strconv.Itoa(userB.IDNo) + strconv.Itoa(userB.phoneNumber)
	} else {
		record = userL.fName + userL.lName + strconv.Itoa(userL.IDNo) + strconv.Itoa(userL.phoneNumber)
	}
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func (blist *borrowersList) regBorrower(fName, lname, businessPermitNo, email string, idNo, phoneNumber int) (*borrower, error) {
	newB := &borrower{}
	newB.fName = fName
	newB.lName = lname
	newB.businessPermitNo = businessPermitNo
	newB.email = email
	newB.phoneNumber = phoneNumber
	newB.IDNo = idNo
	newB.uid = generateUID(newB, nil)

	if blist.user == nil {
		blist.user = newB
	} else {
		current := blist.user
		newB.prev = current
		blist.user = newB
	}

	return newB, nil
}

func (loanNode *loan)isloanValid() bool {
	
	if loanNode.prev != nil {
		if loanNode.index != loanNode.prev.index+1 {
			
			return false
		}
		if loanNode.previousHash != loanNode.prev.hash {
			//fmt.Println("here")
			return false
		}
	}
	if loanNode.calculateHash() != loanNode.hash{
		return false
	}

	//fmt.Println("here")
	return true
}

func main() {
	//intantiate lenderslist... linked with unique id
	lendersList := &lendersList{}
	lenderAdd, _ := lendersList.regLender("sheila", "juma", "sheillajuma@gmail.com", 87878787, 7809099021)
	lenderAdd2, _ := lendersList.regLender("muthoni", "kamau", "kamaumuthoni@gmail.com", 8000098, 7809099021)
	//fmt.Println(lenderAdd)
	//fmt.Println(lenderAdd2)

	//intantiate borrowers list... linked with unique id
	borrowerList := &borrowersList{}
	borrowerAdd, _ := borrowerList.regBorrower("hamza", "nuru", "2011/234523", "hamzanuru@gmail.com", 89878889, 790554853)
	borrowerAdd2, _ := borrowerList.regBorrower("ziza", "bin", "2041/230000", "zizabin@gmail.com", 90020209, 720202022)
	//fmt.Println(borrowerAdd)
	//fmt.Println(borrowerAdd2)

	//intantiate loanlist... linked with unique id
	loanList := &loanList{}
	loanAdd, _ := loanList.generateloan(borrowerAdd.uid, lenderAdd.uid, 2500)
	loanAdd2, _ := loanList.generateloan(borrowerAdd2.uid, lenderAdd2.uid, 4500)
	if loanAdd2.isloanValid(){
		fmt.Println("success")
	}
	fmt.Println(loanAdd)
	fmt.Println(loanAdd2)

}
