package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Request struct {
	Recipient  string `json:"recipient"`
	Amount     string `json:"amount"`
	PrivateKey string `json:"privateKey"`
}

type Response struct {
	Message string `json:"message"`
	Balance string `json:"balance"`
}

func main() {
	http.HandleFunc("/sendEth", sendEthHandler)
	log.Println("Server is listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func sendEthHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client, err := ethclient.Dial("HTTP://127.0.0.1:7545")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	privateKey, err := crypto.HexToECDSA(req.PrivateKey)
	if err != nil {
		http.Error(w, "Invalid private key", http.StatusBadRequest)
		return
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		http.Error(w, "Cannot assert type: publicKey is not of type *ecdsa.PublicKey", http.StatusBadRequest)
		return
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	value := new(big.Int)
	value.SetString(req.Amount, 10)
	value = value.Mul(value, big.NewInt(1e18)) // Convert to Wei

	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	toAddress := common.HexToAddress(req.Recipient)
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID := big.NewInt(1337)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("Transaction sent: %s", signedTx.Hash().Hex())
	balance := getBalanceInEth(client, fromAddress)

	response := Response{
		Message: message,
		Balance: balance,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getBalanceInEth(client *ethclient.Client, address common.Address) string {
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return fmt.Sprintf("Error getting balance: %v", err)
	}
	ethValue := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))
	return fmt.Sprintf("Balance: %s ETH", ethValue.Text('f', 18))
}
