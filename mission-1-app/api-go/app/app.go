package app

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"upcycleconnect/api-go/db"
	"upcycleconnect/api-go/models"
	"upcycleconnect/api-go/passwordHashing"
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
		http.Error(w, "Erreur serveur lors de la récupération de l'utilisateur", http.StatusInternalServerError)
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

func ValidateUserInscription(userDto models.User) []string {
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

	if len(userDto.CodePostal) != 5 {
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

func ValidateUserModify(userDto models.User) []string {
	var errsMsg []string

	if len(userDto.Prenom) < 2 || len(userDto.Prenom) > 50 {
		errsMsg = append(errsMsg, "Le prénom doit faire entre 2 et 50 caractères")
	}
	if len(userDto.Nom) < 2 || len(userDto.Nom) > 50 {
		errsMsg = append(errsMsg, "Le nom doit faire entre 2 et 50 caractères")
	}

	if len(userDto.Adresse) > 0 {
		if len(userDto.Adresse) < 5 || len(userDto.Adresse) > 100 {
			errsMsg = append(errsMsg, "L'adresse doit faire entre 5 et 100 caractères")
		}
	}

	if len(userDto.Ville) > 0 {
		if len(userDto.Ville) < 2 || len(userDto.Ville) > 50 {
			errsMsg = append(errsMsg, "La ville doit faire entre 2 et 50 caractères")
		}
	}

	if len(userDto.CodePostal) != 5 {
		errsMsg = append(errsMsg, "Le code postal doit contenir exactement 5 chiffres")
	}

	if userDto.DateNaissance != "" {
		birthDate, err := time.Parse("2006-01-02", userDto.DateNaissance)
		if err != nil {
			errsMsg = append(errsMsg, "Format de date de naissance invalide (attendu: AAAA-MM-JJ)")
		} else if !birthDate.Before(time.Now()) {
			errsMsg = append(errsMsg, "La date de naissance doit être dans le passé")
		}
	}

	validRoles := map[string]bool{"Particulier": true, "Prestataire": true, "Admin": true}
	if !validRoles[userDto.Role] {
		errsMsg = append(errsMsg, "Rôle invalide")
	}

	forbiddenChars := "$<>[{}]*%"
	fields := map[string]string{
		"Le prénom": userDto.Prenom,
		"Le nom":    userDto.Nom,
		"L'adresse": userDto.Adresse,
		"La ville":  userDto.Ville,
	}

	for label, value := range fields {
		if strings.ContainsAny(value, forbiddenChars) {
			errsMsg = append(errsMsg, fmt.Sprintf("%s contient des caractères interdits ($ < > [ ] { } * %%)", label))
		}
	}

	return errsMsg
}

func ValidateUserPassword(userDto models.User) []string {
	var errsMsg []string

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

	return errsMsg
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var userDto models.User
	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Incorrect body format"})
		return
	}

	errsMsg := ValidateUserInscription(userDto)

	exists, err := db.EmailExists(userDto.Mail)
	if err != nil {
		fmt.Println("❌ ERREUR TECHNIQUE EmailExists :", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Erreur lors de la vérification du mail"})
		return
	}
	if exists {
		errsMsg = append(errsMsg, "Cet email est déjà utilisé par un autre compte.")
	}

	if len(errsMsg) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errsMsg)
		return
	}

	err = db.CreateUser(userDto)
	if err != nil {
		fmt.Println("Erreur(s) :", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Problème d'insertion"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Utilisateur créé avec succès !"})
}

func ModifyUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("id")
	token := r.Header.Get("Authorization")

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

	if !db.VerifyUserByToken(userId, token) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userDto.Id = userId

	errsMsg := ValidateUserModify(userDto)

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

func ModifyUserPassword(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("id")
	token := r.Header.Get("Authorization")

	var input struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	userId, _ := strconv.Atoi(userIDStr)

	if !db.VerifyUserByToken(userId, token) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	currentHash, err := db.GetPasswordHashed(userId)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération du compte", http.StatusInternalServerError)
		return
	}

	if !passwordHashing.VerifyPassword(input.OldPassword, currentHash) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode([]string{"L'ancien mot de passe est incorrect"})
		return
	}

	var userDto models.User
	userDto.Password = input.NewPassword

	errsMsg := ValidateUserPassword(userDto)
	if len(errsMsg) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errsMsg)
		return
	}

	newHashed, err := passwordHashing.HashPassword(input.NewPassword)
	if err != nil {
		http.Error(w, "Erreur de sécurité lors du hashage", http.StatusInternalServerError)
		return
	}
	userDto.Password = newHashed

	err = db.ModifyUserPassword(userId, userDto)
	if err != nil {
		http.Error(w, "Erreur serveur lors de la mise à jour", http.StatusInternalServerError)
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

func UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Format JSON invalide"})
		return
	}

	user, err := db.GetUserByEmail(loginData.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Email ou mot de passe incorrect"})
		return
	}

	if !passwordHashing.VerifyPassword(loginData.Password, user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Email ou mot de passe incorrect"})
		return
	}

	b := make([]byte, 16)
	rand.Read(b)
	randomToken := fmt.Sprintf("%x", b)

	err = db.UpdateUserToken(user.Id, randomToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Erreur lors de la création de session"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Bienvenue !",
		"token":   randomToken,
		"userId":  user.Id,
		"prenom":  user.Prenom,
		"nom":     user.Nom,
	})
}

