package db

import (
	"upcycleconnect/api-go/models"
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
			SELECT m.id, m.id_forum, m.id_utilisateur, m.contenu, m.date_envoi, u.prenom, u.nom, u.role 
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
			
			if err := msgRows.Scan(&m.ID, &topicID, &m.UserID, &m.Content, &m.PostedAt, &prenom, &nom, &m.Role); err == nil {
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