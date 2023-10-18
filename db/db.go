package db

import (
	"database/sql"
	"log"
	"fmt"

	"http-golang-api/types"

	_ "github.com/lib/pq"
)

//// Define an interface for database operations
//type DBOperations interface {
//Connect() (*sql.DB, error)
//AddUser(user User) error
//GetUser
//}

type DatabaseManager struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// func Connect() (*sql.DB, error) {
func (db DatabaseManager) Connect() (*sql.DB, error) {
	//connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", db.Host, db.Port, db.User, db.Password, db.DBName)
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", db.User, db.Password, db.Host, db.DBName)

	//db, err := sql.Open("postgres", "postgres://golanguser:golangpassword@localhost/golangdb?sslmode=disable")
	conn, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (db DatabaseManager) AddUser(user types.User) int {
//func AddUser(user types.User) int {
	var id int

	// TODO: this url must be in .env file <17-10-23, modernpacifist> //
	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("db.AddUser: Could not connect to database %v", err)
	}

	err = conn.QueryRow("INSERT INTO users (name, age, salary, occupation) VALUES ($1, $2, $3, $4) RETURNING id", user.Name, user.Age, user.Salary, user.Occupation).Scan(&id)
	if err != nil {
		log.Printf("db.AddUser.QueryRow: %v", err)
	}

	log.Printf("db.AddUser: Successfully added user with id %d\n", id)

	return id
}

func (db DatabaseManager) GetUser(userID string) types.User {
	// TODO: must check if the id exists in the first place in the db <17-10-23, modernpacifist> //
	var u types.User

	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("db.GetAllRecords: Could not connect to database %v", err)
	}

	err = conn.QueryRow("SELECT * FROM users  WHERE id=$1", userID).Scan(&u.ID, &u.Name, &u.Age, &u.Salary, &u.Occupation)
	if err != nil {
		log.Printf("db.GetUser.QueryRow: %v", err)
	}

	log.Printf("db.GetUser: Successfully retrieved user with id %d\n", u.ID)
	return u
}

func (db DatabaseManager) GetAllRecords() []types.User {
	var res []types.User

	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("db.GetAllRecords: Could not connect to database:%v", err)
	}

	rows, err := conn.Query("SELECT * FROM users")
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