func CheckSession(w http.ResponseWriter, r *http.Request) {
	id_Str := r.URL.Query().Get("id")
	token := r.Header.Get("Authorization")

	id_Int, err := strconv.Atoi(id_Str)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	isValid := db.VerifyUserByToken(id_Int, token)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"isValid": isValid})
}

func GetUserStatsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.Atoi(idStr)

	stats, err := db.GetUserStats(id)
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// Annonces
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

func ValidateAnnonce(a models.Annonce) []string {
	var errsMsg []string

	if len(a.Titre) < 5 || len(a.Titre) > 100 {
		errsMsg = append(errsMsg, "Le titre doit faire entre 5 et 100 caractères")
	}

	if len(a.Description) < 10 {
		errsMsg = append(errsMsg, "La description doit faire au moins 10 caractères")
	}

	if a.IdCategorie <= 0 {
		errsMsg = append(errsMsg, "Veuillez sélectionner une catégorie valide")
	}

	if a.Type == "Vente" && a.Prix <= 0 {
		errsMsg = append(errsMsg, "Pour une vente, le prix doit être supérieur à 0€")
	}
	if a.Type == "Don" && a.Prix != 0 {
		a.Prix = 0
	}

	if a.PoidsEstimeKg < 0 {
		errsMsg = append(errsMsg, "Le poids ne peut pas être négatif")
	}

	if len(a.CodePostal) != 5 {
		errsMsg = append(errsMsg, "Le code postal doit contenir exactement 5 chiffres")
	}
	if len(a.Ville) < 2 || len(a.Ville) > 50 {
		errsMsg = append(errsMsg, "La ville doit faire entre 2 et 50 caractères")
	}

	validTypes := map[string]bool{"Don": true, "Vente": true}
	if !validTypes[a.Type] {
		errsMsg = append(errsMsg, "Type d'annonce invalide (Don ou Vente attendu)")
	}

	validEtats := map[string]bool{"Neuf": true, "Bon etat": true, "Usage": true}
	if !validEtats[a.EtatObjet] {
		errsMsg = append(errsMsg, "État de l'objet invalide")
	}

	forbiddenChars := "$<>[{}]*%"
	fields := map[string]string{
		"Le titre":       a.Titre,
		"La description": a.Description,
		"Le matériau":    a.TypeMateriau,
		"La ville":       a.Ville,
		"L'adresse":      a.Adresse,
	}

	for label, value := range fields {
		if strings.ContainsAny(value, forbiddenChars) {
			errsMsg = append(errsMsg, fmt.Sprintf("%s contient des caractères interdits ($ < > [ ] { } * %%)", label))
		}
	}

	return errsMsg
}

