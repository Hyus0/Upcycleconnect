package db

import (
	"database/sql"
	"fmt"
	"strings"
	"upcycleconnect/api-go/models"
)

func GetAllEvenements() ([]models.Evenement, error) {
	if Conn == nil {
		return nil, fmt.Errorf("connexion DB non initialisee")
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
		FROM EVENEMENT
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
		return nil, fmt.Errorf("connexion DB non initialisee")
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
		FROM EVENEMENT
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
		return fmt.Errorf("connexion DB non initialisee")
	}

	query := `
		INSERT INTO EVENEMENT (
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
		return fmt.Errorf("connexion DB non initialisee")
	}

	query := `
		UPDATE EVENEMENT SET
			titre = ?,
			description = ?,
			adresse = ?,
			ville = ?,
			code_postal = ?,
			date_evenement = ?,
			type = ?
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
		return fmt.Errorf("aucun evenement trouve avec l'ID %d", id)
	}
	return nil
}

func DeleteEvenement(id int) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisee")
	}

	result, err := Conn.Exec("DELETE FROM EVENEMENT WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("DeleteEvenement: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteEvenement RowsAffected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("aucun evenement trouve avec l'ID %d", id)
	}
	return nil
}

// À AJOUTER à la fin de db/evenement.go

func JoinEvenement(userID int, evenementID int) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisee")
	}

	query := `
		INSERT INTO EVENEMENT_INSCRIPTION (id_utilisateur, id_evenement)
		VALUES (?, ?)
	`

	_, err := Conn.Exec(query, userID, evenementID)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return fmt.Errorf("utilisateur déjà inscrit à cet événement")
		}
		return fmt.Errorf("JoinEvenement: %v", err)
	}
	return nil
}

func QuitEvenement(userID int, evenementID int) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisee")
	}

	result, err := Conn.Exec(
		"DELETE FROM EVENEMENT_INSCRIPTION WHERE id_utilisateur = ? AND id_evenement = ?",
		userID, evenementID,
	)
	if err != nil {
		return fmt.Errorf("QuitEvenement: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("QuitEvenement RowsAffected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("inscription introuvable pour cet utilisateur et cet événement")
	}
	return nil
}

func IsUserInscritEvenement(userID int, evenementID int) (bool, error) {
	if Conn == nil {
		return false, fmt.Errorf("connexion DB non initialisee")
	}

	var count int
	err := Conn.QueryRow(
		"SELECT COUNT(*) FROM EVENEMENT_INSCRIPTION WHERE id_utilisateur = ? AND id_evenement = ?",
		userID, evenementID,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
