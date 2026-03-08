package app

import (
	"upcycleconnect/api-go/db"
	_ "upcycleconnect/api-go/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

func ValidateUser(userDto models.User) []string {
	var errsMsg []string

	if len(userDto.Prenom) < 4 {
		errsMsg = append(errsMsg, "Prenom length must be at least 4")
	}
	if len(userDto.Prenom) > 25 {
		errsMsg = append(errsMsg, "Prenom length must not be longer than 25")
	}
	
	if len(userDto.Nom) < 4 {
		errsMsg = append(errsMsg, "Nom length must be at least 4")
	}
	if len(userDto.Nom) > 25 {
		errsMsg = append(errsMsg, "Nom length must not be longer than 25")
	}

	if strings.ContainsAny(userDto.Prenom, "$<>[{}]*%") {
		errsMsg = append(errsMsg, "Prenom can't contain these char ($ < > [ ] { } * %)")
	}
	
	if strings.ContainsAny(userDto.Nom, "$<>[{}]*%") {
		errsMsg = append(errsMsg, "Name can't contain these char ($ < > [ ] { } * %)")
	}
    return errsMsg
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var userDto models.User

	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		http.Error(w, "Incorrect body format", http.StatusBadRequest)
		return
	}

	errsMsg := ValidateUser(userDto)

	if len(errsMsg) > 0 {
		encoded, _ := json.Marshal(errsMsg)
		http.Error(w, string(encoded), http.StatusBadRequest)
		return
	}

	err = db.CreateUser(userDto)
	if err != nil {
		http.Error(w, "pb d'insertion", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func ModifyUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("id")
	var userDto models.User

	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		http.Error(w, "Format JSON du corps invalide", http.StatusBadRequest)
		return
	}
	
	userId, err := strconv.Atoi(userIDStr)
	if err != nil || userId <= 0 {
		http.Error(w, "ID de jeu invalide dans l'URI", http.StatusBadRequest)
		return
	}
	
	userDto.Id = userId

	errsMsg := ValidateUser(userDto)

	if len(errsMsg) > 0 {
		encoded, _ := json.Marshal(errsMsg)
		http.Error(w, string(encoded), http.StatusBadRequest)
		return
	}

	err = db.ModifyUser(userId ,userDto)
	if err != nil {
		if strings.Contains(err.Error(), "aucun utilisateur trouvé") {
			http.Error(w, fmt.Sprintf("Utilisateur non trouvé avec l'ID : %d", userId), http.StatusNotFound)
			return
		}
		
		fmt.Println(err.Error())
		http.Error(w, "Erreur serveur lors de la mise à jour de l'utilisateur", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent) 
}