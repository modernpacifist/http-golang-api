package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"http-golang-api/types"

	_ "github.com/lib/pq"
)

type DatabaseManager struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func (db *DatabaseManager) Connect() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", db.Host, db.Port, db.User, db.Password, db.DBName)

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

func (db *DatabaseManager) AddUser(user types.User) int {
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

// func (db *DatabaseManager) GetUser(userID string) types.User {
func (db *DatabaseManager) GetUser(userID string) (*types.User, error) {
	// TODO: must check if the id exists in the first place in the db <17-10-23, modernpacifist> //
	var user types.User

	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("db.GetAllRecords: Could not connect to database %v", err)
	}

	err = conn.QueryRow("SELECT * FROM users  WHERE id=$1", userID).Scan(&user.ID, &user.Name, &user.Age, &user.Salary, &user.Occupation)
	if err != nil {
		log.Printf("db.GetUser.QueryRow: %v", err)
		//return nil,
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
	}

	log.Printf("db.GetUser: Successfully retrieved user with id %d\n", user.ID)
	return &user, nil
}

func (db *DatabaseManager) GetAllRecords() []types.User {
	var users []types.User

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
		users = append(users, u)
	}

	log.Printf("db.GetAllRecords.sql.Open: Successfully retrieved total %d records", len(users))

	return users
}
