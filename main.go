package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"text/template"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	jwtKey      = []byte("my_secret_key")
	credentials = map[string]struct {
		Password string
		Name     string
		Phone    string
		Email    string
		Address  string
	}{
		"kada@peerfund.com": {"xkUbdz6r", "Kennedy Ada", "0704513552", "adakennedy@outlook.com", "0xcef41520D00132677de7cFC89956B212169109C4"},
		"ann@peerfund.com":           {"bSkinGurl", "Ann Maina", "0724318117", "nyagoh@gmail.com", "0x803b88327972D9ad11170152E0A826Fe3B0BF469"},
		"josotieno@peerfund.com":     {"fyaman42", "Joseph Otieno", "0722549387", "jokumu25@gmail.com", "0x3e743E3d2728513f065d56Eead3A4Fd657966D64"},
	}
	// templates = template.Must(template.ParseGlob("templates/*.html"))
)

type Claims struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Address  string `json:"address"`
	jwt.StandardClaims
}

type Request struct {
	Recipient  string `json:"recipient"`
	Amount     string `json:"amount"`
	PrivateKey string `json:"privateKey"`
}

type Response struct {
	Message string `json:"message"`
	Balance string `json:"balance"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var username, password string
	if r.URL.Path != "/login" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if r.Method == http.MethodGet {
		t, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		username = r.FormValue("email")
		password = r.FormValue("password")

		user, ok := credentials[username]
		if !ok || user.Password != password {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		expirationTime := time.Now().Add(5 * time.Minute)
		claims := &Claims{
			Username: username,
			Name:     user.Name,
			Phone:    user.Phone,
			Email:    user.Email,
			Address:  user.Address,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		// Create and sign the token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Set the token as a cookie
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})

		// Redirect to the dashboard
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if r.Method == http.MethodGet {
		t, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
	} else {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	tokenStr := cookie.Value
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if !token.Valid {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Connect to Ganache and get balance
	client, err := ethclient.Dial("HTTP://127.0.0.1:7545")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	formAddress := common.HexToAddress(claims.Address)
	balance := getBalanceInEth(client, formAddress)

	data := struct {
		Username string
		Name     string
		Phone    string
		Email    string
		Address  string
		Balance  string
	}{
		Username: claims.Username,
		Name:     claims.Name,
		Phone:    claims.Phone,
		Email:    claims.Email,
		Address:  claims.Address,
		Balance:  balance,
	}

	t, err := template.ParseFiles("templates/dashboard.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

func sendEthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse form values
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// Retrieve values from form
		amount := r.FormValue("amount")
		recipient := r.FormValue("address")
		privateKeyHex := r.FormValue("key")

		// Create a new Ethereum client
		client, err := ethclient.Dial("HTTP://127.0.0.1:7545")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Convert private key from hex to ECDSA
		privateKey, err := crypto.HexToECDSA(privateKeyHex)
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

		// Fetch nonce, gas, and create the transaction
		nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		value := new(big.Int)
		value.SetString(amount, 10)
		value = value.Mul(value, big.NewInt(1e18)) // Convert to Wei

		gasLimit := uint64(21000)
		gasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		toAddress := common.HexToAddress(recipient)
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

		// Fetch balances
		balanceSender := getBalanceInEth(client, fromAddress)
		balanceRecipient := getBalanceInEth(client, toAddress)

		// Create response message
		response := map[string]string{
			"message":          "Transaction was successful!",
			"balanceSender":    balanceSender,
			"balanceRecipient": balanceRecipient,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func getBalanceInEth(client *ethclient.Client, address common.Address) string {
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return fmt.Sprintf("Error getting balance: %v", err)
	}
	ethValue := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))
	return fmt.Sprintf("Balance: %s ETH", ethValue.Text('f', 18))
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
