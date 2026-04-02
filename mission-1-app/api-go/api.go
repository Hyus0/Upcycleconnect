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

	http.HandleFunc("GET /", healthCheck)
	http.HandleFunc("GET /health", healthCheck)

	// Admin API consumed by the Vue backoffice.
	http.HandleFunc("GET /api/admin/metrics", app.AdminMetrics)
	http.HandleFunc("GET /{$}", healthCheck)
	http.HandleFunc("GET /api/admin/users", app.AdminListUsers)
	http.HandleFunc("POST /api/admin/users", app.AdminCreateUser)
	http.HandleFunc("PUT /api/admin/users/{id}", app.AdminUpdateUser)
	http.HandleFunc("PATCH /api/admin/users/{id}/status", app.AdminToggleUserStatus)
	http.HandleFunc("DELETE /api/admin/users/{id}", app.AdminDeleteUser)
	http.HandleFunc("GET /api/admin/prestations", app.AdminListPrestations)
	http.HandleFunc("POST /api/admin/prestations", app.AdminCreatePrestation)
	http.HandleFunc("PUT /api/admin/prestations/{id}", app.AdminUpdatePrestation)
	http.HandleFunc("DELETE /api/admin/prestations/{id}", app.AdminDeletePrestation)
	http.HandleFunc("GET /api/admin/categories", app.AdminListCategories)
	http.HandleFunc("POST /api/admin/categories", app.AdminCreateCategory)
	http.HandleFunc("PUT /api/admin/categories/{id}", app.AdminUpdateCategory)
	http.HandleFunc("DELETE /api/admin/categories/{id}", app.AdminDeleteCategory)
	http.HandleFunc("GET /api/admin/events", app.AdminListEvents)
	http.HandleFunc("POST /api/admin/events", app.AdminCreateEvent)
	http.HandleFunc("PUT /api/admin/events/{id}", app.AdminUpdateEvent)
	http.HandleFunc("DELETE /api/admin/events/{id}", app.AdminDeleteEvent)

	// Legacy routes.
	http.HandleFunc("GET /api/admin/user/{id}", app.GetUser)
	
	//Annonce
	http.HandleFunc("GET /annonces", app.GetAllAnnonces)
	http.HandleFunc("POST /annonces", app.CreateAnnonce)
	http.HandleFunc("GET /annonces/{id}", app.GetAnnonce)
	http.HandleFunc("PUT /annonces/{id}", app.ModifyAnnonce)
	http.HandleFunc("DELETE /annonces/{id}", app.DeleteAnnonce)
	http.HandleFunc("PATCH /annonces/{id}", app.ValidAnnonce)
	
	//Evenement
	http.HandleFunc("GET /evenements", app.GetAllEvenements)
	http.HandleFunc("GET /evenements/{id}", app.GetEvenement)
	http.HandleFunc("POST /evenements", app.CreateEvenement)
	http.HandleFunc("PUT /evenements/{id}", app.ModifyEvenement)
	http.HandleFunc("DELETE /evenements/{id}", app.DeleteEvenement)
	
	//Categorie
	http.HandleFunc("GET /categories", app.GetAllCategories)
	http.HandleFunc("GET /category/{id}", app.GetCategory)
	http.HandleFunc("POST /category", app.CreateCategory)
	http.HandleFunc("PUT /category/{id}", app.ModifyCategory)
	http.HandleFunc("DELETE /category/{id}", app.DeleteCategory)

	//Formation
	http.HandleFunc("GET /formations", app.GetAllFormations)
	http.HandleFunc("GET /formation/{id}", app.GetFormation)
	http.HandleFunc("POST /formation", app.CreateFormation)
	http.HandleFunc("PUT /formation/{id}", app.ModifyFormation)
	http.HandleFunc("DELETE /formation/{id}", app.DeleteFormation)
	http.HandleFunc("POST /api/formations/{id}/join", app.JoinFormation)

	fmt.Println("Listening at http://localhost:8081")
	http.ListenAndServe(":8081", nil)
}
