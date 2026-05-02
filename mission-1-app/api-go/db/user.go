package db

import (
	"database/sql"
	"fmt"
	"upcycleconnect/api-go/models"
	"upcycleconnect/api-go/passwordHashing"
)

func GetAllUsers() ([]models.GetUser, error) {
	users := []models.GetUser{}
	query := `SELECT id, prenom, nom, mail, adresse, ville, code_postal,
		          COALESCE(date_naissance, ''),
		          COALESCE(date_inscription, ''),
		          statut, role, id_langue
		          FROM UTILISATEUR`
	rows, err := Conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetAllUsers query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u models.GetUser
		var idLangue sql.NullInt64
		err := rows.Scan(&u.Id, &u.Prenom, &u.Nom, &u.Mail, &u.Adresse, &u.Ville, &u.CodePostal, &u.DateNaissance, &u.DateInscription, &u.Statut, &u.Role, &idLangue)
		if err != nil {
			return nil, fmt.Errorf("GetAllUsers scan: %v", err)
		}
		if idLangue.Valid {
			u.IdLangue = int(idLangue.Int64)
		}
		users = append(users, u)
	}
	return users, nil
}

func GetUser(userId int) (*models.GetUser, error) {
	var u models.GetUser
	var idLangue sql.NullInt64
	query := `SELECT id, prenom, nom, mail, adresse, ville, code_postal,
		          COALESCE(date_naissance, ''),
		          COALESCE(date_inscription, ''),
					COALESCE(date_update_password, ''),
		          statut, role, id_langue
		          FROM UTILISATEUR WHERE id = ?`
	row := Conn.QueryRow(query, userId)

	err := row.Scan(&u.Id, &u.Prenom, &u.Nom, &u.Mail, &u.Adresse, &u.Ville, &u.CodePostal, &u.DateNaissance, &u.DateInscription, &u.DateUpdatePassword, &u.Statut, &u.Role, &idLangue)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GetUser scan: %v", err)
	}
	if idLangue.Valid {
		u.IdLangue = int(idLangue.Int64)
	}

	return &u, nil
}

