package db

import (
	"upcycleconnect/api-go/models"
	"upcycleconnect/api-go/passwordHashing"
	"database/sql"
	"fmt"
)

func GetAllUsers() ([]models.GetUser, error) {
	users := []models.GetUser{}
	query := `SELECT id, prenom, nom, mail, adresse, ville, code_postal,
		          COALESCE(date_naissance, ''),
		          COALESCE(date_inscription, ''),
		          role, id_langue
		          FROM UTILISATEUR`
	rows, err := Conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetAllUsers query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u models.GetUser
		err := rows.Scan(&u.Id, &u.Prenom, &u.Nom, &u.Mail, &u.Adresse, &u.Ville, &u.CodePostal, &u.DateNaissance, &u.DateInscription, &u.Role, &u.IdLangue)
		if err != nil {
			return nil, fmt.Errorf("GetAllUsers scan: %v", err)
		}
		users = append(users, u)
	}
	return users, nil
}

func GetUser(userId int) (*models.GetUser, error) {
	var u models.GetUser
	query := `SELECT id, prenom, nom, mail, adresse, ville, code_postal,
		          COALESCE(date_naissance, ''),
		          COALESCE(date_inscription, ''),
					COALESCE(date_update_password, ''),
		          role, id_langue
		          FROM UTILISATEUR WHERE id = ?`
	row := Conn.QueryRow(query, userId)

	err := row.Scan(&u.Id, &u.Prenom, &u.Nom, &u.Mail, &u.Adresse, &u.Ville, &u.CodePostal, &u.DateNaissance, &u.DateInscription, &u.DateUpdatePassword, &u.Role, &u.IdLangue)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GetUser scan: %v", err)
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
	if user.IdLangue <= 0 {
		user.IdLangue = 1
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
		user.IdLangue,
	)
	if err != nil {
		return err 
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	queryScore := `INSERT INTO UPCYCLING_SCORE (id_utilisateur, ressources_economisees, co2_total_evite_kg, nb_objets_recycles, total_points, niveau) 
	               VALUES (?, 0, 0, 0, 0, 'Débutant 🌱')`

	_, err = Conn.Exec(queryScore, lastInsertId)
	return err
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
	if err != nil { return err }

	count, _ := result.RowsAffected()
	if count == 0 {
		return fmt.Errorf("aucun utilisateur trouvé") }
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
                niveau,
                co2_total_evite_kg,
                nb_objets_recycles,
                ressources_economisees
	          FROM UPCYCLING_SCORE
              WHERE id_utilisateur = ?`

	err := Conn.QueryRow(query, userId).Scan(
		&s.Points,
		&s.Niveau,
		&s.Co2Evite,
		&s.ObjetsRecycles,
		&s.ArgentEconomise,
	)

	if err == sql.ErrNoRows {
		return &models.UserStats{Niveau: "Débutant"}, nil
	}
	if err != nil {
		return nil, err
	}
	if s.Points > 500 {
		s.Niveau = "Héro 🌿"
	} else if s.Points > 200 {
		s.Niveau = "Protecteur 🍃"
	} else {
		s.Niveau = "Débutant 🌱"
	}

	return &s, nil
}
