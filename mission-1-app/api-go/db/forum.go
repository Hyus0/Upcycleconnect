package db

import (
	"upcycleconnect/api-go/models"
	"sort"
)

func GetAllForums() ([]models.ForumCategory, error) {
	salonRows, err := Conn.Query("SELECT id, nom, description FROM FORUM_SALON")
	if err != nil { 
		return nil, err 
	}
	defer salonRows.Close()

	var salons []models.ForumCategory
	salonsMap := make(map[int]*models.ForumCategory)

	for salonRows.Next() {
		var s models.ForumCategory
		if err := salonRows.Scan(&s.ID, &s.Name, &s.Description); err == nil {
			s.Topics = []models.ForumTopic{}
			salons = append(salons, s)
		}
	}
	for i := range salons {
		salonsMap[salons[i].ID] = &salons[i]
	}

	topicRows, err := Conn.Query("SELECT id, id_salon, titre, sujet, date_creation FROM FORUM WHERE id_salon IS NOT NULL ORDER BY date_creation DESC")
	if err != nil { 
		return nil, err 
	}
	defer topicRows.Close()

	topicsMap := make(map[int]*models.ForumTopic)
	for topicRows.Next() {
		var t models.ForumTopic
		var idSalon int
		if err := topicRows.Scan(&t.ID, &idSalon, &t.Title, &t.Preview, &t.LastActivity); err == nil {
			t.Messages = []models.ForumMessage{}
			if salon, ok := salonsMap[idSalon]; ok {
				salon.Topics = append(salon.Topics, t)
				topicsMap[t.ID] = &salon.Topics[len(salon.Topics)-1]
			}
		}
	}

	msgQuery := `
			SELECT m.id, m.id_forum, m.id_utilisateur, m.contenu, m.date_envoi, u.prenom, u.nom, u.role, COALESCE(u.image_profil, '') as image_profil
			FROM FORUM_MESSAGE m
			JOIN UTILISATEUR u ON m.id_utilisateur = u.id
			ORDER BY m.date_envoi ASC`
			
		msgRows, err := Conn.Query(msgQuery)
		if err != nil { 
			return nil, err 
		}
		defer msgRows.Close()

		for msgRows.Next() {
			var m models.ForumMessage
			var topicID int
			var prenom, nom string
			
			if err := msgRows.Scan(&m.ID, &topicID, &m.UserID, &m.Content, &m.PostedAt, &prenom, &nom, &m.Role, &m.ImageProfil,); err == nil {
				m.Author = prenom + " " + nom
				if topic, ok := topicsMap[topicID]; ok {
					topic.Messages = append(topic.Messages, m)
					topic.LastActivity = m.PostedAt
				}
			} else {
				println("Erreur de scan message:", err.Error())
			}
		}

		return salons, nil
}

func CreateForumTopic(userID int, salonID int, titre string, sujet string) error {
	res, err := Conn.Exec("INSERT INTO FORUM (id_utilisateur, id_salon, titre, sujet) VALUES (?, ?, ?, ?)", userID, salonID, titre, sujet)
	if err != nil { 
		return err 
	}

	forumID, err := res.LastInsertId()
	if err != nil { 
		return err 
	}

	_, err = Conn.Exec("INSERT INTO FORUM_MESSAGE (id_utilisateur, id_forum, contenu) VALUES (?, ?, ?)", userID, forumID, sujet)
	return err
}

func SendMessageForum(userID int, forumID int, contenu string) error {
	_, err := Conn.Exec("INSERT INTO FORUM_MESSAGE (id_utilisateur, id_forum, contenu) VALUES (?, ?, ?)", userID, forumID, contenu)
	return err
}

func SignalerMessageForum(messageID int, userID int, motif string) error {
	_, err := Conn.Exec(`
		INSERT IGNORE INTO MESSAGE_SIGNALEMENT (id_message, id_utilisateur, motif) 
		VALUES (?, ?, ?)`, messageID, userID, motif)
	return err
}

func GetTopMessagesSignales() ([]models.ReportedMessage, error) {
	query := `
		SELECT 
			m.id, m.contenu, m.date_envoi, 
			u.id as author_id, u.prenom as author_prenom, u.nom as author_nom, COALESCE(u.image_profil, '') as author_image,
			r.prenom as reporter_prenom, r.nom as reporter_nom, s.motif
		FROM FORUM_MESSAGE m
		JOIN MESSAGE_SIGNALEMENT s ON m.id = s.id_message
		JOIN UTILISATEUR u ON m.id_utilisateur = u.id
		JOIN UTILISATEUR r ON s.id_utilisateur = r.id
	`
	rows, err := Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	msgMap := make(map[int]*models.ReportedMessage)

	for rows.Next() {
		var msgID, authorID int
		var contenu, dateEnvoi, authorP, authorN, authorImage, repP, repN, motif string
		
		if err := rows.Scan(&msgID, &contenu, &dateEnvoi, &authorID, &authorP, &authorN, &authorImage, &repP, &repN, &motif); err == nil {
			if _, exists := msgMap[msgID]; !exists {
				msgMap[msgID] = &models.ReportedMessage{
					MessageID: msgID, Contenu: contenu, DateEnvoi: dateEnvoi,
					AuthorID: authorID, AuthorName: authorP + " " + authorN,
					AuthorImage: authorImage, 
					NbSignalements: 0,
					Details: []models.ReportDetail{},
				}
			}
			msgMap[msgID].NbSignalements++
			msgMap[msgID].Details = append(msgMap[msgID].Details, models.ReportDetail{
				ReporterName: repP + " " + repN,
				Motif: motif,
			})
		}
	}

	var reported []models.ReportedMessage
	for _, msg := range msgMap {
		reported = append(reported, *msg)
	}

	sort.Slice(reported, func(i, j int) bool {
		return reported[i].NbSignalements > reported[j].NbSignalements
	})

	return reported, nil
}

