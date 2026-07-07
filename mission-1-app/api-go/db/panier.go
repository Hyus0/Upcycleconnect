package db

import (
	"fmt"
	"upcycleconnect/api-go/models"
)

type CheckoutResult struct {
    CommandeID    int  `json:"commande_id"`
    TransactionID int  `json:"transaction_id"`
    FactureID     int  `json:"facture_id"`
    NumeroFacture string `json:"numero_facture"`
}

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
	var count int
	err := Conn.QueryRow(`
		SELECT COUNT(*)
		FROM PANIER_ITEM
		WHERE id_utilisateur = ? AND type_item = ? AND reference_id = ?
	`, userID, typeItem, refID).Scan(&count)

	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("cet article est déjà dans votre panier")
	}

	_, err = Conn.Exec(`
		INSERT INTO PANIER_ITEM (id_utilisateur, type_item, reference_id, prix_unitaire)
	    VALUES (?, ?, ?, ?)
	`, userID, typeItem, refID, prix)

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
		return fmt.Errorf("aucun article trouve avec cet ID pour cet utilisateur")
	}

	return nil
}

func Checkout(userID int, stripePaymentID string) (*CheckoutResult, error) {
	tx, err := Conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() 

	rows, err := tx.Query("SELECT id, type_item, reference_id, prix_unitaire FROM PANIER_ITEM WHERE id_utilisateur = ?", userID)
	if err != nil {
		return nil, err
	}

	var items []models.PanierItem
	var montantTotalTTC float64

	for rows.Next() {
		var item models.PanierItem
		if err := rows.Scan(&item.ID, &item.TypeItem, &item.ReferenceID, &item.PrixUnitaire); err == nil {
			items = append(items, item)
			
			prixTTC := item.PrixUnitaire
			if item.TypeItem == "Annonce" || item.TypeItem == "Projet" {
				prixTTC += item.PrixUnitaire * 0.05
			}
			montantTotalTTC += prixTTC
		}
	}
	rows.Close()

	if len(items) == 0 {
		return nil, fmt.Errorf("le panier est vide")
	}

	res, err := tx.Exec("INSERT INTO COMMANDE (id_utilisateur, montant_total, statut, date_commande) VALUES (?, ?, 'Payee', NOW())", userID, montantTotalTTC)
	if err != nil {
		return nil, err
	}
	commandeID, _ := res.LastInsertId()

	for _, item := range items {
		var commission float64
		if item.TypeItem == "Annonce" || item.TypeItem == "Projet" {
			commission = item.PrixUnitaire * 0.05
		}

		_, err = tx.Exec(`
			INSERT INTO LIGNE_COMMANDE (id_commande, type_item, reference_id, prix_unitaire, commission_upc)
			VALUES (?, ?, ?, ?, ?)`, 
			commandeID, item.TypeItem, item.ReferenceID, item.PrixUnitaire, commission)
		if err != nil {
			return nil, err
		}

		if item.TypeItem == "Formation" {
			_, err = tx.Exec("INSERT IGNORE INTO FORMATION_INSCRIPTION (id_utilisateur, id_formation) VALUES (?, ?)", userID, item.ReferenceID)
			if err != nil {
				return nil, fmt.Errorf("erreur lors de l'inscription à la formation : %v", err)
			}
		}

		if item.TypeItem == "Abonnement" {
			_, err = tx.Exec(`
				INSERT INTO ABONNEMENT (id_acheteur, id_type_abonnement, date_debut, date_fin, statut, stripe_subscription_id)
				SELECT ?, id, NOW(), DATE_ADD(NOW(), INTERVAL duree_mois MONTH), 'Actif', ?
				FROM TYPE_ABONNEMENT
				WHERE nom = 'DM Plus'
			`, userID, stripePaymentID)
			if err != nil {
				return nil, fmt.Errorf("erreur lors de l'activation de l'abonnement : %v", err)
			}
		}

		if item.TypeItem == "Projet" {
			_, err = tx.Exec(`
				UPDATE PROJET_UPCYCLING 
				SET statut = 'Vendu', id_acheteur = ?, date_achat = NOW() 
				WHERE id = ?`, userID, item.ReferenceID)
			if err == nil {
				_, err = tx.Exec("UPDATE DM_SALE SET status = 'Payee' WHERE id_projet = ? AND id_buyer = ?", item.ReferenceID, userID)
			}
			if err != nil {
				return nil, fmt.Errorf("erreur lors de la validation du projet : %v", err)
			}
		}

		if item.TypeItem == "Annonce" {
			_, err = tx.Exec(`
				UPDATE ANNONCE 
				SET statut = 'Paye', id_acheteur = ?, date_achat = NOW() 
				WHERE id = ?`, userID, item.ReferenceID)
			if err == nil {
				_, err = tx.Exec("UPDATE DM_SALE SET status = 'Payee' WHERE id_annonce = ? AND id_buyer = ?", item.ReferenceID, userID)
			}
			if err != nil {
				return nil, fmt.Errorf("erreur lors de la validation de l'annonce : %v", err)
			}
		}
	}

	transactionRes, err := tx.Exec(`
		INSERT INTO `+"`TRANSACTION`"+` (id_acheteur, id_commande, montant_total, statut_paiement, stripe_payment_id, date_transaction)
		VALUES (?, ?, ?, 'Valide', ?, NOW())`, 
		userID, commandeID, montantTotalTTC, stripePaymentID)
	if err != nil {
		return nil, err
	}
	transactionID, _ := transactionRes.LastInsertId()

	factureID, numeroFacture, err := CreateFactureForTransaction(tx, transactionID)

	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("DELETE FROM PANIER_ITEM WHERE id_utilisateur = ?", userID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &CheckoutResult{
		CommandeID:    int(commandeID),    
		TransactionID: int(transactionID), 
		FactureID:     factureID,          	
		NumeroFacture: numeroFacture,
	}, nil
}