package db

import (
	"log"
	"database/sql"
	_ "github.com/lib/pq"
)

func DbAddUser() int {
	var id int

	db, err := sql.Open("postgres", "postgres://golanguser:golangpassword@localhost/golangdb?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//stmt, err := db.Prepare("INSERT INTO users (name, age, salary, occupation) VALUES ($1, $2, $3, $4) RETURNING id").Scan(&id)
	//if err != nil {
		//log.Fatal(err)
	//}
	//defer stmt.Close()

	//_, err = stmt.Exec("user1", "21", "10001", "programmer")
	//if err != nil {
		//log.Fatal(err)
	//}
	err := db.QueryRow("INSERT INTO users (name, age, salary, occupation) VALUES ($1, $2, $3, $4) RETURNING id", "user1", "21", "10001", "programmer").Scan(&id)
	if err != nil {
		log.Fatal(err)
	}

	return id
}
