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

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5174")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	db.Conn = db.NewDB()

	http.HandleFunc("GET /", healthCheck)
	http.HandleFunc("GET /health", healthCheck)

	// Admin API consumed by the Vue backoffice.
	http.HandleFunc("GET /api/admin/metrics", app.AdminMetrics)
	http.HandleFunc("GET /{$}", healthCheck)
	http.HandleFunc("GET /api/admin/users", app.AdminListUsers)
	http.HandleFunc("GET /api/admin/users/{id}", app.AdminGetUser)
	http.HandleFunc("POST /api/admin/users", app.AdminCreateUser)
	http.HandleFunc("PUT /api/admin/users/{id}", app.AdminUpdateUser)
	http.HandleFunc("PATCH /api/admin/users/{id}/status", app.AdminToggleUserStatus)
	http.HandleFunc("DELETE /api/admin/users/{id}", app.AdminDeleteUser)
	http.HandleFunc("GET /api/admin/prestations", app.AdminListPrestations)
	http.HandleFunc("GET /api/admin/prestations/{id}", app.AdminGetPrestation)
	http.HandleFunc("POST /api/admin/prestations", app.AdminCreatePrestation)
	http.HandleFunc("PUT /api/admin/prestations/{id}", app.AdminUpdatePrestation)
	http.HandleFunc("DELETE /api/admin/prestations/{id}", app.AdminDeletePrestation)
	http.HandleFunc("GET /api/admin/categories", app.AdminListCategories)
	http.HandleFunc("GET /api/admin/categories/{id}", app.AdminGetCategory)
	http.HandleFunc("POST /api/admin/categories", app.AdminCreateCategory)
	http.HandleFunc("PUT /api/admin/categories/{id}", app.AdminUpdateCategory)
	http.HandleFunc("DELETE /api/admin/categories/{id}", app.AdminDeleteCategory)
	http.HandleFunc("GET /api/admin/events", app.AdminListEvents)
	http.HandleFunc("GET /api/admin/events/{id}", app.AdminGetEvent)
	http.HandleFunc("POST /api/admin/events", app.AdminCreateEvent)
	http.HandleFunc("PUT /api/admin/events/{id}", app.AdminUpdateEvent)
	http.HandleFunc("DELETE /api/admin/events/{id}", app.AdminDeleteEvent)
	http.HandleFunc("GET /api/admin/moderation/queue", app.AdminModerationQueue)
	http.HandleFunc("PATCH /api/admin/moderation/prestations/{id}/publish", app.AdminPublishPrestation)
	http.HandleFunc("PATCH /api/admin/moderation/prestations/{id}/archive", app.AdminArchivePrestation)
	http.HandleFunc("PATCH /api/admin/moderation/events/{id}/publish", app.AdminPublishEvent)
	http.HandleFunc("PATCH /api/admin/moderation/events/{id}/archive", app.AdminArchiveEvent)
	http.HandleFunc("GET /api/admin/finance/overview", app.AdminFinanceOverview)
	http.HandleFunc("GET /api/admin/notifications", app.AdminListNotifications)
	http.HandleFunc("POST /api/admin/notifications", app.AdminCreateNotification)
	http.HandleFunc("PATCH /api/admin/notifications/{id}/status", app.AdminUpdateNotificationStatus)
	http.HandleFunc("DELETE /api/admin/notifications/{id}", app.AdminDeleteNotification)

	// Legacy routes.
	http.HandleFunc("GET /api/admin/user/{id}", app.GetUser)

	//Utilisateur
	http.HandleFunc("GET /users", app.GetAllUsers)
	http.HandleFunc("POST /users", app.CreateUser)
	http.HandleFunc("GET /users/{id}", app.GetUser)
	http.HandleFunc("PUT /users/{id}", app.ModifyUser)
	http.HandleFunc("DELETE /users/{id}", app.DeleteUser)
	http.HandleFunc("POST /login", app.UserLogin)
	http.HandleFunc("GET /check-session", app.CheckSession)
	http.HandleFunc("PUT /users/{id}/password", app.ModifyUserPassword)

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
	http.ListenAndServe(":8081", enableCORS(http.DefaultServeMux))
}
