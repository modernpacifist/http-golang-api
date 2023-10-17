package main

import (
	"fmt"
	"net/http"
	"log"

	"http-golang-api/db"

	"github.com/swaggo/http-swagger"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
)

type User struct {
	ID int `json:"id"`
	Name int `json:"name"`
	Age int `json:"age"`
	Salary int `json:"salary"`
	Occupation int `json:"occupation"`
}

func addUserHandler(w http.ResponseWriter, req *http.Request) {
	user := User{
		ID: "1", 
		Name: "john"
		Age: "age"
		Salary: "salary"
		Occupation: "occupation"
	}
}

func getUserHandler(w http.ResponseWriter, req *http.Request) {

}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func main() {
	db.DbConnect()

	http.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/adduser/", addUserHandler)
	http.HandleFunc("/getuser/{id}", getUserHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
