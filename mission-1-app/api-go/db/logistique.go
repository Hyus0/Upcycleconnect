package db

import (
	"fmt"
	"database/sql"
	"upcycleconnect/api-go/models"
	"crypto/rand"
	"math/big"
)

func generateSecureToken(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

func ReserverUnCasier(annonceID int, siteID int) (string, error) {
	var poidsAnnonce float64
	err := Conn.QueryRow("SELECT poids_estime_kg FROM ANNONCE WHERE id = ?", annonceID).Scan(&poidsAnnonce)
	if err != nil {
		return "", fmt.Errorf("impossible de trouver le poids de l'annonce")
	}

	query := `
		SELECT c.id FROM CASIER c
		JOIN CONTENEUR co ON c.id_conteneur = co.id
		WHERE co.id_site = ? 
		AND c.statut = 'Libre' 
		AND co.statut = 'Operationnel'
		AND (co.niveau_remplissage + ?) <= co.capacite_max_kg
		LIMIT 1`

	var casierID int
	err = Conn.QueryRow(query, siteID, poidsAnnonce).Scan(&casierID)
	if err != nil {
		return "", fmt.Errorf("site complet ou limite de poids atteinte")
	}

	tokenDepot, err := generateSecureToken(20)
	if err != nil {
		return "", err
	}

	tx, err := Conn.Begin()
	if err != nil { return "", err }

	_, err = tx.Exec("UPDATE CASIER SET statut = 'Reserve' WHERE id = ?", casierID)
	if err != nil { tx.Rollback(); return "", err }

	_, err = tx.Exec(`
		UPDATE ANNONCE SET 
			id_casier = ?, id_site = ?, code_barre_depot = ?, statut = 'Reserve' 
		WHERE id = ?`, casierID, siteID, tokenDepot, annonceID)
	if err != nil { tx.Rollback(); return "", err }

	return tokenDepot, tx.Commit()
}

func DeposerObjet(codeBarreDepot string, siteID int) (string, error) {
	var annonceID, conteneurID int
	var poids float64
	
	var acheteurID int
	var titreObjet, nomSite string
	
	err := Conn.QueryRow(`
		SELECT a.id, c.id_conteneur, a.poids_estime_kg, a.id_acheteur, a.titre, s.nom
		FROM ANNONCE a
		JOIN CASIER c ON a.id_casier = c.id
		JOIN SITE s ON a.id_site = s.id
		WHERE a.code_barre_depot = ? 
		AND a.id_site = ? 
		AND a.statut = 'Reserve'`, codeBarreDepot, siteID).Scan(&annonceID, &conteneurID, &poids, &acheteurID, &titreObjet, &nomSite)
	
	if err != nil {
		return "", fmt.Errorf("token invalide ou mauvaise borne") 
	}

	tokenRetrait, err := generateSecureToken(20)
	if err != nil { return "", err }

	tx, err := Conn.Begin()
	if err != nil { return "", err }

	_, err = tx.Exec("UPDATE CONTENEUR SET niveau_remplissage = niveau_remplissage + ? WHERE id = ?", poids, conteneurID)
	if err != nil { 
		tx.Rollback()
		return "", err 
	}

	_, err = tx.Exec(`
		UPDATE ANNONCE SET 
			statut = 'Depose', 
			date_depot_effective = NOW(),
			code_barre_retrait = ?
		WHERE id = ?`, tokenRetrait, annonceID)
	if err != nil { 
		tx.Rollback()
		return "", err 
	}

	titreNotif := "Votre objet est disponible!"
	messageNotif := fmt.Sprintf(
		"Bonne nouvelle ! Le vendeur a déposé votre objet '%s' au site '%s'. Vous pouvez aller le récupérer dès maintenant.", 
		titreObjet, 
		nomSite,
	)
	
	err = CreerNotification(acheteurID, 0, "Casier", titreNotif, messageNotif)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("erreur lors de la notification acheteur : %v", err)
	}

	return tokenRetrait, tx.Commit()
}

const COMMISSION_RATE = 0.05

