package db

import (
	"database/sql"
	"fmt"
	"strings"
)

type DMConversation struct {
	ID               int     `json:"id"`
	OtherUserID      int     `json:"other_user_id"`
	OtherUserName    string  `json:"other_user_name"`
	OtherUserRole    string  `json:"other_user_role"`
	OtherUserAvatar  string  `json:"other_user_avatar"`
	OtherUserPremium bool    `json:"other_user_premium"`
	AnnonceID        *int    `json:"annonce_id,omitempty"`
	AnnonceTitle     string  `json:"annonce_title"`
	AnnoncePrice     float64 `json:"annonce_price"`
	LastMessage      string  `json:"last_message"`
	LastMessageAt    string  `json:"last_message_at"`
	UnreadCount      int     `json:"unread_count"`
	UpdatedAt        string  `json:"updated_at"`
}

type DMMessage struct {
	ID             int    `json:"id"`
	ConversationID int    `json:"conversation_id"`
	SenderID       int    `json:"sender_id"`
	Content        string `json:"content"`
	Read           bool   `json:"read"`
	CreatedAt      string `json:"created_at"`
}

type DMStartResult struct {
	ConversationID int    `json:"conversation_id"`
	Allowed        bool   `json:"allowed"`
	IsSubscriber   bool   `json:"is_subscriber"`
	Used           int    `json:"used"`
	Limit          int    `json:"limit"`
	Message        string `json:"message"`
}

type DMOffer struct {
	ID             int     `json:"id"`
	ConversationID int     `json:"conversation_id"`
	AnnonceID      *int    `json:"annonce_id,omitempty"`
	BuyerID        int     `json:"buyer_id"`
	SellerID       int     `json:"seller_id"`
	Amount         float64 `json:"amount"`
	Status         string  `json:"status"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

type DMSale struct {
	ID             int     `json:"id"`
	OfferID        int     `json:"offer_id"`
	ConversationID int     `json:"conversation_id"`
	AnnonceID      *int    `json:"annonce_id,omitempty"`
	BuyerID        int     `json:"buyer_id"`
	SellerID       int     `json:"seller_id"`
	Amount         float64 `json:"amount"`
	Status         string  `json:"status"`
	ReceivedAt     string  `json:"received_at"`
	ReviewedAt     string  `json:"reviewed_at"`
}

type DMThreadState struct {
	Offers []DMOffer `json:"offers"`
	Sales  []DMSale  `json:"sales"`
}

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

func EnsureDMSubscription(userID int) error {
	_, err := Conn.Exec(`
		INSERT INTO ABONNEMENT (id_acheteur, id_type_abonnement, date_debut, date_fin, statut, stripe_subscription_id)
		SELECT ?, id, NOW(), DATE_ADD(NOW(), INTERVAL duree_mois MONTH), 'Actif', 'local-checkout'
		FROM TYPE_ABONNEMENT
		WHERE nom = 'DM Plus'
		ON DUPLICATE KEY UPDATE statut = statut
	`, userID)
	return err
}

func StartConversation(currentUserID int, targetUserID int, annonceID *int) (DMStartResult, error) {
	result := DMStartResult{Allowed: false, Limit: freeDMVendorLimit}
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
	if annonceID != nil {
		existingErr = Conn.QueryRow(`
			SELECT id FROM DM_CONVERSATION
			WHERE ((id_user_one = ? AND id_user_two = ?) OR (id_user_one = ? AND id_user_two = ?))
			  AND id_annonce = ?
			LIMIT 1
		`, currentUserID, targetUserID, targetUserID, currentUserID, *annonceID).Scan(&existingID)
	} else {
		existingErr = Conn.QueryRow(`
			SELECT id FROM DM_CONVERSATION
			WHERE ((id_user_one = ? AND id_user_two = ?) OR (id_user_one = ? AND id_user_two = ?))
			  AND id_annonce IS NULL
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

	if annonceID == nil && !isSubscriber {
		result.Message = "L'abonnement DM Plus est requis pour contacter directement n'importe quel membre."
		return result, nil
	}

	used, err := CountDistinctAnnonceVendorsContacted(currentUserID)
	if err != nil {
		return result, err
	}
	result.Used = used

	if !isSubscriber && used >= freeDMVendorLimit {
		result.Message = "Limite gratuite atteinte : vous pouvez contacter 5 vendeurs via des annonces. Passez a DM Plus pour continuer."
		return result, nil
	}

	res, err := Conn.Exec(`
		INSERT INTO DM_CONVERSATION (id_user_one, id_user_two, id_annonce, initiator_id)
		VALUES (?, ?, ?, ?)
	`, currentUserID, targetUserID, annonceID, currentUserID)
	if err != nil {
		return result, err
	}
	id, _ := res.LastInsertId()
	result.Allowed = true
	result.ConversationID = int(id)
	result.Message = "Conversation creee."
	return result, nil
}

