package main

import (
	"fmt"
	"go/postgresql-demo/db"
	"go/postgresql-demo/router"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type User struct {
	Id   uint32
	Name string
}

// CORS middleware
func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins (use "*" for any domain)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("failed to load .env fil", err)
	}

	// Connection string
	port := os.Getenv("PORT")
	connectString := os.Getenv("CONNECTION_STRING")
	fmt.Println("Connection String : ", connectString)

	db := db.NewPostgreSqlDbConnection(connectString)
	fmt.Println("successfully connected to postgres db.", db.Db)

	router := router.Router(db.Db)
	fmt.Println("all routes are successfully set.")

	fmt.Println("starting server on port :8080")

	log.Fatalln(http.ListenAndServe(port, withCORS(router)))
}

// func main() {
// 	_ = godotenv.Load()

// 	// Connection string
// 	connectString := os.Getenv("CONNECTION_STRING")
// 	fmt.Println("Connection String : ", connectString)

// 	// Open connection
// 	db, err := sql.Open("postgres", connectString)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	// Verify connection
// 	err = db.Ping()
// 	if err != nil {
// 		log.Fatal("Could not connect: ", err)
// 	}

// 	fmt.Println("✅ Connected to PostgreSQL!")

// 	// user1 := User{
// 	// 	Id:   1,
// 	// 	Name: "Droga",
// 	// }
// 	// user2 := User{
// 	// 	Id:   4,
// 	// 	Name: "KK",
// 	// }

// 	// id1 := insertUser(db, user1)
// 	// id2 := insertUser(db, user2)

// 	// println(id1)
// 	// println(id2)

// 	users := fectchAllUsers(db)

// 	for _, user := range users {
// 		fmt.Println(user.Id, user.Name)
// 	}

// }

// func fectchAllUsers(db *sql.DB) []*User {
// 	qurey := `SELECT * FROM users`
// 	cursor, err := db.Query(qurey)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer cursor.Close()

// 	var users []*User

// 	for cursor.Next() {
// 		var user User
// 		err := cursor.Scan(&user.Id, &user.Name)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		users = append(users, &user)
// 	}

// 	return users
// }

// func insertUser(db *sql.DB, user User) uint32 {
// 	query := `INSERT INTO users (id,name) values ($1,$2) RETURNING id`

// 	var pk uint32
// 	err := db.QueryRow(query, user.Id, user.Name).Scan(&pk)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return pk
// }
