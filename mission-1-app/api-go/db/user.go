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

