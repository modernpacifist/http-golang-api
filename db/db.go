package db

import (
	"database/sql"
	"log"

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
		log.Println(err)
	}

	return u
}

func GetAllRecords() []types.User {
	// TODO: must check if the id exists in the first place in the db <17-10-23, modernpacifist> //
	var res []types.User

	db, err := sql.Open("postgres", "postgres://golanguser:golangpassword@localhost/golangdb?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var u types.User
		err := rows.Scan(&u.ID, &u.Name, &u.Age, &u.Salary, &u.Occupation)
		if err != nil {
			log.Println(err)
		}
		res = append(res, u)
	}

	return res
}
