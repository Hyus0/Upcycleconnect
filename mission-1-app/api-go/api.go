package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"upcycleconnect/api-go/app"
	"upcycleconnect/api-go/db"
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
	http.HandleFunc("PUT /users/{id}", app.ModifyUser)
	http.HandleFunc("DELETE /users/{id}", app.DeleteUser)
	http.HandleFunc("GET /annonces", app.GetAllAnnonces)
	http.HandleFunc("POST /annonces", app.CreateAnnonce)
	http.HandleFunc("GET /annonces/{id}", app.GetAnnonce)

	fmt.Println("Listening at http://localhost:8081")
	http.ListenAndServe(":8081", nil)
}