func AcheterAnnonce(annonceID int, acheteurID int, prixBase float64, stripePaymentID string) (int, string, error) {
	var prenomAcheteur string
	err := Conn.QueryRow("SELECT prenom FROM UTILISATEUR WHERE id = ?", acheteurID).Scan(&prenomAcheteur)
	if err != nil {
		return 0, "", fmt.Errorf("l'acheteur n'existe pas ou ID manquant (reçu: %d)", acheteurID)
	}

	var vendeurID int
	var titreAnnonce string
	err = Conn.QueryRow("SELECT id_vendeur, titre FROM ANNONCE WHERE id = ? AND statut = 'Disponible'", annonceID).Scan(&vendeurID, &titreAnnonce)
	if err != nil {
		return 0, "", fmt.Errorf("cette annonce n'est plus disponible")
	}

	tx, err := Conn.Begin()
	if err != nil {
		return 0, "", err
	}
	defer tx.Rollback()

	query := `
		UPDATE ANNONCE 
		SET id_acheteur = ?, statut = 'Paye', date_achat = NOW() 
		WHERE id = ? AND statut = 'Disponible'`
	res, err := tx.Exec(query, acheteurID, annonceID)
	if err != nil {
		return 0, "", err
	}
	
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return 0, "", fmt.Errorf("l'annonce n'est plus disponible (conflit de concurrence)")
	}

	commission := prixBase * COMMISSION_RATE
	montantTotalTTC := prixBase + commission

	res, err = tx.Exec(`
		INSERT INTO COMMANDE (id_utilisateur, montant_total, statut, date_commande) 
		VALUES (?, ?, 'Payee', NOW())`, 
		acheteurID, montantTotalTTC)
	if err != nil {
		return 0, "", err
	}
	commandeID, _ := res.LastInsertId()

	_, err = tx.Exec(`
		INSERT INTO LIGNE_COMMANDE (id_commande, id_vendeur, type_item, reference_id, prix_unitaire, commission_upc) 
		VALUES (?, ?, 'Annonce', ?, ?, ?)`, 
		commandeID, vendeurID, annonceID, prixBase, commission)
	if err != nil {
		return 0, "", err
	}

	res, err = tx.Exec(`
		INSERT INTO TRANSACTION (id_commande, id_acheteur, montant_total, statut_paiement, date_transaction, stripe_payment_id)
		VALUES (?, ?, ?, 'Valide', NOW(), ?)`,
		commandeID, acheteurID, montantTotalTTC, stripePaymentID)
	if err != nil {
		return 0, "", err
	}
	transactionID, _ := res.LastInsertId()

	factureID, numeroFacture, err := CreateFactureForTransaction(tx, transactionID)
	if err != nil {
		return 0, "", err
	}

	_, err = tx.Exec("UPDATE DM_SALE SET status = 'Payee' WHERE id_annonce = ? AND id_buyer = ?", annonceID, acheteurID)
	if err != nil {
		return 0, "", err
	}

	_, err = tx.Exec("UPDATE DM_SALE SET status = 'Annulee' WHERE id_annonce = ? AND id_buyer != ?", annonceID, acheteurID)
	if err != nil {
		return 0, "", err
	}

	_, err = tx.Exec("UPDATE DM_OFFER SET status = 'Annulee' WHERE id_annonce = ? AND status = 'En attente'", annonceID)
	if err != nil {
		return 0, "", err
	}

	titreNotif := fmt.Sprintf("Votre '%s' a été vendu", titreAnnonce)
	messageNotif := fmt.Sprintf("%s a acheté votre objet. Réservez un casier dès maintenant pour aller le déposer !", prenomAcheteur)
	_, err = tx.Exec(`
		INSERT INTO NOTIFICATION (id_utilisateur, id_emetteur, type, titre, message) 
		VALUES (?, ?, 'Rappel', ?, ?)`, 
		vendeurID, acheteurID, titreNotif, messageNotif)
	if err != nil {
		return 0, "", err
	}

	if err := tx.Commit(); err != nil {
		return 0, "", err
	}

	return factureID, numeroFacture, nil
}

