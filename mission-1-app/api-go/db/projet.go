package db

import (
	"fmt"
	"database/sql"
	"upcycleconnect/api-go/models"
)

func GetAllProjets() ([]models.ProjetUpcycling, error) {
	if Conn == nil {
		return nil, fmt.Errorf("connexion DB non initialisee")
	}

	query := `
		SELECT	
			p.id,
			p.id_createur,
			p.image_url,
			p.titre,
			p.description_courte,
			p.date_creation,
			p.score_impact,
			p.nb_vues,
			p.nb_likes,
			p.co2_evite_kg,
			p.visible_public,
			COUNT(pi.id_utilisateur) AS nb_participants
		FROM PROJET_UPCYCLING p
		LEFT JOIN PROJET_INSCRIPTION pi ON pi.id_projet = p.id
		WHERE p.visible_public = 1
		GROUP BY
			p.id, p.id_createur, p.image_url, p.titre, p.description_courte, p.date_creation,
			p.score_impact, p.nb_vues, p.nb_likes, p.co2_evite_kg, p.visible_public
	`

	rows, err := Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projets []models.ProjetUpcycling
	for rows.Next() {
		var p models.ProjetUpcycling
		var nbParticipants int
		err := rows.Scan(
			&p.ID,
			&p.IdCreateur,
			&p.ImageUrl,
			&p.Titre,
			&p.DescriptionCourte,
			&p.DateCreation,
			&p.ScoreImpact,
			&p.NbVues,
			&p.NbLikes,
			&p.Co2EviteKg,
			&p.VisiblePublic,
			&nbParticipants,
		)
		if err != nil {
			fmt.Println("Erreur Scan Projet:", err)
			return nil, err
		}
		projets = append(projets, p)
	}

	return projets, nil
}

func UpdateProjet(projetID int, p models.ProjetUpcycling) error {
	tx, err := Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
 
	var prixValue interface{}
	if p.Prix != nil && *p.Prix > 0 {
		prixValue = *p.Prix
	} else {
		prixValue = nil
	}
 
	_, err = tx.Exec(`
		UPDATE PROJET_UPCYCLING
		SET titre=?, description_courte=?, visible_public=?, co2_evite_kg=?, prix=?
		WHERE id=?
	`, p.Titre, p.DescriptionCourte, p.VisiblePublic, p.Co2EviteKg, prixValue, projetID)
	if err != nil {
		return err
	}
 
	_, err = tx.Exec("DELETE FROM ETAPE WHERE id_projet = ?", projetID)
	if err != nil {
		return err
	}
 
	for _, e := range p.Etapes {
		_, err = tx.Exec(`INSERT INTO ETAPE (id_projet, numero_ordre, titre, description, image_url)
                          VALUES (?, ?, ?, ?, ?)`, projetID, e.NumeroOrdre, e.Titre, e.Description, e.ImageUrl)
		if err != nil {
			return err
		}
	}
 
	return tx.Commit()
}

func GetProjet(projetID int, userID int) (*models.ProjetUpcycling, error) {
	if Conn == nil {
		return nil, fmt.Errorf("connexion DB non initialisee")
	}
 
	query := `
		SELECT
			id, id_createur, image_url, titre, description_courte, date_creation,
			score_impact, nb_vues, nb_likes, co2_evite_kg, visible_public,
			prix, COALESCE(statut, 'Disponible'), id_acheteur
		FROM PROJET_UPCYCLING
		WHERE id = ?`
 
	var p models.ProjetUpcycling
	var prix sql.NullFloat64
	var idAcheteur sql.NullInt64
 
	err := Conn.QueryRow(query, projetID).Scan(
		&p.ID, &p.IdCreateur, &p.ImageUrl, &p.Titre, &p.DescriptionCourte, &p.DateCreation,
		&p.ScoreImpact, &p.NbVues, &p.NbLikes, &p.Co2EviteKg, &p.VisiblePublic,
		&prix, &p.Statut, &idAcheteur,
	)
	if err != nil {
		return nil, err
	}
 
	if prix.Valid {
		v := prix.Float64
		p.Prix = &v
		p.EnVente = true
	}
	if idAcheteur.Valid {
		v := int(idAcheteur.Int64)
		p.IdAcheteur = &v
	}
 
	etapeQuery := `SELECT id, numero_ordre, titre, description, image_url FROM ETAPE WHERE id_projet = ? ORDER BY numero_ordre`
	rows, err := Conn.Query(etapeQuery, projetID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var e models.EtapeProjet
			rows.Scan(&e.ID, &e.NumeroOrdre, &e.Titre, &e.Description, &e.ImageUrl)
			p.Etapes = append(p.Etapes, e)
		}
	}
 
	var count int
	checkQuery := "SELECT COUNT(*) FROM PROJET_INSCRIPTION WHERE id_utilisateur = ? AND id_projet = ?"
	Conn.QueryRow(checkQuery, userID, projetID).Scan(&count)
 
	return &p, nil
}

