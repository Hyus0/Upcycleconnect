package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"upcycleconnect/api-go/db"
	"upcycleconnect/api-go/models"
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

	if len(userDto.Prenom) < 2 {
		errsMsg = append(errsMsg, "Prenom length must be at least 2")
	}
	if len(userDto.Prenom) > 50 {
		errsMsg = append(errsMsg, "Prenom length must not be longer than 50")
	}

	if len(userDto.Nom) < 2 {
		errsMsg = append(errsMsg, "Nom length must be at least 2")
	}
	if len(userDto.Nom) > 50 {
		errsMsg = append(errsMsg, "Nom length must not be longer than 50")
	}

	if len(userDto.Mail) < 5 {
		errsMsg = append(errsMsg, "Mail length must be at least 5")
	}
	if !strings.Contains(userDto.Mail, "@") || !strings.Contains(userDto.Mail, ".") {
		errsMsg = append(errsMsg, "Mail format is invalid")
	}

	if len(userDto.Password) < 8 {
		errsMsg = append(errsMsg, "Password must be at least 8 characters")
	}

	hasUpper := false
	hasSpecial := false
	for _, char := range userDto.Password {
		if char >= 'A' && char <= 'Z' {
			hasUpper = true
		}
	}
	if strings.ContainsAny(userDto.Password, "!@#$%^&*()-_=+[]{}|;:,.<>?") {
		hasSpecial = true
	}

	if !hasUpper {
		errsMsg = append(errsMsg, "Password must contain at least one uppercase letter")
	}
	if !hasSpecial {
		errsMsg = append(errsMsg, "Password must contain at least one special character")
	}

	if len(userDto.Adresse) < 5 {
		errsMsg = append(errsMsg, "Adresse length must be at least 5")
	}
	if len(userDto.Adresse) > 100 {
		errsMsg = append(errsMsg, "Adresse length must not be longer than 100")
	}

	if len(userDto.Ville) < 2 {
		errsMsg = append(errsMsg, "Ville length must be at least 2")
	}
	if len(userDto.Ville) > 50 {
		errsMsg = append(errsMsg, "Ville length must not be longer than 50")
	}

	if len(userDto.Code_postal) != 5 {
		errsMsg = append(errsMsg, "Code postal must be exactly 5 numbers")
	}

	if userDto.Role != "Particulier" && userDto.Role != "Prestataire" && userDto.Role != "Admin" {
		errsMsg = append(errsMsg, "Role must be Particulier, Prestataire or Admin")
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

	err = db.ModifyUser(userId, userDto)
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

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("id")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		http.Error(w, "ID de utilisateur invalide dans l'URI (doit être un entier positif)", http.StatusBadRequest)
		return
	}

	err = db.DeleteUser(userID)

	if err != nil {
		if strings.Contains(err.Error(), "aucun utilisateur trouvé") {
			http.Error(w, fmt.Sprintf("utilisateur non trouvé avec l'ID %d", userID), http.StatusNotFound)
			return
		}

		fmt.Println("Erreur DB lors de la suppression:", err.Error())
		http.Error(w, "Erreur serveur lors de la suppression du utilisateur", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetAllAnnonces(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	annonces, err := db.GetAllAnnonces()
	if err != nil {
		fmt.Println("Erreur DB GetAllAnnonces:", err.Error())
		http.Error(w, "Erreur de récupération des annonces", http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(annonces)
	if err != nil {
		http.Error(w, "Erreur d'encodage JSON des annonces", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", string(res))
}

func GetAnnonce(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "ID d'annonce invalide (doit être un entier positif)", http.StatusBadRequest)
		return
	}

	annonce, err := db.GetAnnonce(id)
	if err != nil {
		fmt.Println("Erreur DB GetAnnonce:", err.Error())
		http.Error(w, "Erreur serveur lors de la récupération de l'annonce", http.StatusInternalServerError)
		return
	}
	if annonce == nil {
		http.Error(w, fmt.Sprintf("Annonce non trouvée avec l'ID : %d", id), http.StatusNotFound)
		return
	}

	res, err := json.Marshal(annonce)
	if err != nil {
		http.Error(w, "Erreur d'encodage JSON", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", string(res))
}

func CreateAnnonce(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var a models.Annonce
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	if err := db.CreateAnnonce(a); err != nil {
		fmt.Println("Erreur DB CreateAnnonce:", err.Error())
		http.Error(w, "Erreur serveur lors de la création de l'annonce", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func ModifyAnnonce(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "ID d'annonce invalide dans l'URI", http.StatusBadRequest)
		return
	}

	var a models.Annonce
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, "Format JSON du corps invalide", http.StatusBadRequest)
		return
	}
	a.ID = id

	if err := db.ModifyAnnonce(id, a); err != nil {
		if strings.Contains(err.Error(), "aucune annonce trouvée") {
			http.Error(w, fmt.Sprintf("Annonce non trouvée avec l'ID : %d", id), http.StatusNotFound)
			return
		}
		fmt.Println("Erreur DB ModifyAnnonce:", err.Error())
		http.Error(w, "Erreur serveur lors de la mise à jour de l'annonce", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteAnnonce(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "ID d'annonce invalide dans l'URI (doit être un entier positif)", http.StatusBadRequest)
		return
	}

	if err := db.DeleteAnnonce(id); err != nil {
		if strings.Contains(err.Error(), "aucune annonce trouvée") {
			http.Error(w, fmt.Sprintf("Annonce non trouvée avec l'ID %d", id), http.StatusNotFound)
			return
		}
		fmt.Println("Erreur DB DeleteAnnonce:", err.Error())
		http.Error(w, "Erreur serveur lors de la suppression de l'annonce", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetAllEvenements(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	evenements, err := db.GetAllEvenements()
	if err != nil {
		fmt.Println("Erreur DB GetAllEvenements:", err.Error())
		http.Error(w, "Erreur de récupération des événements", http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(evenements)
	if err != nil {
		http.Error(w, "Erreur d'encodage JSON", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", string(res))
}

func GetEvenement(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "ID d'événement invalide (doit être un entier positif)", http.StatusBadRequest)
		return
	}

	evenement, err := db.GetEvenement(id)
	if err != nil {
		fmt.Println("Erreur DB GetEvenement:", err.Error())
		http.Error(w, "Erreur serveur lors de la récupération de l'événement", http.StatusInternalServerError)
		return
	}
	if evenement == nil {
		http.Error(w, fmt.Sprintf("Événement non trouvé avec l'ID : %d", id), http.StatusNotFound)
		return
	}

	res, err := json.Marshal(evenement)
	if err != nil {
		http.Error(w, "Erreur d'encodage JSON", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", string(res))
}

func CreateEvenement(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var e models.Evenement
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	if err := db.CreateEvenement(e); err != nil {
		fmt.Println("Erreur DB CreateEvenement:", err.Error())
		http.Error(w, "Erreur serveur lors de la création de l'événement", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func ModifyEvenement(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "ID d'événement invalide dans l'URI", http.StatusBadRequest)
		return
	}

	var e models.Evenement
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, "Format JSON du corps invalide", http.StatusBadRequest)
		return
	}
	e.ID = id

	if err := db.ModifyEvenement(id, e); err != nil {
		if strings.Contains(err.Error(), "aucun evenement trouvé") {
			http.Error(w, fmt.Sprintf("Événement non trouvé avec l'ID : %d", id), http.StatusNotFound)
			return
		}
		fmt.Println("Erreur DB ModifyEvenement:", err.Error())
		http.Error(w, "Erreur serveur lors de la mise à jour de l'événement", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteEvenement(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "ID d'événement invalide dans l'URI (doit être un entier positif)", http.StatusBadRequest)
		return
	}

	if err := db.DeleteEvenement(id); err != nil {
		if strings.Contains(err.Error(), "aucun evenement trouvé") {
			http.Error(w, fmt.Sprintf("Événement non trouvé avec l'ID %d", id), http.StatusNotFound)
			return
		}
		fmt.Println("Erreur DB DeleteEvenement:", err.Error())
		http.Error(w, "Erreur serveur lors de la suppression de l'événement", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
