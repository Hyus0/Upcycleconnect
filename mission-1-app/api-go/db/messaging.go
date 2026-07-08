package db

import (
	"database/sql"
	"fmt"
	"strings"
	"upcycleconnect/api-go/models"
)

const freeDMVendorLimit = 5

func HasActiveDMSubscription(userID int) (bool, error) {
	var count int
	err := Conn.QueryRow(`
		SELECT COUNT(*)
		FROM ABONNEMENT a
		JOIN TYPE_ABONNEMENT t ON t.id = a.id_type_abonnement
		WHERE a.id_acheteur = ?
		  AND a.statut = 'Actif'
		  AND (a.date_fin IS NULL OR a.date_fin >= NOW())
		  AND (t.nom = 'DM Plus' OR t.prix_ht <= 3.00)
	`, userID).Scan(&count)
	return count > 0, err
}

func StartConversation(currentUserID int, targetUserID int, annonceID *int, projetID *int) (models.DMStartResult, error) {
	result := models.DMStartResult{Allowed: false, Limit: freeDMVendorLimit}
	if currentUserID <= 0 {
		result.Message = "Vous devez etre connecte pour contacter un membre."
		return result, nil
	}
	if targetUserID <= 0 || currentUserID == targetUserID {
		result.Message = "Conversation impossible avec cet utilisateur."
		return result, nil
	}

	var existingID int
	var existingErr error

	switch {
	case annonceID != nil:
		existingErr = Conn.QueryRow(`
			SELECT id FROM DM_CONVERSATION
			WHERE ((id_user_one = ? AND id_user_two = ?) OR (id_user_one = ? AND id_user_two = ?))
			  AND id_annonce = ?
			LIMIT 1
		`, currentUserID, targetUserID, targetUserID, currentUserID, *annonceID).Scan(&existingID)
	case projetID != nil:
		existingErr = Conn.QueryRow(`
			SELECT id FROM DM_CONVERSATION
			WHERE ((id_user_one = ? AND id_user_two = ?) OR (id_user_one = ? AND id_user_two = ?))
			  AND id_projet = ?
			LIMIT 1
		`, currentUserID, targetUserID, targetUserID, currentUserID, *projetID).Scan(&existingID)
	default:
		existingErr = Conn.QueryRow(`
			SELECT id FROM DM_CONVERSATION
			WHERE ((id_user_one = ? AND id_user_two = ?) OR (id_user_one = ? AND id_user_two = ?))
			  AND id_annonce IS NULL AND id_projet IS NULL
			LIMIT 1
		`, currentUserID, targetUserID, targetUserID, currentUserID).Scan(&existingID)
	}

	if existingErr == nil {
		result.Allowed = true
		result.ConversationID = existingID
		result.Message = "Conversation existante."
		return result, nil
	}
	if existingErr != sql.ErrNoRows {
		return result, existingErr
	}

	isSubscriber, err := HasActiveDMSubscription(currentUserID)
	if err != nil {
		return result, err
	}
	result.IsSubscriber = isSubscriber

	isLinkedToItem := annonceID != nil || projetID != nil

	if !isLinkedToItem && !isSubscriber {
		result.Message = "L'abonnement DM Plus est requis pour contacter directement n'importe quel membre."
		return result, nil
	}

	if isLinkedToItem && !isSubscriber {
		used, err := CountDistinctVendorsContacted(currentUserID)
		if err != nil {
			return result, err
		}
		result.Used = used
		if used >= freeDMVendorLimit {
			result.Message = "Limite gratuite atteinte : vous pouvez contacter 5 vendeurs via annonces/projets. Passez a DM Plus pour continuer."
			return result, nil
		}
	}

	res, err := Conn.Exec(`
		INSERT INTO DM_CONVERSATION (id_user_one, id_user_two, id_annonce, id_projet, initiator_id)
		VALUES (?, ?, ?, ?, ?)
	`, currentUserID, targetUserID, annonceID, projetID, currentUserID)
	if err != nil {
		return result, err
	}
	id, _ := res.LastInsertId()
	result.Allowed = true
	result.ConversationID = int(id)
	result.Message = "Conversation creee."
	return result, nil
}

