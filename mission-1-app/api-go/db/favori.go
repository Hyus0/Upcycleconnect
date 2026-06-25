package db

import (
	"fmt"
	"upcycleconnect/api-go/models"
)

func GetFavoriStatus(idAnnonce int, idUtilisateur int) (int, bool, error) {
	if Conn == nil {
		return 0, false, fmt.Errorf("connexion DB non initialisee")
	}

	var total int
	queryTotal := "SELECT COUNT(*) FROM FAVORIS WHERE id_annonce = ?"
	err := Conn.QueryRow(queryTotal, idAnnonce).Scan(&total)
	if err != nil {
		return 0, false, err
	}

	var isFavorited bool = false
	if idUtilisateur > 0 {
		var count int
		queryUser := "SELECT COUNT(*) FROM FAVORIS WHERE id_annonce = ? AND id_utilisateur = ?"
		err = Conn.QueryRow(queryUser, idAnnonce, idUtilisateur).Scan(&count)
		if err == nil {
			isFavorited = (count > 0)
		}
	}

	return total, isFavorited, nil
}

func ToggleFavori(idAnnonce int, idUtilisateur int) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisee")
	}

	var count int
	queryCheck := "SELECT COUNT(*) FROM FAVORIS WHERE id_annonce = ? AND id_utilisateur = ?"
	err := Conn.QueryRow(queryCheck, idAnnonce, idUtilisateur).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		queryDelete := "DELETE FROM FAVORIS WHERE id_annonce = ? AND id_utilisateur = ?"
		_, err = Conn.Exec(queryDelete, idAnnonce, idUtilisateur)
	} else {
		queryInsert := "INSERT INTO FAVORIS (id_annonce, id_utilisateur) VALUES (?, ?)"
		_, err = Conn.Exec(queryInsert, idAnnonce, idUtilisateur)
	}

	return err
}

func GetFavorisByUserId(userID int) ([]models.Annonce, error) {
	query := `
		SELECT a.id, a.id_vendeur, a.titre, a.description, a.prix, a.type, a.ville, a.code_postal, a.date_creation, COALESCE(a.image, '')
		FROM ANNONCE a
		INNER JOIN FAVORIS f ON a.id = f.id_annonce
		WHERE f.id_utilisateur = ?
		ORDER BY f.date_ajout DESC
	`

	rows, err := Conn.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var annonces []models.Annonce
	for rows.Next() {
		var a models.Annonce
		if err := rows.Scan(&a.ID, &a.IdVendeur, &a.Titre, &a.Description, &a.Prix, &a.Type, &a.Ville, &a.CodePostal, &a.DateCreation, &a.Image); err == nil {
			annonces = append(annonces, a)
		}
	}
	return annonces, nil
}