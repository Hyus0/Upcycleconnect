package db

import (
	"fmt"
	"upcycleconnect/api-go/models"
)

func GetAllNotifications() ([]models.Notification, error) {
	rows, err := Conn.Query(`
		SELECT id, id_utilisateur, type, titre, message, lu, date_envoi 
		FROM NOTIFICATION 
		ORDER BY date_envoi DESC`)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifs []models.Notification
	for rows.Next() {
		var n models.Notification
		var luInt int
		var dateEnvoi []byte

		err := rows.Scan(&n.ID, &n.IDUtilisateur, &n.Type, &n.Titre, &n.Message, &luInt, &dateEnvoi)
		if err != nil {
			return nil, err
		}
		
		n.Lu = (luInt == 1)
		n.DateEnvoi = string(dateEnvoi)
		
		notifs = append(notifs, n)
	}
	return notifs, nil
}

func GetNotificationByID(id int) (*models.Notification, error) {
	var n models.Notification
	var luInt int
	var dateEnvoi []byte

	err := Conn.QueryRow(`
		SELECT id, id_utilisateur, type, titre, message, lu, date_envoi 
		FROM NOTIFICATION 
		WHERE id = ?`, id).Scan(&n.ID, &n.IDUtilisateur, &n.Type, &n.Titre, &n.Message, &luInt, &dateEnvoi)

	if err != nil {
		return nil, err
	}

	n.Lu = (luInt == 1)
	n.DateEnvoi = string(dateEnvoi)

	return &n, nil
}

func CreerNotification(idUtilisateur int, idEmetteur int, typeNotif, titre, message string) error {
	_, err := Conn.Exec(`
		INSERT INTO NOTIFICATION (id_utilisateur, id_emetteur, type, titre, message) 
		VALUES (?, ?, ?, ?, ?)`, 
		idUtilisateur, idEmetteur, typeNotif, titre, message)
	return err
}

func GetNotificationsByUser(userID int) ([]models.Notification, error) {
	rows, err := Conn.Query(`
		SELECT id, id_utilisateur, id_emetteur, type, titre, message, lu, date_envoi 
		FROM NOTIFICATION 
		WHERE id_utilisateur = ? 
		ORDER BY date_envoi DESC`, userID)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifs []models.Notification
	for rows.Next() {
		var n models.Notification
		var luInt int
		var dateEnvoi []byte

		err := rows.Scan(&n.ID, &n.IDUtilisateur, &n.IDEmetteur, &n.Type, &n.Titre, &n.Message, &luInt, &dateEnvoi)
		
		if err != nil {
			return nil, err 
		}
		
		n.Lu = (luInt == 1)
		n.DateEnvoi = string(dateEnvoi)
		
		notifs = append(notifs, n)
	}
	return notifs, nil
}

func MarquerNotificationLue(notifID int, userID int) error {
	res, err := Conn.Exec("UPDATE NOTIFICATION SET lu = 1 WHERE id = ? AND id_utilisateur = ?", notifID, userID)
	if err != nil {
		return err
	}
	
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("notification introuvable ou vous n'avez pas les droits")
	}
	return nil
}