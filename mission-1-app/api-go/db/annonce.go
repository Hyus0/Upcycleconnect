package db

import (
	"database/sql"
	"fmt"
	"upcycleconnect/api-go/models"
)

func CreateAnnonce(a *models.Annonce) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisée")
	}

	query := `
		INSERT INTO ANNONCE (
			id_vendeur, id_categorie, titre, description, 
			type_materiau, poids_estime_kg, prix, etat_objet, 
			statut, est_valide, type, ville, code_postal, adresse
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := Conn.Exec(
		query,
		a.IdVendeur,
		a.IdCategorie,
		a.Titre,
		a.Description,
		a.TypeMateriau,
		a.PoidsEstimeKg,
		a.Prix,
		a.EtatObjet,
		"Disponible", 
		"En attente",
		a.Type,
		a.Ville,
		a.CodePostal,
		a.Adresse,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	a.ID = int(id)
	return nil
}

func GetAllAnnonces() ([]models.Annonce, error) {
	if Conn == nil {
		return nil, fmt.Errorf("connexion DB non initialisee")
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
		var idAcheteur sql.NullInt64

		err := rows.Scan(
			&a.ID,
			&a.IdVendeur,
			&idAcheteur,
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

		if idAcheteur.Valid {
			v := int(idAcheteur.Int64)
			a.IdAcheteur = &v
		}

		annonces = append(annonces, a)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return annonces, nil
}

func GetAnnonce(id int) (*models.Annonce, error) {
	if Conn == nil {
		return nil, fmt.Errorf("connexion DB non initialisee")
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
		WHERE id = ?
	`

	row := Conn.QueryRow(query, id)

	var a models.Annonce
	var idAcheteur sql.NullInt64
	err := row.Scan(
		&a.ID,
		&a.IdVendeur,
		&idAcheteur,
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
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if idAcheteur.Valid {
		v := int(idAcheteur.Int64)
		a.IdAcheteur = &v
	}

	return &a, nil
}

func ModifyAnnonce(id int, a models.Annonce) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisee")
	}

	query := `
		UPDATE ANNONCE SET
			id_vendeur = ?,
			id_acheteur = ?,
			titre = ?,
			description = ?,
			statut = ?,
			est_valide = ?,
			prix = ?,
			etat_objet = ?,
			adresse = ?,
			ville = ?,
			code_postal = ?,
			type = ?
		WHERE id = ?
	`

	result, err := Conn.Exec(
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
		id,
	)
	if err != nil {
		return fmt.Errorf("ModifyAnnonce: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ModifyAnnonce RowsAffected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("aucune annonce trouvee avec l'ID %d", id)
	}
	return nil
}

func DeleteAnnonce(id int) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisee")
	}

	result, err := Conn.Exec("DELETE FROM ANNONCE WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("DeleteAnnonce: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteAnnonce RowsAffected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("aucune annonce trouvee avec l'ID %d", id)
	}
	return nil
}

func GetAnnoncesByUser(userId int) ([]models.Annonce, error) {
    var annonces []models.Annonce

    query := `SELECT id, id_vendeur, id_acheteur, titre, description, statut, est_valide, 
                     prix, etat_objet, adresse, ville, code_postal, date_creation, type 
              FROM ANNONCE 
              WHERE id_vendeur = ? 
              ORDER BY date_creation DESC`

    rows, err := Conn.Query(query, userId)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var a models.Annonce
        err := rows.Scan(
            &a.ID, &a.IdVendeur, &a.IdAcheteur, &a.Titre, &a.Description, &a.Statut, 
            &a.EstValide, &a.Prix, &a.EtatObjet, &a.Adresse, &a.Ville, &a.CodePostal, 
            &a.DateCreation, &a.Type,
        )
        if err != nil {
            return nil, err
        }
        annonces = append(annonces, a)
    }

    return annonces, nil
}