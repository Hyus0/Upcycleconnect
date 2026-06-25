package db

import (
	"upcycleconnect/api-go/models"
)

func GetLangues() ([]models.Langue, error) {
	rows, err := Conn.Query("SELECT id, code, nom_langue FROM LANGUE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var langues []models.Langue
	for rows.Next() {
		var l models.Langue
		if err := rows.Scan(&l.ID, &l.Code, &l.NomLangue); err != nil {
			return nil, err
		}
		langues = append(langues, l)
	}
	return langues, nil
}

func GetTraductionsByCode(code string) (map[string]string, error) {
	rows, err := Conn.Query(`
		SELECT t.cle_traduction, t.text_traduit 
		FROM TRADUCTION t
		JOIN LANGUE l ON t.id_langue = l.id
		WHERE l.code = ?`, code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	traductions := make(map[string]string)

	for rows.Next() {
		var cle, texte string
		if err := rows.Scan(&cle, &texte); err != nil {
			return nil, err
		}
		traductions[cle] = texte
	}
	
	return traductions, nil
}

func UpdateLangueUtilisateur(userID int, langueID int) error {
	_, err := Conn.Exec("UPDATE UTILISATEUR SET id_langue = ? WHERE id = ?", langueID, userID)
	return err
}