func RecupererObjet(codeBarreRetrait string, siteID int) error {
	var annonceID, idCasier, idConteneur int
	var poids float64
	
	var idVendeur int
	var idAcheteur sql.NullInt64 
	var titreAnnonce string
	var dbSiteID int 

	err := Conn.QueryRow(`
		SELECT a.id, a.poids_estime_kg, a.id_casier, c.id_conteneur, a.id_vendeur, a.id_acheteur, a.titre, a.id_site
		FROM ANNONCE a 
		JOIN CASIER c ON a.id_casier = c.id 
		WHERE a.code_barre_retrait = ? AND a.statut = 'Depose'`, 
		codeBarreRetrait).Scan(&annonceID, &poids, &idCasier, &idConteneur, &idVendeur, &idAcheteur, &titreAnnonce, &dbSiteID)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("code-barres invalide ou objet déjà récupéré")
		}
		return fmt.Errorf("erreur base de données : %v", err)
	}

	if dbSiteID != siteID {
		return fmt.Errorf("vous êtes à la mauvaise borne ! Cet objet vous attend au site n°%d", dbSiteID)
	}

	tx, err := Conn.Begin()
	if err != nil { return err }

	_, err = tx.Exec("UPDATE CONTENEUR SET niveau_remplissage = niveau_remplissage - ? WHERE id = ?", poids, idConteneur)
	if err != nil { tx.Rollback(); return err }

	_, err = tx.Exec("UPDATE CASIER SET statut = 'Libre' WHERE id = ?", idCasier)
	if err != nil { tx.Rollback(); return err }

	_, err = tx.Exec(`
		UPDATE ANNONCE SET 
			statut = 'Recupere', 
			date_recuperation_effective = NOW(),
			code_barre_retrait = NULL,
			code_barre_depot = NULL
		WHERE id = ?`, annonceID)
	if err != nil { tx.Rollback(); return err }

	co2Evite := poids * 1.5            
	ressourcesEco := poids * 0.8       
	pointsGagnes := 25 + int(poids*5)
	
	_, err = tx.Exec(`
		UPDATE UPCYCLING_SCORE 
		SET ressources_economisees = ressources_economisees + ?,
		    co2_total_evite_kg = co2_total_evite_kg + ?,
		    nb_objets_recycles = nb_objets_recycles + 1,
		    total_points = total_points + ?
		WHERE id_utilisateur = ?
	`, ressourcesEco, co2Evite, pointsGagnes, idVendeur)

	if err != nil { 
		tx.Rollback()
		return fmt.Errorf("erreur lors de la mise à jour du score vendeur: %v", err) 
	}

	if idAcheteur.Valid && idAcheteur.Int64 > 0 {
		_, err = tx.Exec(`
			UPDATE UPCYCLING_SCORE 
			SET ressources_economisees = ressources_economisees + ?,
			    co2_total_evite_kg = co2_total_evite_kg + ?,
			    nb_objets_recycles = nb_objets_recycles + 1,
			    total_points = total_points + ?
			WHERE id_utilisateur = ?
		`, ressourcesEco, co2Evite, pointsGagnes, idAcheteur.Int64)

		if err != nil { 
			tx.Rollback()
			return fmt.Errorf("erreur lors de la mise à jour du score acheteur: %v", err) 
		}
	}

	titreNotifVendeur := fmt.Sprintf("Votre '%s' a été récupéré!", titreAnnonce)
	messageNotifVendeur := "L'acheteur a récupéré votre objet au point de collecte. Votre Upcycling Score a été mis à jour, merci pour votre geste éco-responsable!"

	_, err = tx.Exec(`
		INSERT INTO NOTIFICATION (id_utilisateur, id_emetteur, type, titre, message) 
		VALUES (?, 0, 'Alerte', ?, ?)`, 
		idVendeur, titreNotifVendeur, messageNotifVendeur)
	
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erreur lors de la création de la notification vendeur : %v", err)
	}

	if idAcheteur.Valid && idAcheteur.Int64 > 0 {
		titreNotifAcheteur := fmt.Sprintf("Vous avez récupéré '%s' !", titreAnnonce)
		messageNotifAcheteur := "Félicitations pour cette acquisition ! Vous avez évité l'achat de neuf et gagné des points pour votre Upcycling Score."

		_, err = tx.Exec(`
			INSERT INTO NOTIFICATION (id_utilisateur, id_emetteur, type, titre, message) 
			VALUES (?, 0, 'Alerte', ?, ?)`, 
			idAcheteur.Int64, titreNotifAcheteur, messageNotifAcheteur)
		
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("erreur lors de la création de la notification acheteur : %v", err)
		}
	}

	return tx.Commit()
}