func CreateAnnonce(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var a models.Annonce
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	validationErrors := ValidateAnnonce(a)
	if len(validationErrors) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": validationErrors,
		})
		return
	}

	if err := db.CreateAnnonce(&a); err != nil {
		http.Error(w, "Erreur lors de la création", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(a)
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

func GetUserAnnoncesHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	userId, err := strconv.Atoi(idStr)

	if err != nil || userId <= 0 {
		http.Error(w, "ID utilisateur invalide", http.StatusBadRequest)
		return
	}

	annonces, err := db.GetAnnoncesByUserID(userId)
	if err != nil {
		println("Erreur DB GetAnnoncesByUserID:", err.Error())
		http.Error(w, "Erreur lors de la récupération des annonces", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if annonces == nil {
		w.Write([]byte("[]"))
		return
	}

	json.NewEncoder(w).Encode(annonces)
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

func JoinEvenement(w http.ResponseWriter, r *http.Request) {
	evenementID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || evenementID <= 0 {
		http.Error(w, "ID d'événement invalide", http.StatusBadRequest)
		return
	}

	var body struct {
		UserID int `json:"id_utilisateur"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	if body.UserID <= 0 {
		http.Error(w, "ID utilisateur manquant ou invalide", http.StatusBadRequest)
		return
	}

	err = db.JoinEvenement(body.UserID, evenementID)
	if err != nil {
		if strings.Contains(err.Error(), "déjà inscrit") {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		fmt.Println("Erreur JoinEvenement:", err)
		http.Error(w, "Erreur lors de l'inscription", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Inscription à l'événement enregistrée avec succès"))
}

func QuitEvenement(w http.ResponseWriter, r *http.Request) {
	evenementID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || evenementID <= 0 {
		http.Error(w, "ID d'événement invalide", http.StatusBadRequest)
		return
	}

	var body struct {
		UserID int `json:"id_utilisateur"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	err = db.QuitEvenement(body.UserID, evenementID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Désinscription réussie"))
}
func CheckInscriptionEvenement(w http.ResponseWriter, r *http.Request) {
	evenementID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || evenementID <= 0 {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil || userID <= 0 {
		http.Error(w, "user_id invalide", http.StatusBadRequest)
		return
	}

	inscrit, err := db.IsUserInscritEvenement(userID, evenementID)
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"inscrit": inscrit})
}

func GetAllCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	categories, err := db.GetAllCategories()
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "erreur de récupération des catégorie", http.StatusInternalServerError)
	}

	res, _ := json.Marshal(categories)
	fmt.Fprintf(w, "%s", string(res))
}

func GetCategory(w http.ResponseWriter, r *http.Request) {
	categoryIDStr := r.PathValue("id")

	w.Header().Set("Content-Type", "application/json")

	categoryId, err := strconv.Atoi(categoryIDStr)
	if err != nil || categoryId <= 0 {
		http.Error(w, "ID de user invalide ou manquant (doit être un entier positif)", http.StatusBadRequest)
		return
	}
	category, err := db.GetCategory(categoryId)

	if err != nil {
		fmt.Println("Erreur DB GetUser:", err.Error())
		http.Error(w, "Erreur serveur lors de la récupération de la catégorie", http.StatusInternalServerError)
		return
	}

	if category == nil {
		http.Error(w, fmt.Sprintf("Catégorie non trouvé avec l'ID : %d", categoryId), http.StatusNotFound)
		return
	}

	res, err := json.Marshal(category)
	if err != nil {
		http.Error(w, "Erreur d'encodage JSON du résultat", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", string(res))
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var categoryDto models.Category

	err := json.NewDecoder(r.Body).Decode(&categoryDto)
	if err != nil {
		http.Error(w, "Incorrect body format", http.StatusBadRequest)
		return
	}

	err = db.CreateCategory(categoryDto)
	if err != nil {
		http.Error(w, "pb d'insertion", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func ModifyCategory(w http.ResponseWriter, r *http.Request) {
	categoryIDStr := r.PathValue("id")
	var categoryDto models.Category

	err := json.NewDecoder(r.Body).Decode(&categoryDto)
	if err != nil {
		http.Error(w, "Format JSON du corps invalide", http.StatusBadRequest)
		return
	}

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil || categoryID <= 0 {
		http.Error(w, "ID de la catégorie invalide dans l'URI", http.StatusBadRequest)
		return
	}

	categoryDto.ID = categoryID

	err = db.ModifyCategory(categoryID, categoryDto)
	if err != nil {
		if strings.Contains(err.Error(), "aucune catégorie trouvé") {
			http.Error(w, fmt.Sprintf("Categorie non trouvé avec l'ID : %d", categoryID), http.StatusNotFound)
			return
		}

		fmt.Println(err.Error())
		http.Error(w, "Erreur serveur lors de la mise à jour de la catégorie", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	categoryIDStr := r.PathValue("id")

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil || categoryID <= 0 {
		http.Error(w, "ID de la catégorie invalide dans l'URI (doit être un entier positif)", http.StatusBadRequest)
		return
	}

	err = db.DeleteCategory(categoryID)

	if err != nil {
		if strings.Contains(err.Error(), "aucune catégorie trouvé") {
			http.Error(w, fmt.Sprintf("Catégorie non trouvée avec l'ID %d", categoryID), http.StatusNotFound)
			return
		}

		fmt.Println("Erreur DB lors de la suppression:", err.Error())
		http.Error(w, "Erreur serveur lors de la suppression de la catégorie", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetAllFormations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	formations, err := db.GetAllFormations()
	if err != nil {
		fmt.Println("Erreur GetAllFormations:", err.Error())
		http.Error(w, "Erreur lors de la récupération des formations", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(formations)
}

func GetFormation(w http.ResponseWriter, r *http.Request) {
	formationID, _ := strconv.Atoi(r.PathValue("id"))

	userIDStr := r.URL.Query().Get("user_id")
	userID, _ := strconv.Atoi(userIDStr)

	formation, err := db.GetFormation(formationID, userID)
	if err != nil {
		http.Error(w, "Erreur serveur", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(formation)
}

func ValidateFormation(f models.Formation) []string {
	var errsMsg []string

	if len(f.Titre) < 5 || len(f.Titre) > 100 {
		errsMsg = append(errsMsg, "Le titre doit contenir entre 5 et 100 caractères")
	}
	if len(f.Description) < 10 {
		errsMsg = append(errsMsg, "La description est trop courte (min 10)")
	}

	if f.Capacite_max <= 0 {
		errsMsg = append(errsMsg, "La capacité maximale doit être supérieure à 0")
	}
	if f.Prix_unitaire < 0 {
		errsMsg = append(errsMsg, "Le prix ne peut pas être négatif")
	}

	if f.Type != "Atelier" && f.Type != "Webinaire" {
		errsMsg = append(errsMsg, "Le type doit être 'Atelier' ou 'Webinaire'")
	}
	if f.Statut != "Ouvert" && f.Statut != "Fermé" && f.Statut != "Annulé" {
		errsMsg = append(errsMsg, "Statut invalide (Ouvert, Fermé, Annulé)")
	}

	if f.Type == "Atelier" {
		if f.Ville == "" || len(f.Code_postal) != 5 {
			errsMsg = append(errsMsg, "Un atelier en présentiel nécessite une ville et un code postal valide")
		}
	}

	if f.Date_debut == "" || f.Date_fin == "" {
		errsMsg = append(errsMsg, "Les dates de début et de fin sont obligatoires")
	}

	return errsMsg
}

func CreateFormation(w http.ResponseWriter, r *http.Request) {
	var f models.Formation

	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	errs := ValidateFormation(f)
	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errs)
		return
	}

	if err := db.CreateFormation(f); err != nil {
		fmt.Println("Erreur CreateFormation:", err)
		http.Error(w, "Erreur lors de la création", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func ModifyFormation(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	var f models.Formation
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	errs := ValidateFormation(f)
	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errs)
		return
	}

	if err := db.ModifyFormation(id, f); err != nil {
		if strings.Contains(err.Error(), "aucune formation trouvée") {
			http.Error(w, "Formation non trouvée", http.StatusNotFound)
			return
		}
		http.Error(w, "Erreur lors de la modification", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteFormation(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	if err := db.DeleteFormation(id); err != nil {
		if strings.Contains(err.Error(), "aucune formation trouvée") {
			http.Error(w, "Formation inexistante", http.StatusNotFound)
			return
		}
		http.Error(w, "Erreur lors de la suppression", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func JoinFormation(w http.ResponseWriter, r *http.Request) {
	formationID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || formationID <= 0 {
		http.Error(w, "ID de formation invalide", http.StatusBadRequest)
		return
	}

	var body struct {
		UserID int `json:"id_utilisateur"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	if body.UserID <= 0 {
		http.Error(w, "ID utilisateur manquant ou invalide", http.StatusBadRequest)
		return
	}

	err = db.JoinFormation(body.UserID, formationID)
	if err != nil {
		if strings.Contains(err.Error(), "formation complète") {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		fmt.Println("Erreur JoinFormation:", err)
		http.Error(w, "Erreur lors de l'inscription", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Inscription enregistrée avec succès"))
}

func QuitFormation(w http.ResponseWriter, r *http.Request) {
	formationID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || formationID <= 0 {
		http.Error(w, "ID de formation invalide", http.StatusBadRequest)
		return
	}

	var body struct {
		UserID int `json:"id_utilisateur"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	err = db.QuitFormation(body.UserID, formationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Désinscription réussie"))
}

// Logistique

type ReserveRequest struct {
	SiteID int `json:"site_id"`
}

func ReserverCasier(w http.ResponseWriter, r *http.Request) {
	annonceID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	var req ReserveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON invalide", http.StatusBadRequest)
		return
	}

	pin, err := db.ReserverUnCasier(annonceID, req.SiteID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"pin":    pin,
		"status": "success",
	})
}

func GetAllSites(w http.ResponseWriter, r *http.Request) {
	sites, err := db.GetAllSites()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des sites", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sites)
}

func GetSiteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de site invalide", http.StatusBadRequest)
		return
	}

	site, err := db.GetSiteByID(id)
	if err != nil {
		http.Error(w, "Site introuvable", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(site)
}

func GetConteneurs(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	siteID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de site invalide", http.StatusBadRequest)
		return
	}

	conteneurs, err := db.GetConteneursBySite(siteID)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des conteneurs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conteneurs)
}

func RetireObjetCasierHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	idAnnonce, err := strconv.Atoi(idStr)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "ID d'annonce invalide",
		})
		return
	}

	err = db.RetireObjetCasier(idAnnonce)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Impossible de libérer le casier : " + err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Succès : L'objet est retiré et le poids du site est mis à jour",
	})
}

func GetAllProjets(w http.ResponseWriter, r *http.Request) {
	projets, err := db.GetAllProjets()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(projets)
}

func GetProjet(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	projetID, _ := strconv.Atoi(r.PathValue("id"))
	userID, _ := strconv.Atoi(r.PathValue("userId"))
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de projet invalide", http.StatusBadRequest)
		return
	}

	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	go db.IncrementVue(projetID, userID, ip)
	projet, err := db.GetProjet(id, userID)
	if err != nil || projet == nil {
		http.Error(w, "Projet introuvable", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projet)
}

func JoinProjet(w http.ResponseWriter, r *http.Request) {
	projetID, _ := strconv.Atoi(r.PathValue("id"))
	userID, err := strconv.Atoi(r.PathValue("userId"))

	if err != nil || userID <= 0 {
		http.Error(w, "ID utilisateur manquant ou invalide", http.StatusBadRequest)
		return
	}

	err = db.JoinProjet(userID, projetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Inscription réussie"})
}

func QuitProjet(w http.ResponseWriter, r *http.Request) {
	projetID, _ := strconv.Atoi(r.PathValue("id"))
	userID, err := strconv.Atoi(r.PathValue("userId"))

	if err != nil || userID <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID utilisateur invalide"})
		return
	}

	err = db.QuitProjet(userID, projetID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Vous avez quitté le projet"})
}

func ToggleLike(w http.ResponseWriter, r *http.Request) {
	projetID, _ := strconv.Atoi(r.PathValue("id"))
	userID, err := strconv.Atoi(r.PathValue("userId"))

	if err != nil || userID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = db.IncrementLike(userID, projetID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func CheckLikeStatusHandler(w http.ResponseWriter, r *http.Request) {
	projetID, _ := strconv.Atoi(r.PathValue("id"))
	userID, err := strconv.Atoi(r.PathValue("userId"))

	if err != nil || userID <= 0 {
		json.NewEncoder(w).Encode(map[string]bool{"liked": false})
		return
	}

	liked, _ := db.CheckUserLike(userID, projetID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"liked": liked})
}
