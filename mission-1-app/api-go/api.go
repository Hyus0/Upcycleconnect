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
	http.HandleFunc("GET /api/admin/users", app.GetAllUsers)
	http.HandleFunc("GET /api/admin/user/{id}", app.GetUser)
	http.HandleFunc("POST /api/admin/users", app.CreateUser)
	http.HandleFunc("PUT /api/admin/users/{id}", app.ModifyUser)
	http.HandleFunc("DELETE /api/admin/users/{id}", app.DeleteUser)
	
	http.HandleFunc("GET /annonces", app.GetAllAnnonces)
	http.HandleFunc("POST /annonces", app.CreateAnnonce)
	http.HandleFunc("GET /annonces/{id}", app.GetAnnonce)
	http.HandleFunc("PUT /annonces/{id}", app.ModifyAnnonce)
	http.HandleFunc("DELETE /annonces/{id}", app.DeleteAnnonce)
	http.HandleFunc("PATCH /annonces/{id}", app.ValidAnnonce)

	http.HandleFunc("GET /evenements", app.GetAllEvenements)
	http.HandleFunc("GET /evenements/{id}", app.GetEvenement)
	http.HandleFunc("POST /evenements", app.CreateEvenement)
	http.HandleFunc("PUT /evenements/{id}", app.ModifyEvenement)
	http.HandleFunc("DELETE /evenements/{id}", app.DeleteEvenement)
	
	http.HandleFunc("GET /categories", app.GetAllCategories)
	http.HandleFunc("GET /category/{id}", app.GetCategory)
	http.HandleFunc("POST /category", app.CreateCategory)
	http.HandleFunc("PUT /category/{id}", app.ModifyCategory)
	http.HandleFunc("DELETE /category/{id}", app.DeleteCategory)

	fmt.Println("Listening at http://localhost:8081")
	http.ListenAndServe(":8081", nil)
}
