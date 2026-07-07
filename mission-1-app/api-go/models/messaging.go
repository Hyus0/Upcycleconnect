package models

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
	AnnonceSellerID  *int    `json:"annonce_seller_id,omitempty"`
	AnnonceStatut    string  `json:"annonce_statut"`
	AnnonceBuyerID   *int    `json:"annonce_acheteur_id,omitempty"`

	ProjetID         *int    `json:"projet_id,omitempty"`
	ProjetTitle      string  `json:"projet_title,omitempty"`
	ProjetPrice      float64 `json:"projet_price,omitempty"`
	ProjetSellerID   *int    `json:"projet_seller_id,omitempty"`
	ProjetStatut     string  `json:"projet_statut"`          
	ProjetBuyerID    *int    `json:"projet_acheteur_id,omitempty"` 

	LastMessage   string `json:"last_message"`
	LastMessageAt string `json:"last_message_at"`
	UnreadCount   int    `json:"unread_count"`
	UpdatedAt     string `json:"updated_at"`
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
	ProjetID       *int    `json:"projet_id,omitempty"`
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
	ProjetID       *int    `json:"projet_id,omitempty"`
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