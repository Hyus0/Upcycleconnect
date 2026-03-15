package db

import (
	"database/sql"
	"fmt"
	"upcycleconnect/api-go/models"
)

func GetAllEvenements() ([]models.Evenement, error) {
	if Conn == nil {
		return nil, fmt.Errorf("connexion DB non initialisée")
	}

	query := `
		SELECT
			id,
			titre,
			description,
			adresse,
			ville,
			code_postal,
			date_creation,
			date_evenement,
			type
		FROM evenement
	`

	rows, err := Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	evenements := []models.Evenement{}

	for rows.Next() {
		var e models.Evenement
		err := rows.Scan(
			&e.ID,
			&e.Titre,
			&e.Description,
			&e.Adresse,
			&e.Ville,
			&e.CodePostal,
			&e.DateCreation,
			&e.DateEvenement,
			&e.Type,
		)
		if err != nil {
			return nil, err
		}
		evenements = append(evenements, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return evenements, nil
}

func GetEvenement(id int) (*models.Evenement, error) {
	if Conn == nil {
		return nil, fmt.Errorf("connexion DB non initialisée")
	}

	query := `
		SELECT
			id,
			titre,
			description,
			adresse,
			ville,
			code_postal,
			date_creation,
			date_evenement,
			type
		FROM evenement
		WHERE id = ?
	`

	row := Conn.QueryRow(query, id)

	var e models.Evenement
	err := row.Scan(
		&e.ID,
		&e.Titre,
		&e.Description,
		&e.Adresse,
		&e.Ville,
		&e.CodePostal,
		&e.DateCreation,
		&e.DateEvenement,
		&e.Type,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func CreateEvenement(e models.Evenement) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisée")
	}

	query := `
		INSERT INTO evenement (
			titre,
			description,
			adresse,
			ville,
			code_postal,
			date_evenement,
			type
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := Conn.Exec(
		query,
		e.Titre,
		e.Description,
		e.Adresse,
		e.Ville,
		e.CodePostal,
		e.DateEvenement,
		e.Type,
	)
	if err != nil {
		return fmt.Errorf("CreateEvenement: %v", err)
	}
	return nil
}

func ModifyEvenement(id int, e models.Evenement) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisée")
	}

	query := `
		UPDATE evenement SET
			titre          = ?,
			description    = ?,
			adresse        = ?,
			ville          = ?,
			code_postal    = ?,
			date_evenement = ?,
			type           = ?
		WHERE id = ?
	`

	result, err := Conn.Exec(
		query,
		e.Titre,
		e.Description,
		e.Adresse,
		e.Ville,
		e.CodePostal,
		e.DateEvenement,
		e.Type,
		id,
	)
	if err != nil {
		return fmt.Errorf("ModifyEvenement: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ModifyEvenement RowsAffected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("aucun evenement trouvé avec l'ID %d", id)
	}
	return nil
}

func DeleteEvenement(id int) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisée")
	}

	result, err := Conn.Exec("DELETE FROM evenement WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("DeleteEvenement: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteEvenement RowsAffected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("aucun evenement trouvé avec l'ID %d", id)
	}
	return nil
}
