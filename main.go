package main

import (
	"fmt"
	"net/http"
	"log"
	"encoding/json"

	"http-golang-api/db"
	"http-golang-api/types"

	"github.com/swaggo/http-swagger"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
)

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, http.StatusMethodNotAllowed")
		return
	}

	var user types.User
	err := json.NewDecoder(r.Boby).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	fmt.Println("Received user", user)

	response := struct {

	}

	//user := types.User{
		//ID: 1, 
		//Name: "john",
		//Age: 21,
		//Salary: 100,
		//Occupation: "occupation1",
	//}

	jsonData, _ := json.Marshal(user)
	fmt.Println(jsonData)

	//buffer := bytes.NewBuffer(jsonData)

	//url := "http://localhost:8080/adduser"
	addedUserId := db.AddUser()
	fmt.Println(addedUserId)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/getuser/"):]
	retrievedUser := db.GetUser(id)

	jsonData, err := json.Marshal(retrievedUser)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func handleRequests() {
	http.HandleFunc("/api/adduser/", addUserHandler)
	http.HandleFunc("/api/getuser/", getUserHandler)

	http.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handleRequests()
}
