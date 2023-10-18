package db

import (
	"database/sql"
	"log"

	"http-golang-api/types"

	_ "github.com/lib/pq"
)

func dbConnect() (*sql.DB, error) {
	db, err := sql.Open("postgres", "postgres://golanguser:golangpassword@localhost/golangdb?sslmode=disable")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func AddUser(user types.User) int {
	var id int

	// TODO: this url must be in .env file <17-10-23, modernpacifist> //
	db, err := dbConnect()
	if err != nil {
		log.Fatalf("db.GetAllRecords: Could not connect to database %v", err)
	}

	err = db.QueryRow("INSERT INTO users (name, age, salary, occupation) VALUES ($1, $2, $3, $4) RETURNING id", user.Name, user.Age, user.Salary, user.Occupation).Scan(&id)
	if err != nil {
		log.Printf("db.AddUser.QueryRow: %v", err)
	}

	log.Printf("db.AddUser: Successfully added user with id %d\n", id)

	return id
}

func GetUser(userID string) types.User {
	// TODO: must check if the id exists in the first place in the db <17-10-23, modernpacifist> //
	var u types.User

	db, err := dbConnect()
	if err != nil {
		log.Fatalf("db.GetAllRecords: Could not connect to database %v", err)
	}

	err = db.QueryRow("SELECT * FROM users  WHERE id=$1", userID).Scan(&u.ID, &u.Name, &u.Age, &u.Salary, &u.Occupation)
	if err != nil {
		log.Printf("db.GetUser.QueryRow: %v", err)
	}

	log.Printf("db.GetUser: Successfully retrieved user with id %d\n", u.ID)
	return u
}

func GetAllRecords() []types.User {
	var res []types.User

	db, err := dbConnect()
	if err != nil {
		log.Fatalf("db.GetAllRecords: Could not connect to database:%v", err)
	}

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Printf("db.GetAllRecord.Querys: Could not select from database:%v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u types.User
		err := rows.Scan(&u.ID, &u.Name, &u.Age, &u.Salary, &u.Occupation)
		if err != nil {
			log.Printf("db.GetAllRecords.rows.Scan: could not get info for user with ID:%d\n", u.ID)
		}
		res = append(res, u)
	}

	log.Printf("db.GetAllRecords.sql.Open: Successfully retrieved total %d records", len(res))

	return res
}