func CountDistinctVendorsContacted(userID int) (int, error) {
	var count int
	err := Conn.QueryRow(`
		SELECT COUNT(DISTINCT CASE
			WHEN id_user_one = ? THEN id_user_two
			ELSE id_user_one
		END)
		FROM DM_CONVERSATION
		WHERE initiator_id = ?
		  AND (id_annonce IS NOT NULL OR id_projet IS NOT NULL)
	`, userID, userID).Scan(&count)
	return count, err
}

func ListConversations(userID int) ([]models.DMConversation, error) {
	rows, err := Conn.Query(`
		SELECT
			c.id,
			CASE WHEN c.id_user_one = ? THEN c.id_user_two ELSE c.id_user_one END AS other_id,
			COALESCE(u.prenom, ''), COALESCE(u.nom, ''), COALESCE(u.role, ''), COALESCE(u.image_profil, ''),
			EXISTS (
				SELECT 1 FROM ABONNEMENT ab
				JOIN TYPE_ABONNEMENT tab ON tab.id = ab.id_type_abonnement
				WHERE ab.id_acheteur = u.id
				  AND ab.statut = 'Actif'
				  AND (ab.date_fin IS NULL OR ab.date_fin >= NOW())
				  AND tab.nom = 'DM Plus'
			) AS other_premium,
 
			c.id_annonce, COALESCE(a.titre, ''), COALESCE(a.prix, 0), a.id_vendeur,
			a.statut, a.id_acheteur,
 
			c.id_projet, COALESCE(p.titre, ''), COALESCE(p.prix, 0), p.id_createur,
			p.statut, p.id_acheteur,
 
			COALESCE(last_msg.contenu, '') AS last_message,
			COALESCE(DATE_FORMAT(last_msg.created_at, '%Y-%m-%dT%H:%i:%sZ'), '') AS last_message_at,
			COALESCE(unread.total, 0) AS unread_count,
			DATE_FORMAT(c.updated_at, '%Y-%m-%dT%H:%i:%sZ') AS updated_at
		FROM DM_CONVERSATION c
		JOIN UTILISATEUR u ON u.id = CASE WHEN c.id_user_one = ? THEN c.id_user_two ELSE c.id_user_one END
		LEFT JOIN ANNONCE a ON a.id = c.id_annonce
		LEFT JOIN PROJET_UPCYCLING p ON p.id = c.id_projet
		LEFT JOIN DM_MESSAGE last_msg ON last_msg.id = (
			SELECT m.id FROM DM_MESSAGE m WHERE m.id_conversation = c.id ORDER BY m.created_at DESC, m.id DESC LIMIT 1
		)
		LEFT JOIN (
			SELECT id_conversation, COUNT(*) AS total
			FROM DM_MESSAGE
			WHERE id_sender <> ? AND lu = FALSE
			GROUP BY id_conversation
		) unread ON unread.id_conversation = c.id
		WHERE c.id_user_one = ? OR c.id_user_two = ?
		ORDER BY c.updated_at DESC
	`, userID, userID, userID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
 
	conversations := []models.DMConversation{}
	for rows.Next() {
		var c models.DMConversation
		var prenom, nom string
 
		var annonceID sql.NullInt64
		var annonceTitle string
		var annoncePrix float64
		var annonceSellerID sql.NullInt64
		var annonceStatut sql.NullString
		var annonceBuyerID sql.NullInt64
 
		var projetID sql.NullInt64
		var projetTitle string
		var projetPrix float64
		var projetSellerID sql.NullInt64
		var projetStatut sql.NullString
		var projetBuyerID sql.NullInt64
 
		if err := rows.Scan(
			&c.ID, &c.OtherUserID, &prenom, &nom, &c.OtherUserRole, &c.OtherUserAvatar, &c.OtherUserPremium,
 
			&annonceID, &annonceTitle, &annoncePrix, &annonceSellerID,
			&annonceStatut, &annonceBuyerID,
 
			&projetID, &projetTitle, &projetPrix, &projetSellerID,
			&projetStatut, &projetBuyerID,
 
			&c.LastMessage, &c.LastMessageAt, &c.UnreadCount, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
 
		c.OtherUserName = strings.TrimSpace(prenom + " " + nom)
		if c.OtherUserName == "" {
			c.OtherUserName = fmt.Sprintf("Utilisateur #%d", c.OtherUserID)
		}
 
		if annonceID.Valid {
			val := int(annonceID.Int64)
			c.AnnonceID = &val
			c.AnnonceTitle = annonceTitle
			c.AnnoncePrice = annoncePrix
			if annonceSellerID.Valid {
				s := int(annonceSellerID.Int64)
				c.AnnonceSellerID = &s
			}
			if annonceStatut.Valid {
				c.AnnonceStatut = annonceStatut.String
			}
			if annonceBuyerID.Valid {
				b := int(annonceBuyerID.Int64)
				c.AnnonceBuyerID = &b
			}
		}
 
		if projetID.Valid {
			val := int(projetID.Int64)
			c.ProjetID = &val
			c.ProjetTitle = projetTitle
			c.ProjetPrice = projetPrix
			if projetSellerID.Valid {
				s := int(projetSellerID.Int64)
				c.ProjetSellerID = &s
			}
			if projetStatut.Valid {
				c.ProjetStatut = projetStatut.String
			}
			if projetBuyerID.Valid {
				b := int(projetBuyerID.Int64)
				c.ProjetBuyerID = &b
			}
		}
 
		conversations = append(conversations, c)
	}
	return conversations, nil
}

func GetConversationMessages(userID int, conversationID int) ([]models.DMMessage, error) {
	var allowed int
	if err := Conn.QueryRow(`
		SELECT COUNT(*) FROM DM_CONVERSATION
		WHERE id = ? AND (id_user_one = ? OR id_user_two = ?)
	`, conversationID, userID, userID).Scan(&allowed); err != nil {
		return nil, err
	}
	if allowed == 0 {
		return nil, fmt.Errorf("conversation introuvable ou accès refusé")
	}

	_, err := Conn.Exec(`
		UPDATE DM_MESSAGE 
		SET lu = TRUE 
		WHERE id_conversation = ? AND id_sender != ? AND lu = FALSE
	`, conversationID, userID)
	if err != nil {
		return nil, err
	}

	rows, err := Conn.Query(`
		SELECT 
			id, id_conversation, id_sender, contenu, lu, 
			COALESCE(DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ'), '')
		FROM DM_MESSAGE
		WHERE id_conversation = ?
		ORDER BY created_at ASC, id ASC
	`, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.DMMessage
	for rows.Next() {
		var m models.DMMessage
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.Content, &m.Read, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, nil
}

func SendDMMessage(userID int, conversationID int, content string) (*models.DMMessage, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, fmt.Errorf("message vide")
	}

	var allowed int
	var annonceID, projetID sql.NullInt64
	if err := Conn.QueryRow(`
		SELECT COUNT(*), MAX(id_annonce), MAX(id_projet) FROM DM_CONVERSATION
		WHERE id = ? AND (id_user_one = ? OR id_user_two = ?)
	`, conversationID, userID, userID).Scan(&allowed, &annonceID, &projetID); err != nil {
		return nil, err
	}
	if allowed == 0 {
		return nil, fmt.Errorf("conversation introuvable")
	}
	if !annonceID.Valid && !projetID.Valid {
		isSubscriber, err := HasActiveDMSubscription(userID)
		if err != nil {
			return nil, err
		}
		if !isSubscriber {
			return nil, fmt.Errorf("DM Plus est requis pour envoyer un message direct hors annonce/projet")
		}
	}

	res, err := Conn.Exec(`
		INSERT INTO DM_MESSAGE (id_conversation, id_sender, contenu)
		VALUES (?, ?, ?)
	`, conversationID, userID, content)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	_, _ = Conn.Exec("UPDATE DM_CONVERSATION SET updated_at = NOW() WHERE id = ?", conversationID)

	message := &models.DMMessage{ID: int(id), ConversationID: conversationID, SenderID: userID, Content: content, Read: false}
	_ = Conn.QueryRow("SELECT DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ') FROM DM_MESSAGE WHERE id = ?", id).Scan(&message.CreatedAt)
	return message, nil
}

func CreateDMOffer(userID int, conversationID int, amount float64) (*models.DMOffer, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("montant invalide")
	}

	var userOne, userTwo int
	var annonceID, projetID sql.NullInt64
	if err := Conn.QueryRow(`
		SELECT id_user_one, id_user_two, id_annonce, id_projet
		FROM DM_CONVERSATION
		WHERE id = ? AND (id_user_one = ? OR id_user_two = ?)
	`, conversationID, userID, userID).Scan(&userOne, &userTwo, &annonceID, &projetID); err != nil {
		return nil, fmt.Errorf("conversation introuvable")
	}
	if !annonceID.Valid && !projetID.Valid {
		return nil, fmt.Errorf("la negociation de prix est reservee aux conversations liees a une annonce ou un projet")
	}

	buyerID := userID
	otherID := userOne
	if userOne == userID {
		otherID = userTwo
	}

	var sellerID int
	var insertQuery string
	var insertArgs []interface{}

	if annonceID.Valid {
		if err := Conn.QueryRow("SELECT id_vendeur FROM ANNONCE WHERE id = ?", annonceID.Int64).Scan(&sellerID); err != nil {
			return nil, err
		}
		if sellerID == userID {
			return nil, fmt.Errorf("le vendeur ne peut pas proposer une offre sur sa propre annonce")
		}
		if otherID != sellerID {
			return nil, fmt.Errorf("offre impossible : le vendeur de l'annonce n'est pas dans cette conversation")
		}
		insertQuery = `INSERT INTO DM_OFFER (id_conversation, id_annonce, id_buyer, id_seller, amount, status) VALUES (?, ?, ?, ?, ?, 'En attente')`
		insertArgs = []interface{}{conversationID, annonceID.Int64, buyerID, sellerID, amount}
	} else {
		if err := Conn.QueryRow("SELECT id_createur FROM PROJET_UPCYCLING WHERE id = ?", projetID.Int64).Scan(&sellerID); err != nil {
			return nil, err
		}
		if sellerID == userID {
			return nil, fmt.Errorf("le createur ne peut pas proposer une offre sur son propre projet")
		}
		if otherID != sellerID {
			return nil, fmt.Errorf("offre impossible : le createur du projet n'est pas dans cette conversation")
		}
		insertQuery = `INSERT INTO DM_OFFER (id_conversation, id_projet, id_buyer, id_seller, amount, status) VALUES (?, ?, ?, ?, ?, 'En attente')`
		insertArgs = []interface{}{conversationID, projetID.Int64, buyerID, sellerID, amount}
	}

	res, err := Conn.Exec(insertQuery, insertArgs...)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	_, _ = Conn.Exec("UPDATE DM_CONVERSATION SET updated_at = NOW() WHERE id = ?", conversationID)
	_, _ = Conn.Exec("INSERT INTO DM_MESSAGE (id_conversation, id_sender, contenu) VALUES (?, ?, ?)", conversationID, userID, fmt.Sprintf("Offre proposee : %.2f EUR", amount))

	return GetDMOfferByID(int(id))
}

func GetDMOfferByID(id int) (*models.DMOffer, error) {
	var offer models.DMOffer
	var annonce, projet sql.NullInt64
	err := Conn.QueryRow(`
		SELECT id, id_conversation, id_annonce, id_projet, id_buyer, id_seller, amount, status,
		       DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ'),
		       DATE_FORMAT(updated_at, '%Y-%m-%dT%H:%i:%sZ')
		FROM DM_OFFER WHERE id = ?
	`, id).Scan(&offer.ID, &offer.ConversationID, &annonce, &projet, &offer.BuyerID, &offer.SellerID, &offer.Amount, &offer.Status, &offer.CreatedAt, &offer.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if annonce.Valid {
		v := int(annonce.Int64)
		offer.AnnonceID = &v
	}
	if projet.Valid {
		v := int(projet.Int64)
		offer.ProjetID = &v
	}
	return &offer, nil
}

func RespondDMOffer(userID int, offerID int, action string) (*models.DMOffer, *models.DMSale, error) {
	offer, err := GetDMOfferByID(offerID)
	if err != nil {
		return nil, nil, err
	}
	if offer.SellerID != userID {
		return nil, nil, fmt.Errorf("seul le vendeur peut valider ou refuser cette offre")
	}
	if offer.Status != "En attente" {
		return nil, nil, fmt.Errorf("cette offre a deja ete traitee")
	}

	status := "Refusee"
	message := "Offre refusee"
	if action == "accept" || action == "accepted" || action == "accepter" {
		status = "Acceptee"
		message = "Offre acceptee"
	}

	if _, err := Conn.Exec("UPDATE DM_OFFER SET status = ? WHERE id = ?", status, offerID); err != nil {
		return nil, nil, err
	}
	_, _ = Conn.Exec("INSERT INTO DM_MESSAGE (id_conversation, id_sender, contenu) VALUES (?, ?, ?)", offer.ConversationID, userID, message)
	_, _ = Conn.Exec("UPDATE DM_CONVERSATION SET updated_at = NOW() WHERE id = ?", offer.ConversationID)

	updated, err := GetDMOfferByID(offerID)
	if err != nil {
		return nil, nil, err
	}

	var sale *models.DMSale
	if status == "Acceptee" {
		var annonceValue, projetValue interface{}
		if offer.AnnonceID != nil {
			annonceValue = *offer.AnnonceID
		}
		if offer.ProjetID != nil {
			projetValue = *offer.ProjetID
		}

		res, err := Conn.Exec(`
			INSERT INTO DM_SALE (id_offer, id_conversation, id_annonce, id_projet, id_buyer, id_seller, amount, status)
			VALUES (?, ?, ?, ?, ?, ?, ?, 'Offre acceptee')
			ON DUPLICATE KEY UPDATE status = status
		`, offer.ID, offer.ConversationID, annonceValue, projetValue, offer.BuyerID, offer.SellerID, offer.Amount)
		if err != nil {
			return nil, nil, err
		}
		id, _ := res.LastInsertId()
		if id == 0 {
			sale, _ = GetDMSaleByOffer(offer.ID)
		} else {
			sale, _ = GetDMSaleByID(int(id))
		}
	}

	return updated, sale, nil
}

func GetDMSaleByOffer(offerID int) (*models.DMSale, error) {
	var id int
	if err := Conn.QueryRow("SELECT id FROM DM_SALE WHERE id_offer = ?", offerID).Scan(&id); err != nil {
		return nil, err
	}
	return GetDMSaleByID(id)
}

func GetDMSaleByID(id int) (*models.DMSale, error) {
	var sale models.DMSale
	var annonce, projet sql.NullInt64
	var received, reviewed sql.NullString
	err := Conn.QueryRow(`
		SELECT id, id_offer, id_conversation, id_annonce, id_projet, id_buyer, id_seller, amount, status,
		       COALESCE(DATE_FORMAT(received_at, '%Y-%m-%dT%H:%i:%sZ'), ''),
		       COALESCE(DATE_FORMAT(reviewed_at, '%Y-%m-%dT%H:%i:%sZ'), '')
		FROM DM_SALE WHERE id = ?
	`, id).Scan(&sale.ID, &sale.OfferID, &sale.ConversationID, &annonce, &projet, &sale.BuyerID, &sale.SellerID, &sale.Amount, &sale.Status, &received, &reviewed)
	if err != nil {
		return nil, err
	}
	if annonce.Valid {
		v := int(annonce.Int64)
		sale.AnnonceID = &v
	}
	if projet.Valid {
		v := int(projet.Int64)
		sale.ProjetID = &v
	}
	sale.ReceivedAt = received.String
	sale.ReviewedAt = reviewed.String
	return &sale, nil
}

func ConfirmDMSaleReception(userID int, saleID int) (*models.DMSale, error) {
	sale, err := GetDMSaleByID(saleID)
	if err != nil {
		return nil, err
	}
	if sale.BuyerID != userID {
		return nil, fmt.Errorf("seul l'acheteur peut confirmer la reception")
	}

	_, err = Conn.Exec(`
		UPDATE DM_SALE SET status = 'Recue', received_at = COALESCE(received_at, NOW())
		WHERE id = ?
	`, saleID)
	if err != nil {
		return nil, err
	}
	_, _ = Conn.Exec("INSERT INTO DM_MESSAGE (id_conversation, id_sender, contenu) VALUES (?, ?, 'Reception confirmee')", sale.ConversationID, userID)
	return GetDMSaleByID(saleID)
}

func ReviewDMSale(userID int, saleID int, note int, commentaire string) (*models.DMSale, error) {
	if note < 1 || note > 5 {
		return nil, fmt.Errorf("note invalide")
	}
	sale, err := GetDMSaleByID(saleID)
	if err != nil {
		return nil, err
	}
	if sale.BuyerID != userID {
		return nil, fmt.Errorf("seul l'acheteur peut noter cette vente")
	}
	if sale.Status != "Recue" && sale.Status != "Payee" {
		return nil, fmt.Errorf("confirmez la reception (ou le paiement pour un projet) avant de noter la vente")
	}

	if err := CreateAvis(userID, sale.SellerID, note, commentaire); err != nil {
		return nil, err
	}
	_, err = Conn.Exec("UPDATE DM_SALE SET status = 'Evaluee', reviewed_at = NOW() WHERE id = ?", saleID)
	if err != nil {
		return nil, err
	}
	_, _ = Conn.Exec("INSERT INTO DM_MESSAGE (id_conversation, id_sender, contenu) VALUES (?, ?, 'Vente notee')", sale.ConversationID, userID)
	return GetDMSaleByID(saleID)
}

func GetDMThreadState(userID int, conversationID int) (*models.DMThreadState, error) {
	var allowed int
	if err := Conn.QueryRow(`
		SELECT COUNT(*) FROM DM_CONVERSATION
		WHERE id = ? AND (id_user_one = ? OR id_user_two = ?)
	`, conversationID, userID, userID).Scan(&allowed); err != nil {
		return nil, err
	}
	if allowed == 0 {
		return nil, fmt.Errorf("conversation introuvable")
	}

	offers, err := ListDMOffers(conversationID)
	if err != nil {
		return nil, err
	}
	sales, err := ListDMSales(conversationID)
	if err != nil {
		return nil, err
	}
	return &models.DMThreadState{Offers: offers, Sales: sales}, nil
}

func ListDMOffers(conversationID int) ([]models.DMOffer, error) {
	rows, err := Conn.Query(`SELECT id FROM DM_OFFER WHERE id_conversation = ? ORDER BY created_at DESC, id DESC`, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	offers := []models.DMOffer{}
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		offer, err := GetDMOfferByID(id)
		if err != nil {
			return nil, err
		}
		offers = append(offers, *offer)
	}
	return offers, rows.Err()
}

func ListDMSales(conversationID int) ([]models.DMSale, error) {
	rows, err := Conn.Query(`SELECT id FROM DM_SALE WHERE id_conversation = ? ORDER BY created_at DESC, id DESC`, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	sales := []models.DMSale{}
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		sale, err := GetDMSaleByID(id)
		if err != nil {
			return nil, err
		}
		sales = append(sales, *sale)
	}
	return sales, rows.Err()
}