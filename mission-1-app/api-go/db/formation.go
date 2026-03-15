package db

import (
	"database/sql"
	"fmt"
	"upcycleconnect/api-go/models"
)

func GetAllFormations() ([]models.GetFormation, error) {
	if Conn == nil {
		return nil, fmt.Errorf("connexion DB non initialisée")
	}

	query := `
		SELECT
			id,
			id_formateur,
			type,
			titre,
			description,
			capacite_max,
			nb_inscrit,
			date_debut,
			date_fin,
			statut,
			prix_unitaire,
			adresse,
			ville,
			code_postal
		FROM formation
	`

	rows, err := Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	formations := []models.GetFormation{}

	for rows.Next() {
		var f models.GetFormation
		err := rows.Scan(
			&f.ID,
			&f.ID_formateur,
			&f.Type,
			&f.Titre,
			&f.Description,
			&f.Capacite_max,
			&f.Nb_inscrit,
			&f.Date_debut,
			&f.Date_fin,
			&f.Statut,
			&f.Prix_unitaire,
			&f.Adresse,
			&f.Ville,
			&f.CodePostal,
		)
		if err != nil {
			return nil, err
		}
		formations = append(formations, f)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return formations, nil
}

func GetFormation(id int) (*models.GetFormation, error) {
	if Conn == nil {
		return nil, fmt.Errorf("connexion DB non initialisée")
	}

	query := `
		SELECT
			id,
			id_formateur,
			type,
			titre,
			description,
			capacite_max,
			nb_inscrit,
			date_debut,
			date_fin,
			statut,
			prix_unitaire,
			adresse,
			ville,
			code_postal
		FROM formation WHERE id = ?
	`

	row := Conn.QueryRow(query, id)

	var f models.GetFormation
	err := row.Scan(
		&f.ID,
		&f.ID_formateur,
		&f.Type,
		&f.Titre,
		&f.Description,
		&f.Capacite_max,
		&f.Nb_inscrit,
		&f.Date_debut,
		&f.Date_fin,
		&f.Statut,
		&f.Prix_unitaire,
		&f.Adresse,
		&f.Ville,
		&f.CodePostal,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func CreateFormation(formation models.Formation) error {
	query := `INSERT INTO FORMATION (
		id_formateur,
		type,
		titre,
		description,
		capacite_max,
		date_debut,
		date_fin,
		statut,
		prix_unitaire,
		adresse,
		ville,
		code_postal)
	    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := Conn.Exec(query,
		formation.ID_formateur, formation.Type, formation.Titre, formation.Description,
		formation.Capacite_max, formation.Date_debut, formation.Date_fin, formation.Statut, formation.Prix_unitaire,
		formation.Adresse, formation.Ville, formation.Code_postal,
	)

	if err != nil {
		return fmt.Errorf("CreateFormation: %v", err)
	}
	return nil
}

func ModifyFormation(id int, f models.Formation) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisée")
	}

	query := `
		UPDATE formation SET
			id_formateur  = ?,
			type          = ?,
			titre         = ?,
			description   = ?,
			capacite_max  = ?,
			date_debut    = ?,
			date_fin      = ?,
			statut        = ?,
			prix_unitaire = ?,
			adresse       = ?,
			ville         = ?,
			code_postal   = ?
		WHERE id = ?
	`

	result, err := Conn.Exec(
		query,
		f.ID_formateur,
		f.Type,
		f.Titre,
		f.Description,
		f.Capacite_max,
		f.Date_debut,
		f.Date_fin,
		f.Statut,
		f.Prix_unitaire,
		f.Adresse,
		f.Ville,
		f.Code_postal,
		id,
	)
	if err != nil {
		return fmt.Errorf("ModifyFormation: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ModifyFormation RowsAffected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("aucune formation trouvée avec l'ID %d", id)
	}
	return nil
}

func DeleteFormation(id int) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisée")
	}

	result, err := Conn.Exec("DELETE FROM formation WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("DeleteFormation: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteFormation RowsAffected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("aucune formation trouvée avec l'ID %d", id)
	}
	return nil
}

func JoinFormation(userID int, formationID int) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisée")
	}

	queryInsert := "INSERT INTO FORMATION_INSCRIPTION (id_utilisateur, id_formation) VALUES (?, ?)"
	_, err := Conn.Exec(queryInsert, userID, formationID)
	if err != nil {
		return fmt.Errorf("échec de l'insertion dans FORMATION_INSCRIPTION : %v", err)
	}

	queryUpdate := "UPDATE formation SET nb_inscrit = nb_inscrit + 1 WHERE id = ? AND nb_inscrit < capacite_max"
	result, err := Conn.Exec(queryUpdate, formationID)
	if err != nil {
		return fmt.Errorf("échec de l'incrémentation : %v", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("impossible d'incrémenter : formation complète ou introuvable")
	}

	return nil
}