func GetModerationTopics() ([]models.ModTopic, error) {
	query := `
		SELECT 
			f.id, f.titre, f.sujet, u.id, u.prenom, u.nom, f.date_creation,
			(SELECT COUNT(*) FROM FORUM_MESSAGE WHERE id_forum = f.id) as msg_count
		FROM FORUM f
		JOIN UTILISATEUR u ON f.id_utilisateur = u.id
		ORDER BY f.date_creation DESC
	`
	rows, err := Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []models.ModTopic
	for rows.Next() {
		var t models.ModTopic
		var prenom, nom string
		var authorID int
		
		if err := rows.Scan(&t.ID, &t.Titre, &t.Sujet, &authorID, &prenom, &nom, &t.Date, &t.MsgCount); err == nil {
			t.AuthorID = authorID
			t.AuthorName = prenom + " " + nom
			topics = append(topics, t)
		}
	}
	return topics, nil
}

func GetBannedUsers() ([]models.BannedUser, error) {
	rows, err := Conn.Query("SELECT id, prenom, nom, role, COALESCE(image_profil, '') FROM UTILISATEUR WHERE ban_forum = true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banned []models.BannedUser
	for rows.Next() {
		var u models.BannedUser
		var prenom, nom, imageProfil string
		
		if err := rows.Scan(&u.ID, &prenom, &nom, &u.Role, &imageProfil); err == nil {
			u.FullName = prenom + " " + nom
			u.ImageProfil = imageProfil 
			banned = append(banned, u)
		}
	}
	return banned, nil
}

func ToggleBanForum(userID int, isBanned bool) error {
	_, err := Conn.Exec("UPDATE UTILISATEUR SET ban_forum = ? WHERE id = ?", isBanned, userID)
	if err != nil {
		return err
	}

	titre := "Modération du forum"
	message := "Votre accès au forum a été restreint suite au non-respect de notre charte."
	type_notification := "Alerte"
	if !isBanned {
		message = "Bonne nouvelle ! Votre accès au forum a été rétabli. Merci de respecter la charte."
	}

	_, err = Conn.Exec("INSERT INTO NOTIFICATION (id_utilisateur, type,  titre, message) VALUES (?, ?, ?)", userID, type_notification, titre, message)
	return err
}

func truncateText(s string, max int) string {
	if len(s) > max {
		return s[:max] + "..."
	}
	return s
}

func DeleteMessageForum(messageID int) error {
	var authorID int
	var contenu string
	err := Conn.QueryRow("SELECT id_utilisateur, contenu FROM FORUM_MESSAGE WHERE id = ?", messageID).Scan(&authorID, &contenu)
	if err != nil {
		return err 
	}

	_, err = Conn.Exec("DELETE FROM FORUM_MESSAGE WHERE id = ?", messageID)
	if err != nil {
		return err
	}

	titre := "Modération du forum"
	type_notification := "Alerte"
	messageNotif := "Attention, votre message (\"" + truncateText(contenu, 40) + "\") a été supprimé par l'équipe de modération car il ne respecte pas notre charte."

	_, err = Conn.Exec("INSERT INTO NOTIFICATION (id_utilisateur, type, titre, message) VALUES (?, ?, ?, ?)", authorID, type_notification, titre, messageNotif)
	
	return err
}

func DeleteTopicForum(topicID int) error {
	_, _ = Conn.Exec("DELETE FROM FORUM_MESSAGE WHERE id_forum = ?", topicID)
	
	_, err := Conn.Exec("DELETE FROM FORUM WHERE id = ?", topicID)
	return err
}

func GetRecentForumMessages() ([]models.ForumMessage, error) {
	query := `
		SELECT 
			m.id, 
			m.id_forum, 
			COALESCE(f.titre, 'Discussion supprimée'), 
			m.id_utilisateur, 
			m.contenu, 
			m.date_envoi, 
			u.prenom, 
			u.nom, 
			u.role, 
			COALESCE(u.image_profil, '')
		FROM FORUM_MESSAGE m
		JOIN UTILISATEUR u ON m.id_utilisateur = u.id
		LEFT JOIN FORUM f ON m.id_forum = f.id
		ORDER BY m.date_envoi DESC
		LIMIT 50`

	rows, err := Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.ForumMessage
	for rows.Next() {
		var m models.ForumMessage
		var prenom, nom string
		var imageProfil string

		if err := rows.Scan(
			&m.ID, 
			&m.ForumID, 
			&m.TopicTitle, 
			&m.UserID, 
			&m.Content, 
			&m.PostedAt, 
			&prenom, 
			&nom, 
			&m.Role, 
			&imageProfil,
		); err == nil {
			m.Author = prenom + " " + nom
			m.ImageProfil = imageProfil
			messages = append(messages, m)
		}
	}
	
	if messages == nil {
		messages = []models.ForumMessage{}
	}
	
	return messages, nil
}