package db

import (
	"fmt"
	"upcycleconnect/api-go/models"
)

func GetPanierByUserId(userID int) ([]models.PanierItem, error) {
	query := `SELECT id, id_utilisateur, type_item, reference_id, prix_unitaire, date_ajout
	          FROM PANIER_ITEM WHERE id_utilisateur = ? ORDER BY date_ajout DESC`

	rows, err := Conn.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []models.PanierItem{}
	for rows.Next() {
		var item models.PanierItem
		if err := rows.Scan(&item.ID, &item.IdUtilisateur, &item.TypeItem, &item.ReferenceID, &item.PrixUnitaire, &item.DateAjout); err == nil {
			items = append(items, item)
		}
	}
	return items, nil
}

func AddToPanier(userID int, typeItem string, refID int, prix float64) error {
	_, err := Conn.Exec(`INSERT INTO PANIER_ITEM (id_utilisateur, type_item, reference_id, prix_unitaire)
	    VALUES (?, ?, ?, ?)`, userID, typeItem, refID, prix)
	return err
}

func RemoveFromPanier(itemID int, userID int) error {
	result, err := Conn.Exec("DELETE FROM PANIER_ITEM WHERE id = ? AND id_utilisateur = ?", itemID, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("aucun article trouvé avec cet ID pour cet utilisateur")
	}

	return nil
}

func Checkout(userID int) (int, error) {
	tx, err := Conn.Begin()
	if err != nil {
		return 0, err
	}

	rows, err := tx.Query("SELECT id, type_item, reference_id, prix_unitaire FROM PANIER_ITEM WHERE id_utilisateur = ?", userID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var items []models.PanierItem
	var montantTotal float64 = 0

	for rows.Next() {
		var item models.PanierItem
		if err := rows.Scan(&item.ID, &item.TypeItem, &item.ReferenceID, &item.PrixUnitaire); err == nil {
			items = append(items, item)
			montantTotal += item.PrixUnitaire
		}
	}
	rows.Close()

	if len(items) == 0 {
		tx.Rollback()
		return 0, fmt.Errorf("le panier est vide")
	}

	res, err := tx.Exec("INSERT INTO COMMANDE (id_utilisateur, montant_total, statut) VALUES (?, ?, 'Payee')", userID, montantTotal)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	commandeID, _ := res.LastInsertId()

	for _, item := range items {
		var commission float64 = 0
		if item.TypeItem == "Annonce" {
			commission = item.PrixUnitaire * 0.10
		}

		_, err = tx.Exec(`INSERT INTO LIGNE_COMMANDE (id_commande, type_item, reference_id, prix_unitaire, commission_upc)
		                  VALUES (?, ?, ?, ?, ?)`, commandeID, item.TypeItem, item.ReferenceID, item.PrixUnitaire, commission)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		if item.TypeItem == "Formation" {
			_, err = tx.Exec("INSERT IGNORE INTO FORMATION_INSCRIPTION (id_utilisateur, id_formation) VALUES (?, ?)", userID, item.ReferenceID)
			if err != nil {
				tx.Rollback()
				return 0, fmt.Errorf("erreur lors de l'inscription à la formation : %v", err)
			}
		}
	}

	_, err = tx.Exec(`INSERT INTO TRANSACTION (id_acheteur, id_commande, montant_total, statut_paiement)
	                  VALUES (?, ?, ?, 'Valide')`, userID, commandeID, montantTotal)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	_, err = tx.Exec("DELETE FROM PANIER_ITEM WHERE id_utilisateur = ?", userID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return int(commandeID), nil
}
