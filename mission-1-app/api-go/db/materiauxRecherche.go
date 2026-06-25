package db

import (
	"fmt"
	"strings"
	"upcycleconnect/api-go/models"
)

func GetMateriauRecherche(userID int) (string, error) {
	var materiaux string
	err := Conn.QueryRow(`
		SELECT COALESCE(materiaux_recherches, '')
		FROM UTILISATEUR
		WHERE id = ?
	`, userID).Scan(&materiaux)
	if err != nil {
		return "", fmt.Errorf("GetMateriauRecherche: %v", err)
	}
	return materiaux, nil
}

func UpdateMateriauRecherche(userID int, materiaux string) error {
	res, err := Conn.Exec(`
		UPDATE UTILISATEUR
		SET materiaux_recherches = ?
		WHERE id = ?
	`, materiaux, userID)
	if err != nil {
		return fmt.Errorf("UpdateMateriauRecherche: %v", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("UpdateMateriauRecherche rows affected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("UpdateMateriauRecherche: utilisateur introuvable")
	}
	return nil
}

func RechercheAnnonceMateriau(materiaux []string) ([]models.Annonce, error) {
	if len(materiaux) == 0 {
		return []models.Annonce{}, nil
	}

	query := `
		SELECT id, id_vendeur, titre, description, type_materiau, type, ville, prix, statut, est_valide, date_creation
		FROM ANNONCE
		WHERE est_valide = 'Valide' AND statut = 'Disponible' AND (
	`

	var conditions []string
	var args []interface{}

	for _, mat := range materiaux {
		conditions = append(conditions, "(type_materiau LIKE ? OR titre LIKE ?)")
		motCle := "%" + strings.TrimSpace(mat) + "%"
		args = append(args, motCle, motCle)
	}

	query += strings.Join(conditions, " OR ") + ")"
	query += " ORDER BY date_creation DESC"

	rows, err := Conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("RechercheAnnonceMateriau query: %v", err)
	}
	defer rows.Close()

	var annonces []models.Annonce
	for rows.Next() {
		var a models.Annonce
		err := rows.Scan(
			&a.ID,
			&a.IdVendeur,
			&a.Titre,
			&a.Description,
			&a.TypeMateriau,
			&a.Type,
			&a.Ville,
			&a.Prix,
			&a.Statut,
			&a.EstValide,
			&a.DateCreation,
		)
		if err != nil {
			return nil, fmt.Errorf("RechercheAnnonceMateriau scan: %v", err)
		}
		annonces = append(annonces, a)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("RechercheAnnonceMateriau rows error: %v", err)
	}

	return annonces, nil
}