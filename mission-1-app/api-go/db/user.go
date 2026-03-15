package db

import (
	"upcycleconnect/api-go/models"
	"database/sql"
	"fmt"
)

func GetAllUsers() ([]models.GetUser, error) {
	users := []models.GetUser{}
	query := "SELECT id, prenom, nom, mail, adresse, ville, code_postal, date_naissance, date_inscription, role, id_langue FROM UTILISATEUR"
	rows, err := Conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetAllUsers query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u models.GetUser
		err := rows.Scan(&u.Id, &u.Prenom, &u.Nom, &u.Mail, &u.Adresse, &u.Ville, &u.Code_postal, &u.Date_naissance, &u.Date_inscription, &u.Role, &u.Id_langue)
		if err != nil {
			return nil, fmt.Errorf("GetAllUsers scan: %v", err)
		}
		users = append(users, u)
	}
	return users, nil
}

func GetUser(userId int) (*models.GetUser, error) {
	var u models.GetUser
	query := "SELECT id, prenom, nom, mail, adresse, ville, code_postal, date_naissance, date_inscription, role, id_langue FROM UTILISATEUR WHERE id = ?"
	row := Conn.QueryRow(query, userId)

	err := row.Scan(&u.Id, &u.Prenom, &u.Nom, &u.Mail, &u.Adresse, &u.Ville, &u.Code_postal, &u.Date_naissance, &u.Date_inscription, &u.Role, &u.Id_langue)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GetUser scan: %v", err)
	}

	return &u, nil
}

func CreateUser(user models.User) error {
	query := `INSERT INTO UTILISATEUR (prenom, nom, password, mail, adresse, ville, code_postal, date_naissance, role, id_langue, date_inscription)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`

	_, err := Conn.Exec(query,
		user.Prenom, user.Nom, user.Password, user.Mail,
		user.Adresse, user.Ville, user.Code_postal,
		user.Date_naissance, user.Role, user.Id_langue,
	)

	if err != nil {
		return fmt.Errorf("CreateUser: %v", err)
	}
	return nil
}

func ModifyUser(userId int, user models.User) error {
	query := `UPDATE UTILISATEUR SET
				prenom = ?,
				nom = ?,
				mail = ?,
				adresse = ?,
				ville = ?,
				code_postal = ?,
				role = ?,
				id_langue = ?
			  WHERE id = ?`
	result, err := Conn.Exec(query,
		user.Prenom,
		user.Nom,
		user.Mail,
		user.Adresse,
		user.Ville,
		user.Code_postal,
		user.Role,
		user.Id_langue,
		userId,
	)

	if err != nil {
		return fmt.Errorf("package db ModifyUser : %v", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("package db ModifyUser (RowsAffected) : %v", err.Error())
	}

	if rowsAffected == 0 {
		fmt.Printf("Note: Utilisateur %d non modifié (données identiques ou ID inexistant)\n", userId)
	}

	return nil
}

func DeleteUser(id int) error {
	if id <= 0 {
		return fmt.Errorf("package db DeleteUser : L'ID doit être un entier positif")
	}
	result, err := Conn.Exec("DELETE FROM UTILISATEUR WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("package db DeleteUser : échec de la suppression : %v", err.Error())
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("package db DeleteUser : erreur de RowsAffected : %v", err.Error())
	}
	if rowsAffected == 0 {
		return fmt.Errorf("package db DeleteUser : aucun utilisateur trouvé avec l'ID %d", id)
	}
	return nil
}
