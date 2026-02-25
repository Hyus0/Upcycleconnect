package db

import (
	"upcylcle/models"
	"database/sql"
	"fmt"
)

func GetAllUsers() ([]models.User, error) {
	users := []models.User{}
	rows, err := Conn.Query("SELECT id, name, price FROM users")
	if err != nil {
		return nil, fmt.Errorf("package DB GetAllUsers : %v", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var user models.Game
		err := rows.Scan(&user.Id, &user.Name, &user.Price)
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
	row := Conn.QueryRow("SELECT name, price FROM users WHERE id = ?", userId)
	err := row.Scan(&user.Name, &user.Price)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	if err != nil {
		return nil, fmt.Errorf("package DB GetUser: %v", err.Error())
	}

	return &user, nil 
}

