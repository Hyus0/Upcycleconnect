package db

import (
	"fmt"
	"upcycleconnect/api-go/models"
)

func CreateAnnonce(a models.Annonce) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisée")
	}

	query := `
        INSERT INTO ANNONCE (
            id_vendeur,
            id_acheteur,
            titre,
            description,
            statut,
            est_valide,
            prix,
            etat_objet,
            adresse,
            ville,
            code_postal,
            type
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

	res, err := Conn.Exec(
		query,
		a.IdVendeur,
		a.IdAcheteur,
		a.Titre,
		a.Description,
		a.Statut,
		a.EstValide,
		a.Prix,
		a.EtatObjet,
		a.Adresse,
		a.Ville,
		a.CodePostal,
		a.Type,
	)
	if err != nil {
		return err
	}

	if id, err := res.LastInsertId(); err == nil {
		a.ID = int(id)
	}
	return nil
}

func GetAllAnnonces() ([]models.Annonce, error) {
	if Conn == nil {
		return nil, fmt.Errorf("connexion DB non initialisée")
	}

	query := `
        SELECT
            id,
            id_vendeur,
            id_acheteur,
            titre,
            description,
            statut,
            est_valide,
            prix,
            etat_objet,
            adresse,
            ville,
            code_postal,
            date_creation,
            type
        FROM ANNONCE
    `

	rows, err := Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	annonces := []models.Annonce{}

	for rows.Next() {
		var a models.Annonce

		err := rows.Scan(
			&a.ID,
			&a.IdVendeur,
			&a.IdAcheteur,
			&a.Titre,
			&a.Description,
			&a.Statut,
			&a.EstValide,
			&a.Prix,
			&a.EtatObjet,
			&a.Adresse,
			&a.Ville,
			&a.CodePostal,
			&a.DateCreation,
			&a.Type,
		)
		if err != nil {
			return nil, err
		}

		annonces = append(annonces, a)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return annonces, nil
}
