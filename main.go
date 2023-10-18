package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"sync"
	"log"
	"net/http"
	"os"

	"http-golang-api/db"
	"http-golang-api/types"
	"github.com/pelletier/go-toml"
	_ "http-golang-api/docs"


	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

// TODO: ask about this... <18-10-23, modernpacifist> //
var dbManager db.DatabaseManager

// @Summary		Add new user
// @Description	Add new user with info
// @Tags			Users
// @Accept			json
// @Produce		json
// @Param	name body	string true	"User name"
// @Success		200		{object}	string
// @Failure		400		{object}	nil
// @Router			/api/adduser/ [post]
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

func marshalize(user types.User) [3][]byte {
	var wg sync.WaitGroup
	info := [3][]byte{}

	serializedData := make(chan []byte, 3)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			jsonData, _ := json.Marshal(user)
			xmlData, _ := xml.Marshal(user)
			tomlData, _ := toml.Marshal(user)

			serializedData <- jsonData
			serializedData <- xmlData
			serializedData <- tomlData
		}()
	}

	go func() {
		wg.Wait()
		close(serializedData)
	}

	for data := range serializedData {
		log.Println(data)
	}
}

// @Summary		Get user by ID
// @Description	Get user details by ID
// @Tags			Users
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"User ID"
// @Success		200	{object}	nil
// @Failure		400	{object}	nil
// @Router			/api/getuser/{id} [get]
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

// @Summary		Get all users
// @Description	Get all users
// @Tags			Users
// @Accept			json
// @Produce		json
// @Success		200		{object}	string
// @Failure		400		{object}	nil
// @Router			/api/getallusers/ [get]
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

func handleRequests(port string) {
	router := mux.NewRouter()

	router.HandleFunc("/api/adduser/", addUserHandler)
	router.HandleFunc("/api/getuser/{id}", getUserHandler)
	router.HandleFunc("/api/getallusers/", getSerializedListHandler)

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	servicePort := os.Getenv("SERVICE_PORT")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	port := flag.String("port", servicePort, "New index")
	dbhost := flag.String("dbhost", dbHost, "New index")
	dbport := flag.String("dbport", dbPort, "New index")
	user := flag.String("user", dbUser, "New index")
	password := flag.String("password", dbPassword, "New index")
	dbname := flag.String("dbname", dbName, "New index")

	flag.Parse()

	log.Println("Service started")

	dbManager = db.DatabaseManager{
		Host:     *dbhost,
		Port:     *dbport,
		User:     *user,
		Password: *password,
		DBName:   *dbname,
	}

	handleRequests(*port)
}
