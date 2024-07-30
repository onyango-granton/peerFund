package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
    // Open database connection
    var err error
    dsn := "myuser:mypassword@tcp(127.0.0.1:3306)/mydatabase"
    db, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Ping the database to verify connection
    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }

    // Set up HTTP routes
    http.HandleFunc("/submit", submitHandler)
    http.HandleFunc("/retrieve", retrieveHandler)

    // Start the server
    fmt.Println("Server started at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    // Read the form data
    r.ParseForm()
    data := r.FormValue("data")

    // Insert data into the database
    _, err := db.Exec("INSERT INTO submissions (data) VALUES (?)", data)
    if err != nil {
        http.Error(w, "Failed to store data", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Data stored successfully"))
}

func retrieveHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    // Retrieve data from the database
    rows, err := db.Query("SELECT data FROM submissions")
    if err != nil {
        http.Error(w, "Failed to retrieve data", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var result string
    for rows.Next() {
        var data string
        if err := rows.Scan(&data); err != nil {
            http.Error(w, "Failed to scan data", http.StatusInternalServerError)
            return
        }
        result += data + "\n"
    }

    w.Header().Set("Content-Type", "text/plain")
    w.Write([]byte(result))
}
