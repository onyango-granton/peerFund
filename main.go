package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var client *ethclient.Client

type TransferRequest struct {
	ToAddress string `json:"toAddress"`
	Amount    string `json:"amount"`
}

type TransferResponse struct {
	FromAddress string `json:"fromAddress"`
	ToAddress   string `json:"toAddress"`
	Amount      string `json:"amount"`
	Balance     string `json:"balance"`
	Message     string `json:"message"`
}

func main() {
	var err error
	client, err = connectToEthereum()
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	http.HandleFunc("/transfer", transferHandler)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	fmt.Println("Server started at http://localhost/:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func connectToEthereum() (*ethclient.Client, error) {
	client, err := ethclient.Dial("HTTP://127.0.0.1:7545")
	if err != nil {
		return nil, err
	}
	fmt.Println("We have a connection")
	return client, nil
}

func transferHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	toAddress := req.ToAddress
	amount := req.Amount

	response, err := transferETH(toAddress, amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func checkBalance(client *ethclient.Client, addressStr string) (*big.Int, error) {
	address := common.HexToAddress(addressStr)
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func transferETH(toAddress string, amount string) (*TransferResponse, error) {
	privateKey, err := crypto.HexToECDSA("4237445a301a6e009ce66f468745aeca78dfaab83f6b451c93ad92586392d940")
	if err != nil {
		return nil, fmt.Errorf("error loading private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, fmt.Errorf("error getting nonce: %v", err)
	}

	value := new(big.Int)
	value, ok = value.SetString(amount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid amount")
	}
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting suggested gas price: %v", err)
	}

	toAddressHex := common.HexToAddress(toAddress)
	var data []byte
	tx := types.NewTransaction(nonce, toAddressHex, value, gasLimit, gasPrice, data)

	chainID := big.NewInt(1337)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return nil, fmt.Errorf("error signing transaction: %v", err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, fmt.Errorf("error sending transaction: %v", err)
	}

	receipt, err := waitForTx(client, signedTx.Hash())
	if err != nil {
		return nil, fmt.Errorf("error waiting for transaction: %v", err)
	}

	message := "Transaction failed"
	if receipt.Status == types.ReceiptStatusSuccessful {
		message = "Transaction successful"
	}

	balance, err := checkBalance(client, fromAddress.Hex())
	if err != nil {
		return nil, fmt.Errorf("error checking balance: %v", err)
	}

	response := &TransferResponse{
		FromAddress: fromAddress.Hex(),
		ToAddress:   toAddressHex.Hex(),
		Amount:      value.String(),
		Balance:     balance.String(),
		Message:     message,
	}

	return response, nil
}

func waitForTx(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if err == nil {
			return receipt, nil
		}
		if err != nil && err != ethereum.NotFound {
			return nil, err
		}
		time.Sleep(time.Second)
	}
}
