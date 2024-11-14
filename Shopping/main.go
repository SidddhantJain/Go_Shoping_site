package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"

    _ "github.com/go-sql-driver/mysql"
    "github.com/gorilla/mux"
    "github.com/gorilla/sessions"
    "golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var store = sessions.NewCookieStore([]byte("super-secret-key"))

func initDB() {
    dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
    var err error
    db, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }
}

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Password string `json:"password"`
    Role     string `json:"role"`
}

type Product struct {
    ID          int     `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    ImageURL    string  `json:"image_url"`
}

func handleSignup(w http.ResponseWriter, r *http.Request) {
    var user User
    json.NewDecoder(r.Body).Decode(&user)
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

    stmt, err := db.Prepare("INSERT INTO users(username, password, role) VALUES(?, ?, 'user')")
    if err != nil {
        log.Fatal(err)
    }

    _, err = stmt.Exec(user.Username, string(hashedPassword))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"success": "User registered successfully"})
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
    var user User
    json.NewDecoder(r.Body).Decode(&user)

    stmt := "SELECT id, password, role FROM users WHERE username=?"
    row := db.QueryRow(stmt, user.Username)

    var dbUser User
    err := row.Scan(&dbUser.ID, &dbUser.Password, &dbUser.Role)
    if err != nil {
        http.Error(w, "Invalid login", http.StatusUnauthorized)
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
    if err != nil {
        http.Error(w, "Invalid login", http.StatusUnauthorized)
        return
    }

    session, _ := store.Get(r, "session")
    session.Values["authenticated"] = true
    session.Values["role"] = dbUser.Role
    session.Save(r, w)

    json.NewEncoder(w).Encode(map[string]string{"success": "Login successful", "role": dbUser.Role})
}

func handleAddProduct(w http.ResponseWriter, r *http.Request) {
    var product Product
    json.NewDecoder(r.Body).Decode(&product)

    stmt, err := db.Prepare("INSERT INTO products(name, description, price, image_url) VALUES(?, ?, ?, ?)")
    if err != nil {
        log.Fatal(err)
    }

    _, err = stmt.Exec(product.Name, product.Description, product.Price, product.ImageURL)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode("Product added successfully")
}

func handleAdminDashboard(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "session")
    if session.Values["role"] != "admin" {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }

    // Admin functionality for managing products
}

func main() {
    initDB()

    router := mux.NewRouter()
    router.HandleFunc("/signup", handleSignup).Methods("POST")
    router.HandleFunc("/login", handleLogin).Methods("POST")
    router.HandleFunc("/product", handleAddProduct).Methods("POST")
    router.HandleFunc("/admin", handleAdminDashboard).Methods("GET")

    fmt.Println("Server running on port 8080")
    http.ListenAndServe(":8080", router)
}