func GetAllSites() ([]models.Site, error) {
    if Conn == nil {
        return nil, fmt.Errorf("DB non connectée")
    }

    query := `SELECT id, nom, ville, code_postal, adresse, COALESCE(telephone, ''), type, actif FROM SITE`
    
    rows, err := Conn.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    sites := []models.Site{}

    for rows.Next() {
        var s models.Site
        err := rows.Scan(
            &s.ID, 
            &s.Nom, 
            &s.Ville, 
            &s.CodePostal, 
            &s.Adresse, 
            &s.Telephone, 
            &s.Type, 
            &s.Actif,
        )
        if err != nil {
            fmt.Println("Erreur lors du Scan du site:", err)
            return nil, err
        }
        sites = append(sites, s)
    }

    return sites, nil
}

func GetSiteByID(id int) (map[string]string, error) {
    query := `SELECT nom, ville, code_postal, adresse, COALESCE(telephone, "") FROM SITE WHERE id = ?`
    
    var nom, ville, codePostal, adresse, telephone string
    err := Conn.QueryRow(query, id).Scan(&nom, &ville, &codePostal, &adresse, &telephone)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("site introuvable")
        }
        return nil, err
    }

    return map[string]string{
        "nom":         nom,
        "ville":       ville,
        "code_postal": codePostal,
        "adresse":     adresse,
        "telephone":   telephone,
    }, nil
}

func GetConteneursBySite(siteID int) ([]models.Conteneur, error) {
	query := `SELECT id, id_site, type_dechet, statut, capacite_max_kg, niveau_remplissage 
	          FROM CONTENEUR WHERE id_site = ?`
	rows, err := Conn.Query(query, siteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conteneurs []models.Conteneur
	for rows.Next() {
		var c models.Conteneur
		if err := rows.Scan(&c.ID, &c.IdSite, &c.TypeDechet, &c.Statut, &c.CapaciteMaxKg, &c.NiveauRemplissage); err != nil {
			return nil, err
		}
		conteneurs = append(conteneurs, c)
	}
	return conteneurs, nil
}

func RetireObjetCasier(idAnnonce int) error {
    var poids float64
    var idCasier int
    var idConteneur int

    err := Conn.QueryRow(`
        SELECT a.poids_estime_kg, a.id_casier, c.id_conteneur 
        FROM ANNONCE a 
        JOIN CASIER c ON a.id_casier = c.id 
        WHERE a.id = ?`, idAnnonce).Scan(&poids, &idCasier, &idConteneur)

    if err != nil {
        return err
    }

    tx, err := Conn.Begin()
    if err != nil {
        return err
    }

    _, err = tx.Exec("UPDATE CONTENEUR SET niveau_remplissage = niveau_remplissage - ? WHERE id = ?", poids, idConteneur)
    if err != nil {
        tx.Rollback()
        return err
    }

    _, err = tx.Exec("UPDATE CASIER SET statut = 'Libre' WHERE id = ?", idCasier)
    if err != nil {
        tx.Rollback()
        return err
    }

    _, err = tx.Exec(`
        UPDATE ANNONCE SET 
            statut = 'Disponible', 
            id_casier = NULL, 
            id_site = NULL, 
            code_barre_depot = NULL 
        WHERE id = ?`, idAnnonce)
    if err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit()
}