func CountDistinctAnnonceVendorsContacted(userID int) (int, error) {
	var count int
	err := Conn.QueryRow(`
		SELECT COUNT(DISTINCT CASE
			WHEN id_user_one = ? THEN id_user_two
			ELSE id_user_one
		END)
		FROM DM_CONVERSATION
		WHERE initiator_id = ?
		  AND id_annonce IS NOT NULL
	`, userID, userID).Scan(&count)
	return count, err
}

func ListConversations(userID int) ([]DMConversation, error) {
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
			c.id_annonce,
			COALESCE(a.titre, 'Discussion directe') AS annonce_title,
			COALESCE(a.prix, 0) AS annonce_price,
			COALESCE(last_msg.contenu, '') AS last_message,
			COALESCE(DATE_FORMAT(last_msg.created_at, '%Y-%m-%dT%H:%i:%sZ'), '') AS last_message_at,
			COALESCE(unread.total, 0) AS unread_count,
			DATE_FORMAT(c.updated_at, '%Y-%m-%dT%H:%i:%sZ') AS updated_at
		FROM DM_CONVERSATION c
		JOIN UTILISATEUR u ON u.id = CASE WHEN c.id_user_one = ? THEN c.id_user_two ELSE c.id_user_one END
		LEFT JOIN ANNONCE a ON a.id = c.id_annonce
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

	conversations := []DMConversation{}
	for rows.Next() {
		var c DMConversation
		var prenom, nom string
		if err := rows.Scan(&c.ID, &c.OtherUserID, &prenom, &nom, &c.OtherUserRole, &c.OtherUserAvatar, &c.OtherUserPremium, &c.AnnonceID, &c.AnnonceTitle, &c.AnnoncePrice, &c.LastMessage, &c.LastMessageAt, &c.UnreadCount, &c.UpdatedAt); err != nil {
			return nil, err
		}
		c.OtherUserName = strings.TrimSpace(prenom + " " + nom)
		if c.OtherUserName == "" {
			c.OtherUserName = fmt.Sprintf("Utilisateur #%d", c.OtherUserID)
		}
		conversations = append(conversations, c)
	}
	return conversations, nil
}

