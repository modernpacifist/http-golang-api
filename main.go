package main

import (
	//"fmt"
	//"net/http"

	"http-golang-api/db"
)

//func helloHandler(w http.ResponseWriter, req *http.Request) {
	//fmt.Fprintf(w, "hello\n")
//}

func main() {
	db.DbConnect()

	//fmt.Println("h")
	//db.connectToDb()
}
