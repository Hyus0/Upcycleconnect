package db

import (
	"database/sql"
	"fmt"
	"upcycleconnect/api-go/models"
)

func GetAllFormations() ([]models.GetFormation, error) {
	if Conn == nil {
		return nil, fmt.Errorf("connexion DB non initialisee")
	}

	query := `
		SELECT
			f.id,
			f.id_formateur,
			f.type,
			f.titre,
			f.description,
			f.capacite_max,
			COUNT(fi.id_utilisateur) AS nb_inscrit,
			f.date_debut,
			f.date_fin,
			f.statut,
			f.prix_unitaire,
			f.adresse,
			f.ville,
			f.code_postal
		FROM FORMATION f
		LEFT JOIN FORMATION_INSCRIPTION fi ON fi.id_formation = f.id
		GROUP BY
			f.id, f.id_formateur, f.type, f.titre, f.description, f.capacite_max,
			f.date_debut, f.date_fin, f.statut, f.prix_unitaire, f.adresse, f.ville, f.code_postal
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
		return nil, fmt.Errorf("connexion DB non initialisee")
	}

	query := `
		SELECT
			f.id,
			f.id_formateur,
			f.type,
			f.titre,
			f.description,
			f.capacite_max,
			COUNT(fi.id_utilisateur) AS nb_inscrit,
			f.date_debut,
			f.date_fin,
			f.statut,
			f.prix_unitaire,
			f.adresse,
			f.ville,
			f.code_postal
		FROM FORMATION f
		LEFT JOIN FORMATION_INSCRIPTION fi ON fi.id_formation = f.id
		WHERE f.id = ?
		GROUP BY
			f.id, f.id_formateur, f.type, f.titre, f.description, f.capacite_max,
			f.date_debut, f.date_fin, f.statut, f.prix_unitaire, f.adresse, f.ville, f.code_postal
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
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisee")
	}

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
		code_postal
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := Conn.Exec(
		query,
		formation.ID_formateur,
		formation.Type,
		formation.Titre,
		formation.Description,
		formation.Capacite_max,
		formation.Date_debut,
		formation.Date_fin,
		formation.Statut,
		formation.Prix_unitaire,
		formation.Adresse,
		formation.Ville,
		formation.Code_postal,
	)
	if err != nil {
		return fmt.Errorf("CreateFormation: %v", err)
	}
	return nil
}

func ModifyFormation(id int, f models.Formation) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisee")
	}

	query := `
		UPDATE FORMATION SET
			id_formateur = ?,
			type = ?,
			titre = ?,
			description = ?,
			capacite_max = ?,
			date_debut = ?,
			date_fin = ?,
			statut = ?,
			prix_unitaire = ?,
			adresse = ?,
			ville = ?,
			code_postal = ?
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
		return fmt.Errorf("aucune formation trouvee avec l'ID %d", id)
	}
	return nil
}

func DeleteFormation(id int) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisee")
	}

	result, err := Conn.Exec("DELETE FROM FORMATION WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("DeleteFormation: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteFormation RowsAffected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("aucune formation trouvee avec l'ID %d", id)
	}
	return nil
}

func JoinFormation(userID int, formationID int) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisee")
	}

	var capaciteMax int
	var nbInscrits int
	queryCheck := `
		SELECT f.capacite_max, COUNT(fi.id_utilisateur)
		FROM FORMATION f
		LEFT JOIN FORMATION_INSCRIPTION fi ON fi.id_formation = f.id
		WHERE f.id = ?
		GROUP BY f.id, f.capacite_max
	`
	err := Conn.QueryRow(queryCheck, formationID).Scan(&capaciteMax, &nbInscrits)
	if err == sql.ErrNoRows {
		return fmt.Errorf("formation introuvable")
	}
	if err != nil {
		return fmt.Errorf("verification capacite formation: %v", err)
	}

	if nbInscrits >= capaciteMax {
		return fmt.Errorf("formation complete")
	}

	queryInsert := "INSERT INTO FORMATION_INSCRIPTION (id_utilisateur, id_formation) VALUES (?, ?)"
	if _, err := Conn.Exec(queryInsert, userID, formationID); err != nil {
		return fmt.Errorf("insertion FORMATION_INSCRIPTION: %v", err)
	}

	return nil
}