func GetConversationMessages(userID int, conversationID int) ([]DMMessage, error) {
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

	_, _ = Conn.Exec("UPDATE DM_MESSAGE SET lu = TRUE WHERE id_conversation = ? AND id_sender <> ?", conversationID, userID)

	rows, err := Conn.Query(`
		SELECT id, id_conversation, id_sender, contenu, lu, DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ')
		FROM DM_MESSAGE
		WHERE id_conversation = ?
		ORDER BY created_at ASC, id ASC
	`, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []DMMessage{}
	for rows.Next() {
		var m DMMessage
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.Content, &m.Read, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func SendDMMessage(userID int, conversationID int, content string) (*DMMessage, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, fmt.Errorf("message vide")
	}

	var allowed int
	var annonceID sql.NullInt64
	if err := Conn.QueryRow(`
		SELECT COUNT(*), MAX(id_annonce) FROM DM_CONVERSATION
		WHERE id = ? AND (id_user_one = ? OR id_user_two = ?)
	`, conversationID, userID, userID).Scan(&allowed, &annonceID); err != nil {
		return nil, err
	}
	if allowed == 0 {
		return nil, fmt.Errorf("conversation introuvable")
	}
	if !annonceID.Valid {
		isSubscriber, err := HasActiveDMSubscription(userID)
		if err != nil {
			return nil, err
		}
		if !isSubscriber {
			return nil, fmt.Errorf("DM Plus est requis pour envoyer un message direct hors annonce")
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

	message := &DMMessage{ID: int(id), ConversationID: conversationID, SenderID: userID, Content: content, Read: false}
	_ = Conn.QueryRow("SELECT DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ') FROM DM_MESSAGE WHERE id = ?", id).Scan(&message.CreatedAt)
	return message, nil
}

func CreateDMOffer(userID int, conversationID int, amount float64) (*DMOffer, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("montant invalide")
	}

	var userOne, userTwo int
	var annonceID sql.NullInt64
	if err := Conn.QueryRow(`
		SELECT id_user_one, id_user_two, id_annonce
		FROM DM_CONVERSATION
		WHERE id = ? AND (id_user_one = ? OR id_user_two = ?)
	`, conversationID, userID, userID).Scan(&userOne, &userTwo, &annonceID); err != nil {
		return nil, fmt.Errorf("conversation introuvable")
	}
	if !annonceID.Valid {
		return nil, fmt.Errorf("la negociation de prix est reservee aux conversations liees a une annonce")
	}

	var sellerID int
	if err := Conn.QueryRow("SELECT id_vendeur FROM ANNONCE WHERE id = ?", annonceID.Int64).Scan(&sellerID); err != nil {
		return nil, err
	}
	if sellerID == userID {
		return nil, fmt.Errorf("le vendeur ne peut pas proposer une offre sur sa propre annonce")
	}

	buyerID := userID
	otherID := userOne
	if userOne == userID {
		otherID = userTwo
	}
	if otherID != sellerID {
		return nil, fmt.Errorf("offre impossible : le vendeur de l'annonce n'est pas dans cette conversation")
	}

	res, err := Conn.Exec(`
		INSERT INTO DM_OFFER (id_conversation, id_annonce, id_buyer, id_seller, amount, status)
		VALUES (?, ?, ?, ?, ?, 'En attente')
	`, conversationID, annonceID.Int64, buyerID, sellerID, amount)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	_, _ = Conn.Exec("UPDATE DM_CONVERSATION SET updated_at = NOW() WHERE id = ?", conversationID)
	_, _ = Conn.Exec("INSERT INTO DM_MESSAGE (id_conversation, id_sender, contenu) VALUES (?, ?, ?)", conversationID, userID, fmt.Sprintf("Offre proposee : %.2f EUR", amount))

	return GetDMOfferByID(int(id))
}

func GetDMOfferByID(id int) (*DMOffer, error) {
	var offer DMOffer
	var annonce sql.NullInt64
	err := Conn.QueryRow(`
		SELECT id, id_conversation, id_annonce, id_buyer, id_seller, amount, status,
		       DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ'),
		       DATE_FORMAT(updated_at, '%Y-%m-%dT%H:%i:%sZ')
		FROM DM_OFFER WHERE id = ?
	`, id).Scan(&offer.ID, &offer.ConversationID, &annonce, &offer.BuyerID, &offer.SellerID, &offer.Amount, &offer.Status, &offer.CreatedAt, &offer.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if annonce.Valid {
		value := int(annonce.Int64)
		offer.AnnonceID = &value
	}
	return &offer, nil
}

func RespondDMOffer(userID int, offerID int, action string) (*DMOffer, *DMSale, error) {
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

	var sale *DMSale
	if status == "Acceptee" {
		var annonceValue any
		if offer.AnnonceID != nil {
			annonceValue = *offer.AnnonceID
		}
		res, err := Conn.Exec(`
			INSERT INTO DM_SALE (id_offer, id_conversation, id_annonce, id_buyer, id_seller, amount, status)
			VALUES (?, ?, ?, ?, ?, ?, 'Offre acceptee')
			ON DUPLICATE KEY UPDATE status = status
		`, offer.ID, offer.ConversationID, annonceValue, offer.BuyerID, offer.SellerID, offer.Amount)
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

func GetDMSaleByOffer(offerID int) (*DMSale, error) {
	var id int
	if err := Conn.QueryRow("SELECT id FROM DM_SALE WHERE id_offer = ?", offerID).Scan(&id); err != nil {
		return nil, err
	}
	return GetDMSaleByID(id)
}

func GetDMSaleByID(id int) (*DMSale, error) {
	var sale DMSale
	var annonce sql.NullInt64
	var received, reviewed sql.NullString
	err := Conn.QueryRow(`
		SELECT id, id_offer, id_conversation, id_annonce, id_buyer, id_seller, amount, status,
		       COALESCE(DATE_FORMAT(received_at, '%Y-%m-%dT%H:%i:%sZ'), ''),
		       COALESCE(DATE_FORMAT(reviewed_at, '%Y-%m-%dT%H:%i:%sZ'), '')
		FROM DM_SALE WHERE id = ?
	`, id).Scan(&sale.ID, &sale.OfferID, &sale.ConversationID, &annonce, &sale.BuyerID, &sale.SellerID, &sale.Amount, &sale.Status, &received, &reviewed)
	if err != nil {
		return nil, err
	}
	if annonce.Valid {
		value := int(annonce.Int64)
		sale.AnnonceID = &value
	}
	sale.ReceivedAt = received.String
	sale.ReviewedAt = reviewed.String
	return &sale, nil
}

func ConfirmDMSaleReception(userID int, saleID int) (*DMSale, error) {
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

func ReviewDMSale(userID int, saleID int, note int, commentaire string) (*DMSale, error) {
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
	if sale.Status != "Recue" {
		return nil, fmt.Errorf("confirmez la reception avant de noter la vente")
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

func GetDMThreadState(userID int, conversationID int) (*DMThreadState, error) {
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
	return &DMThreadState{Offers: offers, Sales: sales}, nil
}

func ListDMOffers(conversationID int) ([]DMOffer, error) {
	rows, err := Conn.Query(`
		SELECT id FROM DM_OFFER WHERE id_conversation = ? ORDER BY created_at DESC, id DESC
	`, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	offers := []DMOffer{}
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

func ListDMSales(conversationID int) ([]DMSale, error) {
	rows, err := Conn.Query(`
		SELECT id FROM DM_SALE WHERE id_conversation = ? ORDER BY created_at DESC, id DESC
	`, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	sales := []DMSale{}
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
