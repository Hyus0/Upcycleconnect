package main

import (
	"upcycleconnect/api-go/app"
	"upcycleconnect/api-go/db"
	"encoding/json"
	"fmt"
	"net/http"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("school", "esgi")
	res, _ := json.Marshal("en vie.")
	fmt.Fprintf(w, "%s", string(res))
}

func main() {

	db.Conn = db.NewDB()

	http.HandleFunc("GET /{$}", healthCheck)
	http.HandleFunc("GET /users", app.GetAllUsers)
	http.HandleFunc("GET /user/{id}", app.GetUser)
	http.HandleFunc("POST /users", app.CreateUser)
	//http.HandleFunc("PUT /games/{id}", app.ModifyGame)
	//http.HandleFunc("DELETE /games/{id}", app.DeleteGame)

	fmt.Println("Listening at http://localhost:8081")
	http.ListenAndServe(":8081", nil)
}