func EmailExists(email string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM UTILISATEUR WHERE mail = ?"

	err := Conn.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func CreateUser(user models.User) error {
	if user.Role == "" {
		user.Role = "Particulier"
	}

	hashed, err := passwordHashing.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed

	var dateN interface{} = nil
	if user.DateNaissance != "" {
		dateN = user.DateNaissance
	}
	var idLangue interface{} = nil
	if user.IdLangue > 0 {
		exists, err := languageExists(user.IdLangue)
		if err != nil {
			return err
		}
		if exists {
			idLangue = user.IdLangue
		}
	}

	query := `INSERT INTO UTILISATEUR (prenom, nom, password, mail, adresse, ville, code_postal, date_naissance, role, id_langue)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := Conn.Exec(query,
		user.Prenom,
		user.Nom,
		user.Password,
		user.Mail,
		user.Adresse,
		user.Ville,
		user.CodePostal,
		dateN,
		user.Role,
		idLangue,
	)
	if err != nil {
		return err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return err
	}

	queryScore := `INSERT INTO UPCYCLING_SCORE (id_utilisateur, ressources_economisees, co2_total_evite_kg, nb_objets_recycles, total_points)
	               VALUES (?, 0, 0, 0, 0)`

	_, err = Conn.Exec(queryScore, lastInsertId)
	return err
}

func languageExists(languageID int) (bool, error) {
	var count int
	err := Conn.QueryRow("SELECT COUNT(*) FROM LANGUE WHERE id = ?", languageID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func ModifyUser(userId int, user models.User) error {
	query := `UPDATE UTILISATEUR SET prenom=?, nom=?, mail=?, adresse=?, ville=?, code_postal=?, role=?, id_langue=? WHERE id=?`

	_, err := Conn.Exec(query,
		user.Prenom, user.Nom, user.Mail, user.Adresse,
		user.Ville, user.CodePostal, user.Role, user.IdLangue, userId,
	)
	return err
}

func ModifyUserPassword(userId int, user models.User) error {
	query := `UPDATE UTILISATEUR SET password=?, date_update_password=NOW() WHERE id=?`

	_, err := Conn.Exec(query,
		user.Password,
		userId,
	)
	return err
}

func GetPasswordHashed(userId int) (string, error) {
	var hash string
	query := `SELECT password FROM UTILISATEUR WHERE id = ?`
	err := Conn.QueryRow(query, userId).Scan(&hash)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func DeleteUser(id int) error {
	result, err := Conn.Exec("DELETE FROM UTILISATEUR WHERE id = ?", id)
	if err != nil {
		return err
	}

	count, _ := result.RowsAffected()
	if count == 0 {
		return fmt.Errorf("aucun utilisateur trouve")
	}
	return nil
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, prenom, nom, mail, password, role FROM UTILISATEUR WHERE mail = ?`

	err := Conn.QueryRow(query, email).Scan(
		&user.Id,
		&user.Prenom,
		&user.Nom,
		&user.Mail,
		&user.Password,
		&user.Role,
	)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUserToken(userId int, token string) error {
	query := "UPDATE UTILISATEUR SET token = ? WHERE id = ?"
	_, err := Conn.Exec(query, token, userId)
	return err
}

func VerifyUserByToken(userId int, token string) bool {
	var dbToken string
	query := "SELECT token FROM UTILISATEUR WHERE id = ?"

	err := Conn.QueryRow(query, userId).Scan(&dbToken)

	if err != nil || dbToken != token || token == "" {
		return false
	}

	return true
}

func GetUserStats(userId int) (*models.UserStats, error) {
	var s models.UserStats

	query := `SELECT
                total_points,
                co2_total_evite_kg,
                nb_objets_recycles,
                ressources_economisees
	          FROM UPCYCLING_SCORE
              WHERE id_utilisateur = ?`

	err := Conn.QueryRow(query, userId).Scan(
		&s.Points,
		&s.Co2Evite,
		&s.ObjetsRecycles,
		&s.ArgentEconomise,
	)

	if err == sql.ErrNoRows {
		return &models.UserStats{Niveau: "Debutant"}, nil
	}
	if err != nil {
		return nil, err
	}
	if s.Points > 500 {
		s.Niveau = "Hero"
	} else if s.Points > 200 {
		s.Niveau = "Protecteur"
	} else {
		s.Niveau = "Debutant"
	}

	return &s, nil
}

func GetUserPlanningEntries(userId int) ([]models.UserPlanningEntry, error) {
	entries := []models.UserPlanningEntry{}

	formationRows, err := Conn.Query(`
		SELECT
			f.id,
			'formation' AS kind,
			COALESCE(f.titre, ''),
			f.date_debut,
			COALESCE(f.adresse, ''),
			COALESCE(f.ville, ''),
			COALESCE(f.code_postal, '')
		FROM FORMATION_INSCRIPTION fi
		INNER JOIN FORMATION f ON f.id = fi.id_formation
		WHERE fi.id_utilisateur = ?
	`, userId)
	if err != nil {
		return nil, err
	}
	defer formationRows.Close()

	for formationRows.Next() {
		var entry models.UserPlanningEntry
		var dateTime sql.NullTime
		if err := formationRows.Scan(&entry.ID, &entry.Kind, &entry.Title, &dateTime, &entry.Address, &entry.City, &entry.PostalCode); err != nil {
			return nil, err
		}
		if dateTime.Valid {
			entry.DateTime = dateTime.Time.Format("2006-01-02 15:04:05")
		}
		entries = append(entries, entry)
	}

	eventRows, err := Conn.Query(`
		SELECT
			e.id,
			'event' AS kind,
			COALESCE(e.titre, ''),
			e.date_evenement,
			COALESCE(e.adresse, ''),
			COALESCE(e.ville, ''),
			COALESCE(e.code_postal, '')
		FROM EVENEMENT_INSCRIPTION ei
		INNER JOIN EVENEMENT e ON e.id = ei.id_evenement
		WHERE ei.id_utilisateur = ?
	`, userId)
	if err != nil {
		return nil, err
	}
	defer eventRows.Close()

	for eventRows.Next() {
		var entry models.UserPlanningEntry
		var dateTime sql.NullTime
		if err := eventRows.Scan(&entry.ID, &entry.Kind, &entry.Title, &dateTime, &entry.Address, &entry.City, &entry.PostalCode); err != nil {
			return nil, err
		}
		if dateTime.Valid {
			entry.DateTime = dateTime.Time.Format("2006-01-02 15:04:05")
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func GetPlatformOverview() (*models.PlatformOverview, error) {
	overview := &models.PlatformOverview{}

	if err := Conn.QueryRow(`SELECT COUNT(*) FROM UTILISATEUR`).Scan(&overview.MembersCount); err != nil {
		return nil, err
	}
	if err := Conn.QueryRow(`SELECT COUNT(*) FROM UTILISATEUR WHERE role = 'Prestataire'`).Scan(&overview.PrestatairesCount); err != nil {
		return nil, err
	}
	if err := Conn.QueryRow(`SELECT COALESCE(SUM(nb_objets_recycles), 0), COALESCE(SUM(co2_total_evite_kg), 0) FROM UPCYCLING_SCORE`).Scan(&overview.UpcycledObjectsCount, &overview.Co2SavedKg); err != nil {
		return nil, err
	}
	if err := Conn.QueryRow(`SELECT COUNT(*) FROM SITE WHERE actif = 1`).Scan(&overview.SitesCount); err != nil {
		return nil, err
	}
	if err := Conn.QueryRow(`SELECT COUNT(*) FROM ANNONCE`).Scan(&overview.AnnoncesCount); err != nil {
		return nil, err
	}
	if err := Conn.QueryRow(`SELECT COUNT(*) FROM EVENEMENT`).Scan(&overview.EventsCount); err != nil {
		return nil, err
	}
	if err := Conn.QueryRow(`SELECT COUNT(*) FROM FORMATION`).Scan(&overview.FormationsCount); err != nil {
		return nil, err
	}

	return overview, nil
}
