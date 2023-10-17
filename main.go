package main

import (
	"fmt"
	"net/http"
	"log"
	"encoding/json"
	"encoding/xml"

	"http-golang-api/db"
	"http-golang-api/types"

	"github.com/pelletier/go-toml"
	"github.com/swaggo/http-swagger"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
)

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, http.StatusMethodNotAllowed", 405)
		return
	}

	var user types.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	addedUserId := db.AddUser(user)

	response := struct {
		Message string `json:"id"`
	}{
		Message: fmt.Sprintf("%d", addedUserId),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed, http.StatusMethodNotAllowed", 405)
		return
	}

	id := r.URL.Path[len("/api/getuser/"):]
	retrievedUser := db.GetUser(id)

	// TODO: this must be wrapped in design pattern <17-10-23, modernpacifist> //
	jsonData, err := json.Marshal(retrievedUser)
	fmt.Println(jsonData)
	if err != nil {
		log.Fatal(err)
	}

	xmlData, err := xml.Marshal(retrievedUser)
	fmt.Println(string(xmlData))
	if err != nil {
		log.Fatal(err)
	}

	tomlData, err := toml.Marshal(retrievedUser)
	fmt.Println(string(tomlData))
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// TODO: temp name <17-10-23, modernpacifist> //
func getSerializedListHandler(w http.ResponseWriter, r *http.Request) {
	//id := r.URL.Path[len("/getuser/"):]
	//retrievedUser := db.GetUser(id)

	//jsonData, err := json.Marshal(retrievedUser)
	//if err != nil {
		//log.Fatal(err)
	//}

	//w.Header().Set("Content-Type", "application/json")
	//w.Write(jsonData)
}

func handleRequests() {
	http.HandleFunc("/api/adduser/", addUserHandler)
	http.HandleFunc("/api/getuser/", getUserHandler)
	http.HandleFunc("/api/getlist/", getSerializedListHandler)

	http.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handleRequests()
}
