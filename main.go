package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"

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
	if err != nil {
		log.Println(err)
	}

	xmlData, err := xml.Marshal(retrievedUser)
	if err != nil {
		log.Println(err)
	}

	tomlData, err := toml.Marshal(retrievedUser)
	if err != nil {
		log.Println(err)
	}

	d := types.Data{
		JsonField: string(jsonData),
		XmlField:  string(xmlData),
		TomlField: string(tomlData),
	}
	json_data, err := json.Marshal(d)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json_data)
}

// TODO: temp name <17-10-23, modernpacifist> //
func getSerializedListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed, http.StatusMethodNotAllowed", 405)
		return
	}

	var res []byte

	marshalers := []types.Marshaler{
		&types.JSONMarshaler{},
		&types.XMLMarshaler{},
		&types.TOMLMarshaler{},
	}

	allUsers := db.GetAllRecords()
	//fmt.Println(allUsers)

	for _, user := range allUsers {
		var r []byte
		for _, m := range marshalers {
			data, err := m.Marshal(user)
			fmt.Printf("%T", data)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			r = append(r, data)
		}
		res = append(res, r)
	}

	//fmt.Println(res)

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
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
