package db

import (
	"database/sql"
	"fmt"
)

type SubscriptionStatus struct {
	IsPremium bool   `json:"is_premium"`
	PlanName  string `json:"plan_name,omitempty"`
	DateFin   string `json:"date_fin,omitempty"`
	Statut    string `json:"statut,omitempty"`
	PlanID    int    `json:"plan_id,omitempty"` 
}

type TypeAbonnement struct {
	ID          int     `json:"id"`
	Nom         string  `json:"nom"`
	Description string  `json:"description"`
	PrixHT      float64 `json:"prix_ht"`
	DureeMois   int     `json:"duree_mois"`
}

func GetTypeAbonnementByID(planID int) (TypeAbonnement, error) {
	var p TypeAbonnement
	err := Conn.QueryRow("SELECT id, nom, description, prix_ht, duree_mois FROM TYPE_ABONNEMENT WHERE id = ?", planID).
		Scan(&p.ID, &p.Nom, &p.Description, &p.PrixHT, &p.DureeMois)
	
	if err != nil {
		return p, err
	}
	return p, nil
}


func GetUserSubscription(userID int) (SubscriptionStatus, error) {
	var sub SubscriptionStatus

	err := Conn.QueryRow(`
		SELECT t.id, t.nom, DATE_FORMAT(a.date_fin, '%Y-%m-%dT%H:%i:%sZ'), a.statut
		FROM ABONNEMENT a
		JOIN TYPE_ABONNEMENT t ON t.id = a.id_type_abonnement
		WHERE a.id_acheteur = ? AND a.statut = 'Actif' AND a.date_fin >= NOW()
		ORDER BY t.id DESC, a.date_fin DESC LIMIT 1
	`, userID).Scan(&sub.PlanID, &sub.PlanName, &sub.DateFin, &sub.Statut)

	if err == sql.ErrNoRows {
		return SubscriptionStatus{IsPremium: false}, nil
	} else if err != nil {
		return sub, err
	}

	sub.IsPremium = true
	return sub, nil
}

func SubscribeUser(userID int, planID int, stripePaymentID string) error {
	active, _ := GetUserSubscription(userID)
	
	if active.IsPremium && active.PlanID == planID {
		return fmt.Errorf("utilisateur deja abonne a ce plan")
	}

	tx, err := Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() 

	if active.IsPremium && active.PlanID != planID {
		_, err = tx.Exec(`
			UPDATE ABONNEMENT 
			SET statut = 'Resilie' 
			WHERE id_acheteur = ? AND statut = 'Actif'
		`, userID)
		if err != nil {
			return fmt.Errorf("impossible de résilier l'ancien abonnement: %v", err)
		}
	}

	var prixPlan float64
	err = tx.QueryRow("SELECT prix_ht FROM TYPE_ABONNEMENT WHERE id = ?", planID).Scan(&prixPlan)
	if err != nil {
		return fmt.Errorf("type d'abonnement introuvable")
	}

	_, err = tx.Exec(`
		INSERT INTO ABONNEMENT (id_acheteur, id_type_abonnement, date_debut, date_fin, statut, stripe_subscription_id)
		SELECT ?, id, NOW(), DATE_ADD(NOW(), INTERVAL duree_mois MONTH), 'Actif', ?
		FROM TYPE_ABONNEMENT
		WHERE id = ?
	`, userID, stripePaymentID, planID)
	if err != nil {
		return err
	}

	res, err := tx.Exec("INSERT INTO COMMANDE (id_utilisateur, montant_total, statut, date_commande) VALUES (?, ?, 'Payee', NOW())", userID, prixPlan)
	if err != nil {
		return err
	}
	commandeID, _ := res.LastInsertId()

	_, err = tx.Exec(`INSERT INTO LIGNE_COMMANDE (id_commande, type_item, reference_id, prix_unitaire, commission_upc)
					  VALUES (?, 'Abonnement', ?, ?, 0)`, commandeID, planID, prixPlan)
	if err != nil {
		return err
	}

	transactionRes, err := tx.Exec(`INSERT INTO TRANSACTION (id_acheteur, id_commande, montant_total, statut_paiement, date_transaction, stripe_payment_id)
					  VALUES (?, ?, ?, 'Valide', NOW(), ?)`, userID, commandeID, prixPlan, stripePaymentID)
	if err != nil {
		return err
	}
	transactionID, _ := transactionRes.LastInsertId()

	_, _, err = CreateFactureForTransaction(tx, transactionID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func CancelSubscription(userID int) error {
	_, err := Conn.Exec(`
		UPDATE ABONNEMENT 
		SET statut = 'Resilie' 
		WHERE id_acheteur = ? AND statut = 'Actif'
	`, userID)
	
	return err
}