func CreateProjet(p models.ProjetUpcycling) (int, error) {
	if Conn == nil {
		return 0, fmt.Errorf("connexion DB non initialisee")
	}
 
	var prixValue interface{}
	statutValue := "Disponible"
	if p.Prix != nil && *p.Prix > 0 {
		prixValue = *p.Prix
	} else {
		prixValue = nil
		statutValue = "Disponible" 
	}
 
	query := `INSERT INTO PROJET_UPCYCLING (
		id_createur, titre, description_courte, image_url, score_impact,
		co2_evite_kg, visible_public, prix, statut
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
 
	result, err := Conn.Exec(
		query,
		p.IdCreateur,
		p.Titre,
		p.DescriptionCourte,
		p.ImageUrl,
		p.ScoreImpact,
		p.Co2EviteKg,
		p.VisiblePublic,
		prixValue,
		statutValue,
	)
	if err != nil {
		return 0, fmt.Errorf("CreateProjet error: %v", err)
	}
 
	lastID, _ := result.LastInsertId()
	return int(lastID), nil
}

func CreateEtape(e models.EtapeProjet) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisee")
	}

	query := `INSERT INTO ETAPE (
		id_projet, numero_ordre, titre, description, image_url
	) VALUES (?, ?, ?, ?, ?)`

	_, err := Conn.Exec(
		query,
		e.IdProjet,
		e.NumeroOrdre,
		e.Titre,
		e.Description,
		e.ImageUrl,
	)
	return err
}

func ModifyProjet(id int, p models.ProjetUpcycling) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisee")
	}
 
	var prixValue interface{}
	if p.Prix != nil && *p.Prix > 0 {
		prixValue = *p.Prix
	} else {
		prixValue = nil
	}
 
	query := `
		UPDATE PROJET_UPCYCLING SET
			titre = ?,
			description_courte = ?,
			visible_public = ?,
			co2_evite_kg = ?,
			prix = ?
		WHERE id = ?
	`
 
	result, err := Conn.Exec(query, p.Titre, p.DescriptionCourte, p.VisiblePublic, p.Co2EviteKg, prixValue, id)
	if err != nil {
		return fmt.Errorf("ModifyProjet: %v", err)
	}
 
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("aucun projet trouvé avec l'ID %d", id)
	}
	return nil
}

func DeleteProjet(id int) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisee")
	}

	_, err := Conn.Exec("DELETE FROM PROJET_UPCYCLING WHERE id = ?", id)
	return err
}

func JoinProjet(userID int, projetID int) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisee")
	}

	queryInsert := "INSERT INTO PROJET_INSCRIPTION (id_utilisateur, id_projet) VALUES (?, ?)"
	if _, err := Conn.Exec(queryInsert, userID, projetID); err != nil {
		return fmt.Errorf("erreur inscription projet: %v", err)
	}

	return nil
}

func QuitProjet(userID int, projetID int) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisee")
	}

	query := "DELETE FROM PROJET_INSCRIPTION WHERE id_utilisateur = ? AND id_projet = ?"
	_, err := Conn.Exec(query, userID, projetID)
	return err
}

func IncrementLike(userID int, projetID int) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisee")
	}
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM PROJET_LIKE WHERE id_utilisateur = ? AND id_projet = ?)"
	err := Conn.QueryRow(checkQuery, userID, projetID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("erreur lors de la vérification du like : %v", err)
	}

	if exists {
		_, err = Conn.Exec("DELETE FROM PROJET_LIKE WHERE id_utilisateur = ? AND id_projet = ?", userID, projetID)
		if err != nil {
			return fmt.Errorf("erreur lors de la suppression du like : %v", err)
		}

		_, err = Conn.Exec("UPDATE PROJET_UPCYCLING SET nb_likes = nb_likes - 1 WHERE id = ?", projetID)
		if err != nil {
			return fmt.Errorf("erreur lors de la mise à jour du compteur (decrement) : %v", err)
		}
	} else {
		_, err = Conn.Exec("INSERT INTO PROJET_LIKE (id_utilisateur, id_projet) VALUES (?, ?)", userID, projetID)
		if err != nil {
			return fmt.Errorf("erreur lors de l'ajout du like : %v", err)
		}

		_, err = Conn.Exec("UPDATE PROJET_UPCYCLING SET nb_likes = nb_likes + 1 WHERE id = ?", projetID)
		if err != nil {
			return fmt.Errorf("erreur lors de la mise à jour du compteur (increment) : %v", err)
		}

		var idCreateur int
		var titreProjet string
		
		errInfos := Conn.QueryRow("SELECT id_createur, titre FROM PROJET_UPCYCLING WHERE id = ?", projetID).Scan(&idCreateur, &titreProjet)
		
		if errInfos == nil {
			if idCreateur != userID {
				
				var prenomLiker string
				errUser := Conn.QueryRow("SELECT prenom FROM UTILISATEUR WHERE id = ?", userID).Scan(&prenomLiker)
				
				if errUser != nil || prenomLiker == "" {
					prenomLiker = "Un membre"
				}

				titreNotif := "Nouveau coup de cœur"
				messageNotif := fmt.Sprintf("%s a eu un coup de cœur pour votre projet '%s' !", prenomLiker, titreProjet)
				
				CreerNotification(idCreateur, userID, "Like", titreNotif, messageNotif)
			}
		}
	}
	return nil
}

func CheckUserLike(userID int, projetID int) (bool, error) {
	if Conn == nil {
		return false, fmt.Errorf("connexion DB non initialisee")
	}

	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM PROJET_LIKE WHERE id_utilisateur = ? AND id_projet = ?)"
	err := Conn.QueryRow(query, userID, projetID).Scan(&exists)

	return exists, err
}

func IncrementVue(projetID int, userID int, ipAddress string) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisee")
	}

	var exists bool
	checkQuery := `
		SELECT EXISTS(
			SELECT 1 FROM PROJET_VUE
			WHERE id_projet = ?
			AND (id_utilisateur = ? OR ip_adresse = ?)
			AND date_vue > DATE_SUB(NOW(), INTERVAL 1 DAY)
		)`

	err := Conn.QueryRow(checkQuery, projetID, userID, ipAddress).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}
	_, err = Conn.Exec("INSERT INTO PROJET_VUE (id_projet, id_utilisateur, ip_adresse) VALUES (?, ?, ?)",
		projetID, userID, ipAddress)
	if err != nil {
		return err
	}

	_, err = Conn.Exec("UPDATE PROJET_UPCYCLING SET nb_vues = nb_vues + 1 WHERE id = ?", projetID)

	return err
}

func GetUserProjetsLike(userID int) ([]models.ProjetUpcycling, error) {
	if Conn == nil {
		return nil, fmt.Errorf("connexion DB non initialisee")
	}
 
	query := `
		SELECT
			p.id,
			p.id_createur,
			p.image_url,
			p.titre,
			p.description_courte,
			p.date_creation,
			p.score_impact,
			p.nb_vues,
			p.nb_likes,
			p.co2_evite_kg,
			p.visible_public
		FROM PROJET_UPCYCLING p
		INNER JOIN PROJET_LIKE pl ON pl.id_projet = p.id
		WHERE pl.id_utilisateur = ?
		ORDER BY pl.date_like DESC
	`
 
	rows, err := Conn.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("GetUserProjetsLike: %v", err)
	}
	defer rows.Close()
 
	projets := []models.ProjetUpcycling{}
	for rows.Next() {
		var p models.ProjetUpcycling
		err := rows.Scan(
			&p.ID,
			&p.IdCreateur,
			&p.ImageUrl,
			&p.Titre,
			&p.DescriptionCourte,
			&p.DateCreation,
			&p.ScoreImpact,
			&p.NbVues,
			&p.NbLikes,
			&p.Co2EviteKg,
			&p.VisiblePublic,
		)
		if err != nil {
			fmt.Println("Erreur Scan GetUserProjetsLike:", err)
			return nil, err
		}
		projets = append(projets, p)
	}
 
	return projets, nil
}


func GetProjetsByUserId(userID int) ([]models.ProjetUpcycling, error) {
	query := `SELECT id, id_createur, titre, description_courte, image_url, score_impact, co2_evite_kg, nb_vues, nb_likes, visible_public, date_creation 
              FROM PROJET_UPCYCLING 
			  WHERE id_createur = ? 
			  ORDER BY date_creation DESC`
	
    rows, err := Conn.Query(query, userID)
	if err != nil { 
		return nil, err 
	}
	defer rows.Close()

	var projets []models.ProjetUpcycling
	for rows.Next() {
		var p models.ProjetUpcycling
		err := rows.Scan(&p.ID, &p.IdCreateur, &p.Titre, &p.DescriptionCourte, &p.ImageUrl, &p.ScoreImpact, &p.Co2EviteKg, &p.NbVues, &p.NbLikes, &p.VisiblePublic, &p.DateCreation)
		if err == nil {
			projets = append(projets, p)
		}
	}
	return projets, nil
}

func CheckProjectLimit(userID int) error {
	sub, err := GetUserSubscription(userID)
	if err != nil {
		return fmt.Errorf("erreur vérification abonnement: %v", err)
	}

	var activeProjectsCount int
	query := `
		SELECT COUNT(*) 
		FROM PROJET_UPCYCLING 
		WHERE id_createur = ? AND visible_public = 1
	`
	err = Conn.QueryRow(query, userID).Scan(&activeProjectsCount)
	if err != nil {
		return fmt.Errorf("erreur comptage projets: %v", err)
	}

	limite := 2 
	if sub.IsPremium {
		limite = 3 
	}

	if activeProjectsCount >= limite {
		if sub.IsPremium {
			return fmt.Errorf("Limite atteinte de 20 projets en ligne") 
		}
		return fmt.Errorf("Limite atteinte de 5 projets en ligne") 
	}

	return nil
}

func AcheterProjet(projetID int, acheteurID int, prixBase float64, stripePaymentID string) (int64, string, error) {
	if acheteurID <= 0 {
		return 0, "", fmt.Errorf("ID acheteur invalide")
	}

	var prenomAcheteur string
	err := Conn.QueryRow("SELECT prenom FROM UTILISATEUR WHERE id = ?", acheteurID).Scan(&prenomAcheteur)
	if err != nil {
		return 0, "", fmt.Errorf("l'acheteur n'existe pas dans la base de données")
	}

	var vendeurID int
	var titreProjet string
	err = Conn.QueryRow("SELECT id_createur, titre FROM PROJET_UPCYCLING WHERE id = ? AND statut = 'Disponible'", projetID).Scan(&vendeurID, &titreProjet)
	if err != nil {
		return 0, "", fmt.Errorf("ce projet n'est plus disponible à la vente")
	}

	tx, err := Conn.Begin()
	if err != nil {
		return 0, "", err
	}
	defer tx.Rollback() 

	res, err := tx.Exec(`
		UPDATE PROJET_UPCYCLING 
		SET id_acheteur = ?, statut = 'Vendu', date_achat = NOW() 
		WHERE id = ? AND statut = 'Disponible'`, acheteurID, projetID)
	if err != nil {
		return 0, "", err
	}
	
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return 0, "", fmt.Errorf("le projet a déjà été vendu (conflit de concurrence)")
	}

	commission := prixBase * 0.05 
	montantTotalTTC := prixBase + commission

	res, err = tx.Exec(`
		INSERT INTO COMMANDE (id_utilisateur, montant_total, statut, date_commande) 
		VALUES (?, ?, 'Payee', NOW())`, acheteurID, montantTotalTTC)
	if err != nil {
		return 0, "", err
	}
	commandeID, _ := res.LastInsertId()

	_, err = tx.Exec(`
		INSERT INTO LIGNE_COMMANDE (id_commande, id_vendeur, type_item, reference_id, prix_unitaire, commission_upc) 
		VALUES (?, ?, 'Projet', ?, ?, ?)`, commandeID, vendeurID, projetID, prixBase, commission)
	if err != nil {
		return 0, "", err
	}

	// 5. Transaction
	res, err = tx.Exec(`
		INSERT INTO `+"`TRANSACTION`"+` (id_commande, id_acheteur, montant_total, statut_paiement, date_transaction, stripe_payment_id) 
		VALUES (?, ?, ?, 'Valide', NOW(), ?)`, commandeID, acheteurID, montantTotalTTC, stripePaymentID)
	if err != nil {
		return 0, "", err
	}
	transactionID, _ := res.LastInsertId()

	factureID, numeroFacture, err := CreateFactureForTransaction(tx, transactionID)
	if err != nil {
		return 0, "", err
	}

	_, _ = tx.Exec("UPDATE DM_SALE SET status = 'Payee' WHERE id_projet = ? AND id_buyer = ?", projetID, acheteurID)
	_, _ = tx.Exec("UPDATE DM_SALE SET status = 'Annulee' WHERE id_projet = ? AND id_buyer != ?", projetID, acheteurID)
	_, _ = tx.Exec("UPDATE DM_OFFER SET status = 'Annulee' WHERE id_projet = ? AND status = 'En attente'", projetID)

	// 8. Notification
	titreNotif := fmt.Sprintf("Votre projet '%s' a été vendu !", titreProjet)
	messageNotif := fmt.Sprintf("%s a acheté votre création. Organisez la livraison ou la remise en main propre.", prenomAcheteur)
	_, _ = tx.Exec(`
		INSERT INTO NOTIFICATION (id_utilisateur, id_emetteur, type, titre, message) 
		VALUES (?, ?, 'Alerte', ?, ?)`, vendeurID, acheteurID, titreNotif, messageNotif)

	if err := tx.Commit(); err != nil {
		return 0, "", err
	}

	return int64(factureID), numeroFacture, nil
}