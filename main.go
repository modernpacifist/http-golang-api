package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"http-golang-api/db"
	"http-golang-api/types"

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

	// TODO: this is bullshit line <17-10-23, modernpacifist> //
	id := r.URL.Path[len("/api/getuser/"):]
	retrievedUser := db.GetUser(id)

	marshalers := []types.Marshaler{
		&types.JSONMarshaler{},
		&types.XMLMarshaler{},
		&types.TOMLMarshaler{},
	}

	// TODO: there must be no hardcoded sizes/indexing <17-10-23, modernpacifist> //
	info := [3][]byte{}
	for i, m := range marshalers {
		info[i], _ = m.Marshal(retrievedUser)
	}

	d := types.Data{
		JsonField: string(info[0]),
		XmlField:  string(info[1]),
		TomlField: string(info[2]),
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

	var res []types.Data

	marshalers := []types.Marshaler{
		&types.JSONMarshaler{},
		&types.XMLMarshaler{},
		&types.TOMLMarshaler{},
	}

	allUsers := db.GetAllRecords()

	for _, user := range allUsers {
		data := []byte{}
		temp := [3][]byte{}
		for i, m := range marshalers {
			data, _ = m.Marshal(user)
			temp[i] = data
		}
		d := types.Data {
			JsonField: string(temp[0]),
			XmlField: string(temp[1]),
			TomlField: string(temp[2]),
		}
		res = append(res, d)
	}

	e, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.Write(e)
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
	log.Println("Service started")

	handleRequests()
}
