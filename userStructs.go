package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

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
	prev *lender
}

func generateUID(userB *borrower,userL *lender) string {
	var record string
	if userB != nil{
		record = userB.fName + userB.lName + strconv.Itoa(userB.IDNo) + strconv.Itoa(userB.phoneNumber)
	} else {
		record = userL.fName + userL.lName + strconv.Itoa(userL.IDNo) + strconv.Itoa(userL.phoneNumber)
	}	
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func (newB *borrower) regBorrower(fName,lname,businessPermitNo,email string, idNo, phoneNumber int, prev *borrower) (*borrower, error) {
	newB.fName = fName
	newB.lName = lname
	newB.businessPermitNo = businessPermitNo
	newB.email = email
	newB.phoneNumber = phoneNumber
	newB.IDNo = idNo
	newB.prev = prev
	newB.uid = generateUID(newB,nil)

	return newB, nil
}

func (newL *lender) regLender(fName,lname,email string, idNo, phoneNumber int, prev *lender) (*lender, error) {
	newL.fName = fName
	newL.lName = lname
	newL.email = email
	newL.phoneNumber = phoneNumber
	newL.IDNo = idNo
	newL.prev = prev
	newL.uid = generateUID(nil,newL)

	return newL, nil
}


func ()