package db

import (
	"log"
	"database/sql"

	"http-golang-api/types"

	_ "github.com/lib/pq"
)

func AddUser(user types.User) int {
	var id int

	// TODO: this url must be in .env file <17-10-23, modernpacifist> //
	db, err := sql.Open("postgres", "postgres://golanguser:golangpassword@localhost/golangdb?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.QueryRow("INSERT INTO users (name, age, salary, occupation) VALUES ($1, $2, $3, $4) RETURNING id", user.Name, user.Age, user.Salary, user.Occupation).Scan(&id)
	if err != nil {
		log.Println(err)
	}

	return id
}

func GetUser(userID string) types.User {
	// TODO: must check if the id exists in the first place in the db <17-10-23, modernpacifist> //
	var u types.User

	db, err := sql.Open("postgres", "postgres://golanguser:golangpassword@localhost/golangdb?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.QueryRow("SELECT * FROM users  WHERE id=$1", userID).Scan(&u.ID, &u.Name, &u.Age, &u.Salary, &u.Occupation)
	if err != nil {
		//return *new(types.User)
		//log.Fatal(err)
		log.Println(err)
	}

	return u
}
