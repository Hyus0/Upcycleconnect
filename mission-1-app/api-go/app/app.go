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
	"io"
	"os"
	"path/filepath"
)

//Stats

func GetPlatformStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := db.GetPlatformStats()
	if err != nil {
		http.Error(w, "Erreur lors du calcul des statistiques", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

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

	validRoles := map[string]bool{"Particulier": true, "Prestataire": true, "Salarie": true,  "Admin": true}
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
		"role": user.Role,
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

func UploadUserImages(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	if userID == "" {
		http.Error(w, `{"message": "ID utilisateur manquant"}`, http.StatusBadRequest)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, `{"message": "Fichiers trop lourds"}`, http.StatusBadRequest)
		return
	}

	var cheminProfilDB, cheminBanniereDB string

	fileProfil, headerProfil, errProfil := r.FormFile("profil")
	if errProfil == nil {
		defer fileProfil.Close()

		dossierProfil := "./uploads/profil"
		os.MkdirAll(dossierProfil, os.ModePerm)

		ext := filepath.Ext(headerProfil.Filename)
		cheminComplet := fmt.Sprintf("%s/user_%s_profil%s", dossierProfil, userID, ext)

		out, errCreate := os.Create(cheminComplet)
		if errCreate == nil {
			defer out.Close()
			io.Copy(out, fileProfil)
			cheminProfilDB = fmt.Sprintf("http://localhost:8081/img/profil/user_%s_profil%s", userID, ext)
		}
	}

	fileBanniere, headerBanniere, errBanniere := r.FormFile("banniere")
	if errBanniere == nil {
		defer fileBanniere.Close()

		dossierBanniere := "./uploads/banniere"
		os.MkdirAll(dossierBanniere, os.ModePerm)

		ext := filepath.Ext(headerBanniere.Filename)
		cheminComplet := fmt.Sprintf("%s/user_%s_banniere%s", dossierBanniere, userID, ext)

		out, errCreate := os.Create(cheminComplet)
		if errCreate == nil {
			defer out.Close()
			io.Copy(out, fileBanniere)
			cheminBanniereDB = fmt.Sprintf("http://localhost:8081/img/banniere/user_%s_banniere%s", userID, ext)
		}
	}

	if cheminProfilDB != "" {
		db.Conn.Exec("UPDATE UTILISATEUR SET image_profil = ? WHERE id = ?", cheminProfilDB, userID)
	}

	if cheminBanniereDB != "" {
		db.Conn.Exec("UPDATE UTILISATEUR SET banniere = ? WHERE id = ?", cheminBanniereDB, userID)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Upload terminé",
	})
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

func AcheterAnnonceHandler(w http.ResponseWriter, r *http.Request) {
	annonceID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID d'annonce invalide", http.StatusBadRequest)
		return
	}

	var req struct {
		UserID  int     `json:"user_id"`
		Montant float64 `json:"montant"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	if req.UserID == 0 {
		http.Error(w, "ID utilisateur manquant", http.StatusBadRequest)
		return
	}

	factureID, numeroFacture, err := db.AcheterAnnonce(annonceID, req.UserID, req.Montant)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":        "Annonce achetée avec succès",
		"facture_id":     factureID,
		"numero_facture": numeroFacture,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetUserAchatsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID utilisateur invalide", http.StatusBadRequest)
		return
	}

	achats, err := db.GetAchatsByUser(userID)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des achats", http.StatusInternalServerError)
		return
	}

	if achats == nil {
		achats = []models.Annonce{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(achats)
}

func UploadAnnonceImage(w http.ResponseWriter, r *http.Request) {
	annonceID := r.PathValue("id")
	if annonceID == "" {
		http.Error(w, `{"message": "ID annonce manquant"}`, http.StatusBadRequest)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, `{"message": "Fichier trop lourd"}`, http.StatusBadRequest)
		return
	}

	fileImage, headerImage, errImage := r.FormFile("image")
	if errImage != nil {
		http.Error(w, `{"message": "Aucune image reçue"}`, http.StatusBadRequest)
		return
	}
	defer fileImage.Close()

	dossierAnnonces := "./uploads/annonces"
	os.MkdirAll(dossierAnnonces, os.ModePerm)

	ext := filepath.Ext(headerImage.Filename)
	timestamp := time.Now().Unix() 
	nomFichier := fmt.Sprintf("annonce_%s_%d%s", annonceID, timestamp, ext)
	
	cheminComplet := fmt.Sprintf("%s/%s", dossierAnnonces, nomFichier)
	
	out, errCreate := os.Create(cheminComplet)
	if errCreate != nil {
		http.Error(w, `{"message": "Erreur d'écriture sur le disque"}`, http.StatusInternalServerError)
		return
	}
	defer out.Close()
	io.Copy(out, fileImage)

	cheminImageDB := fmt.Sprintf("http://localhost:8081/img/annonces/%s", nomFichier)
	
	_, errDb := db.Conn.Exec("UPDATE ANNONCE SET image = ? WHERE id = ?", cheminImageDB, annonceID)
	if errDb != nil {
		fmt.Println("❌ Erreur DB image annonce:", errDb)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Image enregistrée avec succès !",
		"image":   cheminImageDB,
	})
}

//Evenements

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

func GetPlanningHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	userID, err := strconv.Atoi(idStr)
	
	if err != nil || userID <= 0 {
		http.Error(w, "ID utilisateur invalide", http.StatusBadRequest)
		return
	}

	planning, err := db.GetUserPlanning(userID)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération du planning", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(planning)
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

func GetEvenementParticipantsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	evenementID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	participants, err := db.GetEvenementParticipants(evenementID)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des participants", http.StatusInternalServerError)
		return
	}

	if participants == nil {
		participants = []models.Participant{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(participants)
}

//Categorie
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

func GetFormationParticipantsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	formationID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	participants, err := db.GetFormationParticipants(formationID)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des participants", http.StatusInternalServerError)
		return
	}

	if participants == nil {
		participants = []models.Participant{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(participants)
}

// Logistique

type ReserveRequest struct {
	SiteID int `json:"site_id"`
}

func ReserverCasierHandler(w http.ResponseWriter, r *http.Request) {
	annonceID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	var req struct {
		SiteID int `json:"site_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	tokenDepot, err := db.ReserverUnCasier(annonceID, req.SiteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Casier réservé",
		"token":   tokenDepot,
	})
}

func DeposerObjetHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token  string `json:"token"`
		SiteID int    `json:"site_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	if req.Token == "" || req.SiteID == 0 {
		http.Error(w, "Token ou Site manquant", http.StatusBadRequest)
		return
	}

	tokenRetrait, err := db.DeposerObjet(req.Token, req.SiteID) 
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Objet déposé avec succès",
		"token_retrait": tokenRetrait,
	})
}

func RetireObjetCasierHandler(w http.ResponseWriter, r *http.Request) {
	conteneurID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID conteneur invalide", http.StatusBadRequest)
		return
	}

	var req struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Token manquant", http.StatusBadRequest)
		return
	}
	err = db.RecupererObjet(req.Token, conteneurID) 
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Objet récupéré avec succès"})
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

func GetProjetsByUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.PathValue("id"))
	projets, err := db.GetProjetsByUserId(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
    if projets == nil { projets = []models.ProjetUpcycling{} }
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projets)
}

func CreateProjetHandler(w http.ResponseWriter, r *http.Request) {
	var p models.ProjetUpcycling
	
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Données JSON invalides : "+err.Error(), http.StatusBadRequest)
		return
	}

	projetID, err := db.CreateProjet(p)
	if err != nil {
		http.Error(w, "Erreur BDD projet : "+err.Error(), http.StatusInternalServerError)
		return
	}

	for _, etape := range p.Etapes {
		etape.IdProjet = projetID 
		
		err := db.CreateEtape(etape)
		if err != nil {
			println("Erreur lors de l'insertion de l'étape ", etape.NumeroOrdre, ": ", err.Error())
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "success",
		"message":   "Projet et étapes enregistrés avec succès",
		"id_projet": projetID,
	})
}

func DeleteProjetHandler(w http.ResponseWriter, r *http.Request) {
	projetID, _ := strconv.Atoi(r.PathValue("id"))
	err := db.DeleteProjet(projetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func UpdateProjetHandler(w http.ResponseWriter, r *http.Request) {
    idStr := r.PathValue("id")
    projetID, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "ID invalide", http.StatusBadRequest)
        return
    }

    var p models.ProjetUpcycling
    if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
        http.Error(w, "Données invalides", http.StatusBadRequest)
        return
    }

    err = db.UpdateProjet(projetID, p)
    if err != nil {
        http.Error(w, "Erreur lors de la mise à jour : "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Projet mis à jour avec succès"})
}

func UploadProjetImageHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(5 << 20)
	if err != nil {
		http.Error(w, "Fichier trop lourd", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Image introuvable dans la requête", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	timestamp := time.Now().UnixNano()
	nouveauNom := fmt.Sprintf("projet_%d%s", timestamp, ext)

	dossierUpload := "./uploads/projets"
	os.MkdirAll(dossierUpload, os.ModePerm)

	cheminComplet := filepath.Join(dossierUpload, nouveauNom)
	dst, err := os.Create(cheminComplet)
	if err != nil {
		http.Error(w, "Erreur serveur lors de la sauvegarde", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Erreur lors de l'écriture", http.StatusInternalServerError)
		return
	}
	urlImage := fmt.Sprintf("http://localhost:8081/img/projets/%s", nouveauNom)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"url": urlImage,
	})
}
//Tips

func GetTipByRoleHandler(w http.ResponseWriter, r *http.Request) {
	role := r.PathValue("role")
	
	if role == "" {
		http.Error(w, "Rôle non spécifié", http.StatusBadRequest)
		return
	}

	tip, err := db.GetRandomTipByRole(role)
	if err != nil {
		messageErreur := fmt.Sprintf("Aucun conseil trouvé pour ce rôle : '%s'", role)
		http.Error(w, messageErreur, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tip)
}

func GetAllTipsHandler(w http.ResponseWriter, r *http.Request) {
	tips, err := db.GetAllTips()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des conseils", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tips)
}

func GetTipByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	
	if err != nil || id <= 0 {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	tip, err := db.GetTipByID(id)
	if err != nil {
		http.Error(w, "Conseil introuvable", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tip)
}

func CreateTipHandler(w http.ResponseWriter, r *http.Request) {
	var t models.Tip 
	
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Données JSON invalides", http.StatusBadRequest)
		return
	}

	err := db.CreerTip(t)
	if err != nil {
		http.Error(w, "Erreur lors de la création du tip", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Tip créé avec succès"})
}

func UpdateTipHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	var t models.Tip
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Données JSON invalides", http.StatusBadRequest)
		return
	}
	
	t.ID = id 

	err = db.ModifierTip(t)
	if err != nil {
		http.Error(w, "Erreur lors de la modification du tip", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Tip modifié avec succès"})
}

func DeleteTipHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	err = db.SupprimerTip(id)
	if err != nil {
		http.Error(w, "Erreur lors de la suppression du tip", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Tip supprimé avec succès"})
}

//Commentaire

func GetAllCommentairesHandler(w http.ResponseWriter, r *http.Request) {
	commentaires, err := db.GetAllCommentaires()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des commentaires", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commentaires)
}

//Forum

func GetForumsHandler(w http.ResponseWriter, r *http.Request) {
	forums, err := db.GetAllForums()
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(forums)
}

func CreateTopicHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID  int    `json:"user_id"`
		SalonID int    `json:"salon_id"` 
		Title   string `json:"title"`
		Sujet   string `json:"sujet"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Erreur de format de requête", http.StatusBadRequest)
		return
	}
	
	err := db.CreateForumTopic(req.UserID, req.SalonID, req.Title, req.Sujet)
	if err != nil {
		http.Error(w, "Erreur lors de la création de la discussion", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusCreated)
}

