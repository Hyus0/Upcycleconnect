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
		return nil, fmt.Errorf("connexion DB non initialisée")
	}

	query := `
		SELECT
			id, id_vendeur, 
			COALESCE(id_acheteur, 0), 
			COALESCE(id_casier, 0), 
			id_categorie, titre, description, type_materiau, poids_estime_kg, prix,
			etat_objet, statut, est_valide, 
			COALESCE(code_pin_depot, ''), 
			COALESCE(code_barre_retrait, ''),
			date_creation, 
			COALESCE(date_depot_effective, '0001-01-01 00:00:00'), 
			COALESCE(date_recuperation_effective, '0001-01-01 00:00:00'),
			type, ville, code_postal, adresse
		FROM ANNONCE
	`

	rows, err := Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var annonces []models.Annonce
	for rows.Next() {
		var a models.Annonce
		
		err := rows.Scan(
			&a.ID, &a.IdVendeur, &a.IdAcheteur, &a.IdCasier, &a.IdCategorie,
			&a.Titre, &a.Description, &a.TypeMateriau, &a.PoidsEstimeKg, &a.Prix,
			&a.EtatObjet, &a.Statut, &a.EstValide, &a.CodePinDepot, &a.CodeBarreRetrait,
			&a.DateCreation, &a.DateDepotEffective, &a.DateRecuperationEffective,
			&a.Type, &a.Ville, &a.CodePostal, &a.Adresse,
		)
		if err != nil {
			fmt.Println("Erreur Scan GetAllAnnonces:", err)
			return nil, err
		}
		annonces = append(annonces, a)
	}

	return annonces, nil
}

func GetAnnoncesByUserID(userID int) ([]models.Annonce, error) {
    query := `
        SELECT
            id, id_vendeur, 
            COALESCE(id_acheteur, 0), 
            COALESCE(id_casier, 0), 
            id_categorie, titre, description, type_materiau, poids_estime_kg, prix,
            etat_objet, statut, est_valide, 
            COALESCE(code_pin_depot, ''), 
            COALESCE(code_barre_retrait, ''),
            date_creation, 
            CAST(COALESCE(date_depot_effective, '0001-01-01 00:00:00') AS DATETIME), 
            CAST(COALESCE(date_recuperation_effective, '0001-01-01 00:00:00') AS DATETIME),
            type, ville, code_postal, adresse
        FROM ANNONCE
        WHERE id_vendeur = ?
    `

    rows, err := Conn.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var annonces []models.Annonce
    for rows.Next() {
        var a models.Annonce

        err := rows.Scan(
            &a.ID, &a.IdVendeur, &a.IdAcheteur, &a.IdCasier, &a.IdCategorie, 
            &a.Titre, &a.Description, &a.TypeMateriau, &a.PoidsEstimeKg,  
            &a.Prix, &a.EtatObjet, &a.Statut, &a.EstValide,                
            &a.CodePinDepot, &a.CodeBarreRetrait, 
            &a.DateCreation, 
            &a.DateDepotEffective, &a.DateRecuperationEffective,
            &a.Type, &a.Ville, &a.CodePostal, &a.Adresse,
        )
        if err != nil {
            fmt.Println("Erreur Scan Index 16 ou autre:", err)
            return nil, err
        }
        annonces = append(annonces, a)
    }
    return annonces, nil
}

func GetAnnonce(id int) (*models.Annonce, error) {
	query := `
		SELECT
			id, id_vendeur, id_acheteur, id_casier, id_categorie,
			titre, description, type_materiau, poids_estime_kg, prix,
			etat_objet, statut, est_valide, code_pin_depot, code_barre_retrait,
			date_creation, date_depot_effective, date_recuperation_effective,
			type, ville, code_postal, adresse
		FROM ANNONCE
		WHERE id = ?
	`

	var a models.Annonce
	err := Conn.QueryRow(query, id).Scan(
		&a.ID, &a.IdVendeur, &a.IdAcheteur, &a.IdCasier, &a.IdCategorie,
		&a.Titre, &a.Description, &a.TypeMateriau, &a.PoidsEstimeKg, &a.Prix,
		&a.EtatObjet, &a.Statut, &a.EstValide, &a.CodePinDepot, &a.CodeBarreRetrait,
		&a.DateCreation, &a.DateDepotEffective, &a.DateRecuperationEffective,
		&a.Type, &a.Ville, &a.CodePostal, &a.Adresse,
	)

	if err == sql.ErrNoRows { return nil, nil }
	if err != nil { return nil, err }

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