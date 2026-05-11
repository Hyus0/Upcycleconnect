package db

import (
	"fmt"
	"upcycleconnect/api-go/models"
)

func CreateAvis(idAuteur int, idCible int, note int, commentaire string) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisee")
	}

	query := "INSERT INTO AVIS (id_auteur, id_cible, note, commentaire) VALUES (?, ?, ?, ?)"
	_, err := Conn.Exec(query, idAuteur, idCible, note, commentaire)
	
	return err
}

func GetAvisByCible(idCible int) ([]models.Avis, error) {
	if Conn == nil {
		return nil, fmt.Errorf("connexion DB non initialisee")
	}

	query := `
		SELECT 
			a.id, 
			a.id_auteur, 
			a.id_cible, 
			a.note, 
			a.commentaire, 
			a.date_creation, 
			u.prenom 
		FROM AVIS a
		JOIN UTILISATEUR u ON a.id_auteur = u.id
		WHERE a.id_cible = ?
		ORDER BY a.date_creation DESC
	`

	rows, err := Conn.Query(query, idCible)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	avisList := []models.Avis{}
	for rows.Next() {
		var a models.Avis
		err := rows.Scan(
			&a.ID,
			&a.IdAuteur,
			&a.IdCible,
			&a.Note,
			&a.Commentaire,
			&a.DateCreation,
			&a.PrenomAuteur,
		)
		if err != nil {
			return nil, err
		}
		avisList = append(avisList, a)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return avisList, nil 
}