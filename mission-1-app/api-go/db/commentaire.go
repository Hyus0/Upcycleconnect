package db

import (
	"upcycleconnect/api-go/models"
)

func GetAllCommentaires() ([]models.Commentaire, error) {
	query := `
		SELECT c.id, c.description, u.prenom, u.nom 
		FROM COMMENTAIRE c
		JOIN UTILISATEUR u ON c.id_utilisateur = u.id
		ORDER BY c.date_creation DESC`

	rows, err := Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Commentaire
	
	for rows.Next() {
		var c models.Commentaire
		var prenom, nom string
		
		err := rows.Scan(&c.ID, &c.Description, &prenom, &nom)
		if err != nil {
			return nil, err
		}
		
		if len(nom) > 0 {
			c.Auteur = prenom + " " + string(nom[0]) + "."
		} else {
			c.Auteur = prenom
		}
		comments = append(comments, c)
	}
	
	if comments == nil {
		comments = []models.Commentaire{}
	}
	return comments, nil
}