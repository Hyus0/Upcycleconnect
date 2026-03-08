package db

import (
	"upcycleconnect/api-go/models"
	"database/sql"
	"fmt"
)

func GetAllUsers() ([]models.User, error) {
	users := []models.User{}
	rows, err := Conn.Query("SELECT id, prenom, nom FROM UTILISATEUR")
	if err != nil {
		return nil, fmt.Errorf("package DB GetAllUsers : %v", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.Prenom, &user.Nom)
		if err != nil {
			return nil, fmt.Errorf("package DB GetAllUsers : %v", err.Error())
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("package DB GetAllUsers : %v", err.Error())
	}
	return users, nil
}

func GetUser(userId int) (*models.User, error) { 
	user := models.User{Id: userId}
	row := Conn.QueryRow("SELECT id, prenom, nom FROM UTILISATEUR WHERE id = ?", userId)
	err := row.Scan(&user.Id, &user.Prenom, &user.Nom)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	if err != nil {
		return nil, fmt.Errorf("package DB GetUser: %v", err.Error())
	}

	return &user, nil 
}

func CreateUser(user models.User) error {
    if user.Id > 0 {
        _, err := Conn.Exec("INSERT INTO UTILISATEUR(id, prenom, nom) VALUES(?, ?, ?)", user.Id, user.Prenom, user.Nom)
        if err != nil {
            return fmt.Errorf("package db CreateUser (avec ID): %v", err.Error())
        }
    } else {
        _, err := Conn.Exec("INSERT INTO UTILISATEUR(prenom, nom) VALUES(?, ?)", user.Prenom, user.Nom)
        if err != nil {
            return fmt.Errorf("package db CreateUser (sans ID): %v", err.Error())
        }
    }
	return nil
}

func ModifyUser(userId int, user models.User) error {
	result, err := Conn.Exec("UPDATE UTILISATEUR SET prenom = ?, nom = ? WHERE id = ?", user.Prenom, user.Nom, userId)
	
	if err != nil {
		return fmt.Errorf("package db ModifyUser : échec de la mise à jour : %v", err.Error())
	}

 	rowsAffected, err := result.RowsAffected()

 	if rowsAffected == 0 {
 		return fmt.Errorf("package db ModifyUser : aucun user trouvé avec l'ID %d", userId)
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