func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID  int    `json:"user_id"`
		ForumID int    `json:"forum_id"`
		Contenu string `json:"contenu"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	err := db.SendMessageForum(req.UserID, req.ForumID, req.Contenu)
	if err != nil {
		http.Error(w, "Erreur envoi", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func SignalerMessageHandler(w http.ResponseWriter, r *http.Request) {
	messageIDStr := r.PathValue("id") 
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		http.Error(w, "ID message invalide", http.StatusBadRequest)
		return
	}

	var req struct {
		UserID int    `json:"user_id"`
		Motif  string `json:"motif"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Erreur de format de requête", http.StatusBadRequest)
		return
	}

	if err := db.SignalerMessageForum(messageID, req.UserID, req.Motif); err != nil {
		http.Error(w, "Erreur lors du signalement", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func TopMessageSignaleHandler(w http.ResponseWriter, r *http.Request) {
	reportedMessages, err := db.GetTopMessagesSignales()
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	if reportedMessages == nil {
		reportedMessages = []models.ReportedMessage{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reportedMessages)
}

func BanUserForumHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("id") 
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "ID utilisateur invalide", http.StatusBadRequest)
		return
	}

	var req struct {
		Ban bool `json:"ban"` 
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Erreur de format de requête", http.StatusBadRequest)
		return
	}

	if err := db.ToggleBanForum(userID, req.Ban); err != nil {
		http.Error(w, "Erreur lors du ban", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func IgnoreSignalementHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	messageID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"message": "ID invalide"}`, http.StatusBadRequest)
		return
	}

	if err := db.IgnoreSignalementForum(messageID); err != nil {
		http.Error(w, `{"message": "Erreur lors de la suppression du signalement"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Signalement ignoré avec succès"}`))
}


func DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	messageID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	if err := db.DeleteMessageForum(messageID); err != nil {
		http.Error(w, "Erreur lors de la suppression", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteTopicHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	topicID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	if err := db.DeleteTopicForum(topicID); err != nil {
		http.Error(w, "Erreur lors de la suppression", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetRecentMessagesHandler(w http.ResponseWriter, r *http.Request) {
	messages, err := db.GetRecentForumMessages()
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	if messages == nil {
		messages = []models.ForumMessage{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func GetBannedUsersHandler(w http.ResponseWriter, r *http.Request) {
	bannedUsers, err := db.GetBannedUsers()
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	if bannedUsers == nil {
		bannedUsers = []models.BannedUser{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bannedUsers)
}

func GetModerationTopicsHandler(w http.ResponseWriter, r *http.Request) {
	topics, err := db.GetModerationTopics()
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	if topics == nil {
		topics = []models.ModTopic{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topics)
}

// Messagerie privee

func GetSubscriptionStatusHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || userID <= 0 {
		http.Error(w, "ID utilisateur invalide", http.StatusBadRequest)
		return
	}

	isSubscriber, err := db.HasActiveDMSubscription(userID)
	if err != nil {
		http.Error(w, "Erreur abonnement", http.StatusInternalServerError)
		return
	}

	used, err := db.CountDistinctAnnonceVendorsContacted(userID)
	if err != nil {
		http.Error(w, "Erreur compteur messagerie", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"is_subscriber": isSubscriber,
		"used":          used,
		"limit":         5,
		"price":         2.99,
	})
}

func GetConversationsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || userID <= 0 {
		http.Error(w, "ID utilisateur invalide", http.StatusBadRequest)
		return
	}

	conversations, err := db.ListConversations(userID)
	if err != nil {
		http.Error(w, "Erreur conversations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversations)
}

func StartConversationHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || userID <= 0 {
		http.Error(w, "ID utilisateur invalide", http.StatusBadRequest)
		return
	}

	var req struct {
		TargetUserID int  `json:"target_user_id"`
		AnnonceID    *int `json:"annonce_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requete invalide", http.StatusBadRequest)
		return
	}

	result, err := db.StartConversation(userID, req.TargetUserID, req.AnnonceID)
	if err != nil {
		http.Error(w, "Erreur creation conversation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if !result.Allowed {
		w.WriteHeader(http.StatusPaymentRequired)
	}
	json.NewEncoder(w).Encode(result)
}

func GetConversationMessagesHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("id"))
	conversationID, errConv := strconv.Atoi(r.PathValue("conversationId"))
	if err != nil || errConv != nil || userID <= 0 || conversationID <= 0 {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	messages, err := db.GetConversationMessages(userID, conversationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func SendDMMessageHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("id"))
	conversationID, errConv := strconv.Atoi(r.PathValue("conversationId"))
	if err != nil || errConv != nil || userID <= 0 || conversationID <= 0 {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requete invalide", http.StatusBadRequest)
		return
	}

	message, err := db.SendDMMessage(userID, conversationID, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

func GetConversationStateHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("id"))
	conversationID, errConv := strconv.Atoi(r.PathValue("conversationId"))
	if err != nil || errConv != nil || userID <= 0 || conversationID <= 0 {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	state, err := db.GetDMThreadState(userID, conversationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

func CreateDMOfferHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("id"))
	conversationID, errConv := strconv.Atoi(r.PathValue("conversationId"))
	if err != nil || errConv != nil || userID <= 0 || conversationID <= 0 {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requete invalide", http.StatusBadRequest)
		return
	}

	offer, err := db.CreateDMOffer(userID, conversationID, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(offer)
}

func RespondDMOfferHandler(w http.ResponseWriter, r *http.Request) {
	userID, errUser := strconv.Atoi(r.PathValue("id"))
	offerID, errOffer := strconv.Atoi(r.PathValue("offerId"))
	if errUser != nil || errOffer != nil || userID <= 0 || offerID <= 0 {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	var req struct {
		Action string `json:"action"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requete invalide", http.StatusBadRequest)
		return
	}

	offer, sale, err := db.RespondDMOffer(userID, offerID, strings.ToLower(strings.TrimSpace(req.Action)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"offer": offer,
		"sale":  sale,
	})
}

func ConfirmDMSaleReceptionHandler(w http.ResponseWriter, r *http.Request) {
	userID, errUser := strconv.Atoi(r.PathValue("id"))
	saleID, errSale := strconv.Atoi(r.PathValue("saleId"))
	if errUser != nil || errSale != nil || userID <= 0 || saleID <= 0 {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	sale, err := db.ConfirmDMSaleReception(userID, saleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sale)
}

func ReviewDMSaleHandler(w http.ResponseWriter, r *http.Request) {
	userID, errUser := strconv.Atoi(r.PathValue("id"))
	saleID, errSale := strconv.Atoi(r.PathValue("saleId"))
	if errUser != nil || errSale != nil || userID <= 0 || saleID <= 0 {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	var req struct {
		Note        int    `json:"note"`
		Commentaire string `json:"commentaire"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requete invalide", http.StatusBadRequest)
		return
	}

	sale, err := db.ReviewDMSale(userID, saleID, req.Note, req.Commentaire)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sale)
}

//Favori

func GetFavoriStatusHandler(w http.ResponseWriter, r *http.Request) {
	idAnnonce, _ := strconv.Atoi(r.PathValue("id"))
	userID, _ := strconv.Atoi(r.PathValue("userId"))

	total, isFavorited, err := db.GetFavoriStatus(idAnnonce, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total":        total,
		"is_favorited": isFavorited,
	})
}

func ToggleFavoriHandler(w http.ResponseWriter, r *http.Request) {
	idAnnonce, _ := strconv.Atoi(r.PathValue("id"))
	userID, err := strconv.Atoi(r.PathValue("userId"))

	if err != nil || userID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = db.ToggleFavori(idAnnonce, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
}

func GetMesFavorisHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID utilisateur invalide", http.StatusBadRequest)
		return
	}

	annonces, err := db.GetFavorisByUserId(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if annonces == nil {
		annonces = []models.Annonce{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(annonces)
}

//Avis

func GetUserAvisHandler(w http.ResponseWriter, r *http.Request) {
	idCibleStr := r.PathValue("id")
	idCible, err := strconv.Atoi(idCibleStr)

	if err != nil || idCible <= 0 {
		http.Error(w, "ID utilisateur invalide", http.StatusBadRequest)
		return
	}

	avisList, err := db.GetAvisByCible(idCible)
	if err != nil {
		http.Error(w, "Erreur serveur lors de la récupération des avis", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(avisList)
}

func AddAvisHandler(w http.ResponseWriter, r *http.Request) {
	idCibleStr := r.PathValue("id")
	idCible, err := strconv.Atoi(idCibleStr)

	if err != nil || idCible <= 0 {
		http.Error(w, "ID utilisateur cible invalide", http.StatusBadRequest)
		return
	}

	var req struct {
		IdAuteur    int    `json:"id_auteur"`
		Note        int    `json:"note"`
		Commentaire string `json:"commentaire"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Format de requête invalide", http.StatusBadRequest)
		return
	}

	err = db.CreateAvis(req.IdAuteur, idCible, req.Note, req.Commentaire)
	if err != nil {
		http.Error(w, "Erreur lors de l'ajout de l'avis", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Avis ajouté avec succès"})
}

//Follow

func GetFollowStatusHandler(w http.ResponseWriter, r *http.Request) {
	idProfil, _ := strconv.Atoi(r.PathValue("id"))
	idConnecte, _ := strconv.Atoi(r.PathValue("userId"))

	followers, following, isFollowing, err := db.GetFollowStatus(idProfil, idConnecte)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"followers":    followers,
		"following":    following, 
		"is_following": isFollowing,
	})
}

func ToggleFollowHandler(w http.ResponseWriter, r *http.Request) {
	idSuivi, _ := strconv.Atoi(r.PathValue("id"))
	idAbonne, err := strconv.Atoi(r.PathValue("userId"))

	if err != nil || idAbonne <= 0 {
		http.Error(w, "Vous devez être connecté", http.StatusUnauthorized)
		return
	}
    
    if idSuivi == idAbonne {
        http.Error(w, "Vous ne pouvez pas vous suivre vous-même", http.StatusBadRequest)
		return
    }

	err = db.ToggleFollowUser(idSuivi, idAbonne)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
}

//Panier

type AddPanierRequest struct {
	TypeItem     string  `json:"type_item"`
	ReferenceID  int     `json:"reference_id"`
	PrixUnitaire float64 `json:"prix_unitaire"`
}

func GetPanierHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.PathValue("id"))

	items, err := db.GetPanierByUserId(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func AddToPanierHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.PathValue("id"))

	var req AddPanierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	err := db.AddToPanier(userID, req.TypeItem, req.ReferenceID, req.PrixUnitaire)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Ajouté au panier avec succès"}`))
}

func RemoveFromPanierHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("🗑️ Requête de suppression reçue !")

	userIDStr := r.PathValue("id")
	itemIDStr := r.PathValue("itemId")

	fmt.Printf("Données URL -> UserID: %s | ItemID: %s\n", userIDStr, itemIDStr)

	userID, err1 := strconv.Atoi(userIDStr)
	itemID, err2 := strconv.Atoi(itemIDStr)

	if err1 != nil || err2 != nil {
		fmt.Println("❌ Erreur : ID invalide dans l'URL")
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	err := db.RemoveFromPanier(itemID, userID)
	if err != nil {
		fmt.Println("❌ Erreur de suppression en DB :", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Article supprimé avec succès !")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Article supprimé"}`))
}

func CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Requête de Checkout reçue !")

	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID utilisateur invalide", http.StatusBadRequest)
		return
	}

	commandeID, err := db.Checkout(userID)
	if err != nil {
		fmt.Println("❌ Erreur lors du Checkout :", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Commande validée ! ID :", commandeID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "Commande créée avec succès",
		"commande_id": commandeID,
	})
}

func CheckoutWithInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || userID <= 0 {
		http.Error(w, "ID utilisateur invalide", http.StatusBadRequest)
		return
	}

	checkout, err := db.Checkout(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":        "Commande creee avec succes",
		"commande_id":    checkout.CommandeID,
		"transaction_id": checkout.TransactionID,
		"facture_id":     checkout.FactureID,
		"numero_facture": checkout.NumeroFacture,
		"facture_url":    fmt.Sprintf("/users/%d/factures/%d/download", userID, checkout.FactureID),
	})
}

func GetFacturesHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || userID <= 0 {
		http.Error(w, "ID utilisateur invalide", http.StatusBadRequest)
		return
	}

	factures, err := db.GetFacturesByUser(userID)
	if err != nil {
		http.Error(w, "Erreur lors de la recuperation des factures", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(factures)
}

func DownloadFactureHandler(w http.ResponseWriter, r *http.Request) {
	userID, errUser := strconv.Atoi(r.PathValue("id"))
	factureID, errFacture := strconv.Atoi(r.PathValue("factureId"))
	if errUser != nil || errFacture != nil || userID <= 0 || factureID <= 0 {
		http.Error(w, "Parametres invalides", http.StatusBadRequest)
		return
	}

	document, facture, err := db.BuildFactureHTML(userID, factureID)
	if err != nil {
		http.Error(w, "Facture introuvable", http.StatusNotFound)
		return
	}

	fileName := fmt.Sprintf("%s.html", strings.ReplaceAll(facture.NumeroFacture, "/", "-"))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(document))
}

func SendFactureByMailHandler(w http.ResponseWriter, r *http.Request) {
	userID, errUser := strconv.Atoi(r.PathValue("id"))
	factureID, errFacture := strconv.Atoi(r.PathValue("factureId"))
	if errUser != nil || errFacture != nil || userID <= 0 || factureID <= 0 {
		http.Error(w, "Parametres invalides", http.StatusBadRequest)
		return
	}

	facture, err := db.GetFactureByIDForUser(userID, factureID)
	if err != nil {
		http.Error(w, "Facture introuvable", http.StatusNotFound)
		return
	}

	message := fmt.Sprintf("Votre facture %s est disponible dans votre espace Factures. Envoi SMTP a brancher pour un mail reel vers %s.", facture.NumeroFacture, facture.AcheteurEmail)
	if err := db.CreerNotification(userID, userID, "Message", "Facture disponible", message); err != nil {
		http.Error(w, "Facture preparee mais notification impossible", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Facture preparee. Sans SMTP configure, une notification locale a ete creee.",
		"email":   facture.AcheteurEmail,
	})
}

//Notification

func GetAllNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	notifs, err := db.GetAllNotifications()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des notifications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifs)
}

func GetNotificationHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de notification invalide", http.StatusBadRequest)
		return
	}

	notif, err := db.GetNotificationByID(id)
	if err != nil {
		http.Error(w, "Notification introuvable", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notif)
}

func GetNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID utilisateur invalide", http.StatusBadRequest)
		return
	}

	notifs, err := db.GetNotificationsByUser(userID)
	if err != nil {
		fmt.Println("Erreur récupération notifications:", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	if notifs == nil {
		notifs = []models.Notification{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifs)
}

func MarquerNotificationLueHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	notifID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID notification invalide", http.StatusBadRequest)
		return
	}

	var body struct {
		UserID int `json:"user_id"`
	}
	
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		fmt.Println("Erreur de lecture JSON :", err)
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	err = db.MarquerNotificationLue(notifID, body.UserID)
	if err != nil {
		http.Error(w, "Erreur BDD", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

//Traduction
func GetLanguesHandler(w http.ResponseWriter, r *http.Request) {
	langues, err := db.GetLangues()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des langues", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(langues)
}

func GetTraductionsHandler(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code") 
	if code == "" {
		code = "fr"
	}

	traductions, err := db.GetTraductionsByCode(code)
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	if traductions == nil {
		traductions = make(map[string]string)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(traductions)
}

func UpdateLangueHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	var req struct {
		IDLangue int `json:"id_langue"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	err = db.UpdateLangueUtilisateur(userID, req.IDLangue)
	if err != nil {
		http.Error(w, "Erreur base de données", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}