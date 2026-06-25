package db

import (
	"upcycleconnect/api-go/models"
)

func GetRandomTipByRole(role string) (models.Tip, error) {
	var tip models.Tip


	query := `
		SELECT id, id_createur, titre, description, video_url, role_cible, date_creation, actif
		FROM tips
		WHERE role_cible = ? AND actif = 1
		ORDER BY RAND()
		LIMIT 1`

	err := Conn.QueryRow(query, role).Scan(
		&tip.ID,
		&tip.ID_createur,
		&tip.Titre,
		&tip.Description,
		&tip.Video_url,
		&tip.Role_cible,
		&tip.Date_creation,
		&tip.Actif,
	)

	if err != nil {
		return tip, err 
	}

	return tip, nil
}

func GetAllTips() ([]models.Tip, error) {
	query := `
		SELECT id, id_createur, titre, description, video_url, role_cible, date_creation, actif 
		FROM tips 
		ORDER BY date_creation DESC`

	rows, err := Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tips []models.Tip
	for rows.Next() {
		var t models.Tip
		err := rows.Scan(
			&t.ID, 
			&t.ID_createur, 
			&t.Titre, 
			&t.Description, 
			&t.Video_url,
			&t.Role_cible, 
			&t.Date_creation, 
			&t.Actif,
		)
		if err != nil {
			return nil, err
		}
		tips = append(tips, t)
	}

	if tips == nil {
		tips = []models.Tip{}
	}

	return tips, nil
}

func GetTipByID(id int) (models.Tip, error) {
	var t models.Tip
	query := `
		SELECT id, id_createur, titre, description, video_url, role_cible, date_creation, actif 
		FROM tips 
		WHERE id = ?`

	err := Conn.QueryRow(query, id).Scan(
		&t.ID, 
		&t.ID_createur, 
		&t.Titre, 
		&t.Description, 
		&t.Video_url,
		&t.Role_cible, 
		&t.Date_creation, 
		&t.Actif,
	)

	return t, err
}

func CreerTip(t models.Tip) error {
	query := `
		INSERT INTO TIPS (id_createur, titre, description, video_url, role_cible, date_creation, actif) 
		VALUES (?, ?, ?, ?, ?, NOW(), ?)`
	
	_, err := Conn.Exec(query, 
		t.ID_createur, 
		t.Titre, 
		t.Description, 
		t.Video_url, 
		t.Role_cible, 
		t.Actif,
	)
	return err
}

func ModifierTip(t models.Tip) error {
	query := `
		UPDATE TIPS 
		SET titre = ?, description = ?, video_url = ?, role_cible = ?, actif = ? 
		WHERE id = ?`
		
	_, err := Conn.Exec(query, 
		t.Titre, 
		t.Description, 
		t.Video_url, 
		t.Role_cible, 
		t.Actif, 
		t.ID, 
	)
	return err
}

func SupprimerTip(id int) error {
	query := `DELETE FROM TIPS WHERE id = ?`
	_, err := Conn.Exec(query, id)
	return err
}