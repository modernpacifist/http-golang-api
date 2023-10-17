package main

import (
	"fmt"
	"net/http"
	"log"
	"encoding/json"

	"http-golang-api/db"

	"github.com/swaggo/http-swagger"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
)

type User struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Age int `json:"age"`
	Salary int `json:"salary"`
	Occupation string `json:"occupation"`
}

func addUserHandler(w http.ResponseWriter, req *http.Request) {
	user := User{
		ID: 1, 
		Name: "john",
		Age: 21,
		Salary: 100,
		Occupation: "occupation1",
	}

	jsonData, _ := json.Marshal(user)
	fmt.Println(jsonData)

	//buffer := bytes.NewBuffer(jsonData)

	//url := "http://localhost:8080/adduser"
}

func getUserHandler(w http.ResponseWriter, req *http.Request) {

}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func handleRequests() {
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/adduser/", addUserHandler)
	http.HandleFunc("/getuser/{id}", getUserHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	id := db.DbAddUser()
	fmt.Println(id)

	http.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	handleRequests()
}
