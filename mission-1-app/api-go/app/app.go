package app

import (
	"upcycleconnect/api-go/db"
	_ "upcycleconnect/api-go/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	_ "strings"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("school", "esgi")

	users, err := db.GetAllUsers()
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "erreur de récupération des utilisateurs", http.StatusInternalServerError)
	}

	res, _ := json.Marshal(users)
	fmt.Fprintf(w, "%s", string(res))
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("id")

	w.Header().Set("Content-Type", "application/json")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		http.Error(w, "ID de user invalide ou manquant (doit être un entier positif)", http.StatusBadRequest)
		return
	}
	user, err := db.GetUser(userID) 
	
	if err != nil {
		fmt.Println("Erreur DB GetUser:", err.Error())
		http.Error(w, "Erreur serveur lors de la récupération du jeu", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, fmt.Sprintf("User non trouvé avec l'ID : %d", userID), http.StatusNotFound)
		return
	}

	res, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Erreur d'encodage JSON du résultat", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK) 
	fmt.Fprintf(w, "%s", string(res))
}