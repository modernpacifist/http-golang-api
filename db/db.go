package db

import (
	"log"
	"database/sql"
	_ "github.com/lib/pq"
)

func AddUser() int {
	var id int

	db, err := sql.Open("postgres", "postgres://golanguser:golangpassword@localhost/golangdb?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.QueryRow("INSERT INTO users (name, age, salary, occupation) VALUES ($1, $2, $3, $4) RETURNING id", "user1", "21", "10001", "programmer").Scan(&id)
	if err != nil {
		log.Fatal(err)
	}

	return id
}

func GetUser(userID int) string {
	var res string

	db, err := sql.Open("postgres", "postgres://golanguser:golangpassword@localhost/golangdb?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.QueryRow("SELECT FROM users (name, age, salary, occupation) where id=$1", userID).Scan(&res)
	if err != nil {
		log.Fatal(err)
	}

	return res
}
