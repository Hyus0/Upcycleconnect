package db

import (
	"upcycleconnect/api-go/models"
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
		          role, id_langue 
		          FROM UTILISATEUR WHERE id = ?`
	row := Conn.QueryRow(query, userId)

	err := row.Scan(&u.Id, &u.Prenom, &u.Nom, &u.Mail, &u.Adresse, &u.Ville, &u.CodePostal, &u.DateNaissance, &u.DateInscription, &u.Role, &u.IdLangue)

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

	query := `INSERT INTO UTILISATEUR (prenom, nom, password, mail, adresse, ville, code_postal, date_naissance, role, id_langue)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	var dateN interface{} = user.DateNaissance
	if user.DateNaissance == "" {
		dateN = nil
	}

	_, err := Conn.Exec(query,
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

func DeleteUser(id int) error {
	result, err := Conn.Exec("DELETE FROM UTILISATEUR WHERE id = ?", id)
	if err != nil { return err }
	
	count, _ := result.RowsAffected()
	if count == 0 { return fmt.Errorf("aucun utilisateur trouvé") }
	return nil
}