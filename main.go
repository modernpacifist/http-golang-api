package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"http-golang-api/db"
	"http-golang-api/types"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/swaggo/http-swagger"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
)

// TODO: ask about this... <18-10-23, modernpacifist> //
var dbManager db.DatabaseManager

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, http.StatusMethodNotAllowed", 405)
		log.Println("main.addUserHandler:Method not allowed, http.StatusMethodNotAllowed Code 405")
		return
	}

	var user types.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		log.Println("main.addUserHandler: Received invalid JSON payload")
		return
	}

	addedUserId := dbManager.AddUser(user)

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
		log.Println("main.getUserHandler:Method not allowed, http.StatusMethodNotAllowed Code 405")
		return
	}

	vars := mux.Vars(r)

	var id string

	if value, ok := vars["id"]; !ok {
		log.Println("main.getUserHandler: did not receive id in url")
		return
	} else {
		id = value
	}

	retrievedUser := dbManager.GetUser(id)

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
		log.Printf("main.getUserHandler: problem with serialization of Data struct:%v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json_data)
}

// TODO: temp name <17-10-23, modernpacifist> //
func getSerializedListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed, http.StatusMethodNotAllowed", 405)
		log.Println("main.getSerializedListHandler:Method not allowed, http.StatusMethodNotAllowed Code 405")
		return
	}

	var res []types.Data

	marshalers := []types.Marshaler{
		&types.JSONMarshaler{},
		&types.XMLMarshaler{},
		&types.TOMLMarshaler{},
	}

	allUsers := dbManager.GetAllRecords()

	for _, user := range allUsers {
		data := []byte{}
		temp := [3][]byte{}
		for i, m := range marshalers {
			data, _ = m.Marshal(user)
			temp[i] = data
		}
		d := types.Data{
			JsonField: string(temp[0]),
			XmlField:  string(temp[1]),
			TomlField: string(temp[2]),
		}
		res = append(res, d)
	}

	e, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.Write(e)
}

func handleRequests() {
	router := mux.NewRouter()

	router.HandleFunc("/api/adduser/", addUserHandler)
	router.HandleFunc("/api/getuser/{id}", getUserHandler)
	router.HandleFunc("/api/getallusers/", getSerializedListHandler)

	http.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	host := flag.String("host", dbHost, "New index")
	port := flag.String("port", dbPort, "New index")
	user := flag.String("user", dbUser, "New index")
	password := flag.String("password", dbPassword, "New index")
	dbname := flag.String("dbname", dbName, "New index")

	flag.Parse()

	if *layouts == "" {
		panic("layouts flag was not specified")
	}

	log.Println("Service started")

	//dbManager = db.DatabaseManager{
		//Host:     "localhost",
		//Port:     "5432",
		//User:     "golanguser",
		//Password: "golangpassword",
		//DBName:   "golangdb",
	//}
	dbManager = db.DatabaseManager{
		Host:     *host,
		Port:     *port,
		User:*,
		Password: "golangpassword",
		DBName:   "golangdb",
	}

	handleRequests()
}
