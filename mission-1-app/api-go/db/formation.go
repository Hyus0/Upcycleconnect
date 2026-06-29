package db

import (
	"database/sql"
	"fmt"
	"time"
	"upcycleconnect/api-go/models"
)

func GetAllFormations() ([]models.GetFormation, error) {
	if Conn == nil {
		return nil, fmt.Errorf("connexion DB non initialisee")
	}

	query := `
		SELECT
			f.id, f.id_formateur, f.type, f.titre, f.description, f.capacite_max, f.est_valide,
			COUNT(DISTINCT fi.id_utilisateur) AS nb_inscrit,
			COALESCE(MIN(fs.date_debut), '1970-01-01 00:00:00') as date_debut, 
			COALESCE(MAX(fs.date_fin), '1970-01-01 00:00:00') as date_fin,
			f.statut, f.prix_unitaire, f.adresse, f.ville, f.code_postal
		FROM FORMATION f
		LEFT JOIN FORMATION_SESSION fs ON fs.id_formation = f.id
		LEFT JOIN FORMATION_INSCRIPTION fi ON fi.id_formation = f.id
		GROUP BY
			f.id, f.id_formateur, f.type, f.titre, f.description, f.capacite_max,
			f.est_valide, f.statut, f.prix_unitaire, f.adresse, f.ville, f.code_postal
	`

	rows, err := Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	formations := []models.GetFormation{}
	for rows.Next() {
		var f models.GetFormation
		err := rows.Scan(
			&f.ID, &f.ID_formateur, &f.Type, &f.Titre, &f.Description, &f.Capacite_max,
			&f.Est_valide, &f.Nb_inscrit, &f.Date_debut, &f.Date_fin,
			&f.Statut, &f.Prix_unitaire, &f.Adresse, &f.Ville, &f.CodePostal,
		)
		if err != nil {
			fmt.Println("Erreur Scan Formation:", err)
			return nil, err
		}
		
		if f.Date_debut == "1970-01-01 00:00:00" { f.Date_debut = "" }
		if f.Date_fin == "1970-01-01 00:00:00" { f.Date_fin = "" }
		
		formations = append(formations, f)
	}

	return formations, nil
}

func GetFormation(formationID int, userID int) (*models.GetFormation, error) {
	if Conn == nil {
		return nil, fmt.Errorf("connexion DB non initialisee")
	}

	query := `
		SELECT 
			f.id, f.id_formateur, u.prenom, u.nom, COALESCE(u.image_profil, ''), f.type, f.titre, f.description, 
			f.capacite_max, 
			COALESCE(MIN(fs.date_debut), '1970-01-01 00:00:00'), 
			COALESCE(MAX(fs.date_fin), '1970-01-01 00:00:00'), 
			f.statut, f.prix_unitaire, f.adresse, f.ville, f.code_postal,
			(SELECT COUNT(DISTINCT id_utilisateur) FROM FORMATION_INSCRIPTION WHERE id_formation = f.id) as nb_inscrit
		FROM FORMATION f
		JOIN UTILISATEUR u ON f.id_formateur = u.id
		LEFT JOIN FORMATION_SESSION fs ON fs.id_formation = f.id
		WHERE f.id = ?
		GROUP BY f.id`

	var f models.GetFormation
	
	err := Conn.QueryRow(query, formationID).Scan(
		&f.ID, &f.ID_formateur, &f.Prenom_formateur, &f.Nom_formateur, &f.Image_formateur,
		&f.Type, &f.Titre, &f.Description, &f.Capacite_max,
		&f.Date_debut, &f.Date_fin, &f.Statut, &f.Prix_unitaire, 
		&f.Adresse, &f.Ville, &f.CodePostal, &f.Nb_inscrit,       
	)
	if err != nil {
		return nil, err
	}

	if f.Date_debut == "1970-01-01 00:00:00" { f.Date_debut = "" }
	if f.Date_fin == "1970-01-01 00:00:00" { f.Date_fin = "" }

	sessionRows, err := Conn.Query("SELECT id, nom, date_debut, date_fin, statut FROM FORMATION_SESSION WHERE id_formation = ?", formationID)
	if err == nil {
		defer sessionRows.Close()
		for sessionRows.Next() {
			var s models.FormationSession
			if err := sessionRows.Scan(&s.ID, &s.Nom, &s.DateDebut, &s.DateFin, &s.Statut); err == nil {
				f.Sessions = append(f.Sessions, s)
			}
		}
	}

	var count int
	checkQuery := "SELECT COUNT(*) FROM FORMATION_INSCRIPTION WHERE id_utilisateur = ? AND id_formation = ?"
	Conn.QueryRow(checkQuery, userID, formationID).Scan(&count)

	f.IsRegistered = (count > 0)

	return &f, nil
}

