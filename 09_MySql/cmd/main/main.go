package main

import (
	"database/sql"
	"fmt"
)

func main() {
	dsn := "user:password@tcp(127.0.0.1:3306)/testdb"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// CREATE
	id, _ := createUser(db, "Alice", "alice@example.com")
	fmt.Println("Inserted user ID:", id)

	// READ ONE
	user, _ := getUser(db, int(id))
	fmt.Println("User:", user)

	// UPDATE
	_ = updateUser(db, int(id), "Alice Updated", "alice@newmail.com")

	// READ ALL
	users, _ := getAllUsers(db)
	fmt.Println("All users:", users)

	// DELETE
	_ = deleteUser(db, int(id))
	fmt.Println("User deleted.")
}

type User struct {
	ID    int
	Name  string
	Email string
}

func getUser(db *sql.DB, id int) (User, error) {
	var user User
	row := db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id)
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	return user, err
}

func createUser(db *sql.DB, name, email string) (int64, error) {
	result, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", name, email)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func getAllUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func updateUser(db *sql.DB, id int, name, email string) error {
	_, err := db.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?", name, email, id)
	return err
}

func deleteUser(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}