func CreateFormation(formation models.Formation) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisee")
	}

	tx, err := Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO FORMATION (
		id_formateur, type, titre, description, capacite_max, est_valide,
		statut, prix_unitaire, adresse, ville, code_postal
	) VALUES (?, ?, ?, ?, ?, "En attente", ?, ?, ?, ?, ?)`

	res, err := tx.Exec(
		query,
		formation.ID_formateur, formation.Type, formation.Titre,
		formation.Description, formation.Capacite_max, formation.Statut,
		formation.Prix_unitaire, formation.Adresse, formation.Ville, formation.Code_postal,
	)
	if err != nil {
		return fmt.Errorf("CreateFormation: %v", err)
	}

	formationID, _ := res.LastInsertId()

	for _, s := range formation.Sessions {
		_, err = tx.Exec(
			"INSERT INTO FORMATION_SESSION (id_formation, nom, date_debut, date_fin, statut) VALUES (?, ?, ?, ?, 'Ouvert')",
			formationID, s.Nom, s.DateDebut, s.DateFin,
		)
		if err != nil {
			return fmt.Errorf("erreur ajout session %s : %v", s.Nom, err)
		}
	}

	return tx.Commit()
}

func ModifyFormation(id int, f models.Formation) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisee")
	}

	tx, err := Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queryFull := `
		UPDATE FORMATION SET
			id_formateur = ?, type = ?, titre = ?, description = ?,
			capacite_max = ?, statut = ?, prix_unitaire = ?,
			adresse = ?, ville = ?, code_postal = ?
		WHERE id = ?
	`
	res, err := tx.Exec(queryFull, f.ID_formateur, f.Type, f.Titre, f.Description, f.Capacite_max, f.Statut, f.Prix_unitaire, f.Adresse, f.Ville, f.Code_postal, id)
	if err != nil {
		return fmt.Errorf("ModifyFormation: %v", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		var exists int
		tx.QueryRow("SELECT 1 FROM FORMATION WHERE id = ?", id).Scan(&exists)
		if exists == 0 {
			return fmt.Errorf("aucune formation trouvee avec l'ID %d", id)
		}
	}

	for _, s := range f.Sessions {
		if s.ID == 0 {
			_, err = tx.Exec(
				"INSERT INTO FORMATION_SESSION (id_formation, nom, date_debut, date_fin, statut) VALUES (?, ?, ?, ?, 'Ouvert')",
				id, s.Nom, s.DateDebut, s.DateFin,
			)
			if err != nil { return fmt.Errorf("insert session: %v", err) }
			
		} else {
			var debutStr string
			err = tx.QueryRow("SELECT date_debut FROM FORMATION_SESSION WHERE id = ? AND id_formation = ?", s.ID, id).Scan(&debutStr)
			if err != nil {
				continue 
			}

			debut, parseErr := time.Parse("2006-01-02 15:04:05", debutStr)
			if parseErr != nil {
				debut, parseErr = time.Parse(time.RFC3339, debutStr)
			}

			if parseErr == nil && debut.After(time.Now()) {
				_, err = tx.Exec("UPDATE FORMATION_SESSION SET nom = ?, date_debut = ?, date_fin = ? WHERE id = ?", s.Nom, s.DateDebut, s.DateFin, s.ID)
			} else {
				_, err = tx.Exec("UPDATE FORMATION_SESSION SET nom = ? WHERE id = ?", s.Nom, s.ID)
			}
			
			if err != nil {
				return fmt.Errorf("update session %d: %v", s.ID, err)
			}
		}
	}

	return tx.Commit()
}

func DeleteFormation(id int) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisee")
	}

	result, err := Conn.Exec("DELETE FROM FORMATION WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("DeleteFormation: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteFormation RowsAffected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("aucune formation trouvee avec l'ID %d", id)
	}
	return nil
}

func JoinFormation(userID int, formationID int, sessionID int) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisee")
	}

	var capaciteMax int
	var nbInscrits int
	queryCheck := `
		SELECT f.capacite_max, COUNT(fi.id_utilisateur)
		FROM FORMATION f
		LEFT JOIN FORMATION_INSCRIPTION fi ON fi.id_session = ?
		WHERE f.id = ?
		GROUP BY f.id, f.capacite_max
	`
	err := Conn.QueryRow(queryCheck, sessionID, formationID).Scan(&capaciteMax, &nbInscrits)
	if err == sql.ErrNoRows {
		return fmt.Errorf("formation ou session introuvable")
	}
	if err != nil {
		return fmt.Errorf("verification capacite formation: %v", err)
	}

	if nbInscrits >= capaciteMax {
		return fmt.Errorf("cette session est complète")
	}

	queryInsert := "INSERT INTO FORMATION_INSCRIPTION (id_utilisateur, id_formation, id_session) VALUES (?, ?, ?)"
	if _, err := Conn.Exec(queryInsert, userID, formationID, sessionID); err != nil {
		return fmt.Errorf("insertion FORMATION_INSCRIPTION: %v", err)
	}

	return nil
}

func QuitFormation(userID int, sessionID int) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisee")
	}

	query := "DELETE FROM FORMATION_INSCRIPTION WHERE id_utilisateur = ? AND id_session = ?"
	result, err := Conn.Exec(query, userID, sessionID)
	if err != nil {
		return fmt.Errorf("erreur lors de la désinscription: %v", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("aucune inscription trouvée pour cet utilisateur à cette session")
	}

	return nil
}

func GetUserFormations(userID int) ([]models.GetFormation, error) {
	query := `
		SELECT f.id, f.id_formateur, f.type, f.titre, f.description, f.capacite_max,
		       fs.date_debut, fs.date_fin, f.statut, f.prix_unitaire, f.adresse, f.ville, f.code_postal
		FROM FORMATION f
		INNER JOIN FORMATION_INSCRIPTION fi ON f.id = fi.id_formation
		INNER JOIN FORMATION_SESSION fs ON fi.id_session = fs.id
		WHERE fi.id_utilisateur = ?`

	rows, err := Conn.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var formations []models.GetFormation
	for rows.Next() {
		var f models.GetFormation
		rows.Scan(&f.ID, &f.ID_formateur, &f.Type, &f.Titre, &f.Description, &f.Capacite_max,
			&f.Date_debut, &f.Date_fin, &f.Statut, &f.Prix_unitaire, &f.Adresse, &f.Ville, &f.CodePostal)
		f.IsRegistered = true
		formations = append(formations, f)
	}
	return formations, nil
}

func GetFormationParticipants(formationID int) ([]models.Participant, error) {
	query := `
		SELECT 
			u.id, u.prenom, u.nom, u.mail, COALESCE(u.image_profil, '') as image_profil, u.role
		FROM UTILISATEUR u
		JOIN FORMATION_INSCRIPTION i ON u.id = i.id_utilisateur
		WHERE i.id_formation = ?
	`
	
	rows, err := Conn.Query(query, formationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []models.Participant
	for rows.Next() {
		var p models.Participant
		if err := rows.Scan(&p.ID, &p.Prenom, &p.Nom, &p.Mail, &p.ImageProfil, &p.Role); err == nil {
			participants = append(participants, p)
		}
	}
	return participants, nil
}

func AnnulerFormation(formationID int, formateurID int) error {
	if Conn == nil {
		return fmt.Errorf("connexion DB non initialisee")
	}

	tx, err := Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var titre string
	err = tx.QueryRow("SELECT titre FROM FORMATION WHERE id = ? AND id_formateur = ?", formationID, formateurID).Scan(&titre)
	if err == sql.ErrNoRows {
		return fmt.Errorf("formation introuvable ou vous n'êtes pas autorisé à l'annuler")
	} else if err != nil {
		return err
	}

	rows, err := tx.Query("SELECT DISTINCT id_utilisateur FROM FORMATION_INSCRIPTION WHERE id_formation = ?", formationID)
	if err != nil {
		return err
	}
	
	var inscrits []int
	for rows.Next() {
		var uID int
		if err := rows.Scan(&uID); err == nil {
			inscrits = append(inscrits, uID)
		}
	}
	rows.Close()

	_, err = tx.Exec("UPDATE FORMATION SET statut = 'Annule' WHERE id = ?", formationID)
	if err != nil { return err }

	_, err = tx.Exec("UPDATE FORMATION_SESSION SET statut = 'Annule' WHERE id_formation = ?", formationID)
	if err != nil { return err }

	_, err = tx.Exec("DELETE FROM FORMATION_INSCRIPTION WHERE id_formation = ?", formationID)
	if err != nil { return err }

	if len(inscrits) > 0 {
		titreNotif := "Formation annulée"
		msgNotif := fmt.Sprintf("La formation \"%s\" a été annulée par le formateur, vous serez remboursé d'ici 3-4 jours ouvrables.", titre)

		for _, uID := range inscrits {
			_, err = tx.Exec(`
				INSERT INTO NOTIFICATION (id_utilisateur, id_emetteur, type, titre, message) 
				VALUES (?, ?, 'Alerte', ?, ?)`, 
				uID, formateurID, titreNotif, msgNotif)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}