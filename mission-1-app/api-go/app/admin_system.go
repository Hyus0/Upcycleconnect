package app

import (
	"encoding/json"
	"net/http"
	"upcycleconnect/api-go/db"
)

func AdminSystemOverview(w http.ResponseWriter, r *http.Request) {
	type metric struct {
		Key   string `json:"key"`
		Label string `json:"label"`
		Value int    `json:"value"`
	}

	count := func(query string) int {
		var value int
		_ = db.Conn.QueryRow(query).Scan(&value)
		return value
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"metrics": []metric{
			{Key: "users", Label: "Utilisateurs", Value: count("SELECT COUNT(*) FROM UTILISATEUR")},
			{Key: "premium", Label: "Abonnes DM Plus", Value: count("SELECT COUNT(*) FROM ABONNEMENT a JOIN TYPE_ABONNEMENT t ON t.id = a.id_type_abonnement WHERE a.statut = 'Actif' AND t.nom = 'DM Plus' AND (a.date_fin IS NULL OR a.date_fin >= NOW())")},
			{Key: "annonces", Label: "Annonces", Value: count("SELECT COUNT(*) FROM ANNONCE")},
			{Key: "conversations", Label: "Conversations privees", Value: count("SELECT COUNT(*) FROM DM_CONVERSATION")},
			{Key: "offers", Label: "Offres negociees", Value: count("SELECT COUNT(*) FROM DM_OFFER")},
			{Key: "sales", Label: "Ventes suivies", Value: count("SELECT COUNT(*) FROM DM_SALE")},
			{Key: "reviews", Label: "Avis", Value: count("SELECT COUNT(*) FROM AVIS")},
			{Key: "orders", Label: "Commandes", Value: count("SELECT COUNT(*) FROM COMMANDE")},
		},
	})
}

func AdminListSubscriptions(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Conn.Query(`
		SELECT a.id, a.id_acheteur, CONCAT_WS(' ', u.prenom, u.nom), u.mail, t.nom,
		       t.prix_ht, a.statut,
		       COALESCE(DATE_FORMAT(a.date_debut, '%Y-%m-%d %H:%i:%s'), ''),
		       COALESCE(DATE_FORMAT(a.date_fin, '%Y-%m-%d %H:%i:%s'), '')
		FROM ABONNEMENT a
		JOIN UTILISATEUR u ON u.id = a.id_acheteur
		JOIN TYPE_ABONNEMENT t ON t.id = a.id_type_abonnement
		ORDER BY a.date_debut DESC, a.id DESC
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	items := []map[string]any{}
	for rows.Next() {
		var id, userID int
		var name, email, plan, status, start, end string
		var price float64
		if err := rows.Scan(&id, &userID, &name, &email, &plan, &price, &status, &start, &end); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		items = append(items, map[string]any{
			"id": id, "user_id": userID, "user": name, "email": email,
			"plan": plan, "price": price, "status": status, "date_debut": start, "date_fin": end,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func AdminListMessages(w http.ResponseWriter, r *http.Request) {
	conversations, err := db.Conn.Query(`
		SELECT c.id, c.id_user_one, CONCAT_WS(' ', u1.prenom, u1.nom),
		       c.id_user_two, CONCAT_WS(' ', u2.prenom, u2.nom),
		       c.id_annonce, COALESCE(a.titre, 'Discussion directe'),
		       DATE_FORMAT(c.updated_at, '%Y-%m-%d %H:%i:%s')
		FROM DM_CONVERSATION c
		JOIN UTILISATEUR u1 ON u1.id = c.id_user_one
		JOIN UTILISATEUR u2 ON u2.id = c.id_user_two
		LEFT JOIN ANNONCE a ON a.id = c.id_annonce
		ORDER BY c.updated_at DESC
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conversations.Close()

	items := []map[string]any{}
	for conversations.Next() {
		var id, userOne, userTwo int
		var annonceID any
		var userOneName, userTwoName, title, updated string
		if err := conversations.Scan(&id, &userOne, &userOneName, &userTwo, &userTwoName, &annonceID, &title, &updated); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		items = append(items, map[string]any{
			"id": id, "user_one_id": userOne, "user_one": userOneName,
			"user_two_id": userTwo, "user_two": userTwoName, "annonce_id": annonceID,
			"title": title, "updated_at": updated,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func AdminListCommerce(w http.ResponseWriter, r *http.Request) {
	offers, _ := db.Conn.Query(`
		SELECT o.id, o.id_conversation, o.id_annonce, o.id_buyer, o.id_seller, o.amount, o.status,
		       DATE_FORMAT(o.created_at, '%Y-%m-%d %H:%i:%s')
		FROM DM_OFFER o ORDER BY o.created_at DESC
	`)
	defer func() {
		if offers != nil {
			offers.Close()
		}
	}()

	offerItems := []map[string]any{}
	if offers != nil {
		for offers.Next() {
			var id, conversationID, buyerID, sellerID int
			var annonceID any
			var amount float64
			var status, created string
			if err := offers.Scan(&id, &conversationID, &annonceID, &buyerID, &sellerID, &amount, &status, &created); err == nil {
				offerItems = append(offerItems, map[string]any{
					"id": id, "conversation_id": conversationID, "annonce_id": annonceID,
					"buyer_id": buyerID, "seller_id": sellerID, "amount": amount, "status": status, "created_at": created,
				})
			}
		}
	}

	sales, _ := db.Conn.Query(`
		SELECT id, id_offer, id_conversation, id_annonce, id_buyer, id_seller, amount, status,
		       COALESCE(DATE_FORMAT(received_at, '%Y-%m-%d %H:%i:%s'), ''),
		       COALESCE(DATE_FORMAT(reviewed_at, '%Y-%m-%d %H:%i:%s'), '')
		FROM DM_SALE ORDER BY updated_at DESC
	`)
	defer func() {
		if sales != nil {
			sales.Close()
		}
	}()

	saleItems := []map[string]any{}
	if sales != nil {
		for sales.Next() {
			var id, offerID, conversationID, buyerID, sellerID int
			var annonceID any
			var amount float64
			var status, received, reviewed string
			if err := sales.Scan(&id, &offerID, &conversationID, &annonceID, &buyerID, &sellerID, &amount, &status, &received, &reviewed); err == nil {
				saleItems = append(saleItems, map[string]any{
					"id": id, "offer_id": offerID, "conversation_id": conversationID, "annonce_id": annonceID,
					"buyer_id": buyerID, "seller_id": sellerID, "amount": amount, "status": status,
					"received_at": received, "reviewed_at": reviewed,
				})
			}
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{"offers": offerItems, "sales": saleItems})
}

func AdminRawDump(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	allowed := map[string]string{
		"users":        "SELECT id, prenom, nom, mail, role, statut FROM UTILISATEUR ORDER BY id DESC LIMIT 200",
		"annonces":     "SELECT id, titre, type, prix, statut, est_valide, id_vendeur FROM ANNONCE ORDER BY id DESC LIMIT 200",
		"events":       "SELECT id, titre, date_evenement, statut, capacite_max FROM EVENEMENT ORDER BY id DESC LIMIT 200",
		"formations":   "SELECT id, titre, type, prix_unitaire, statut, est_valide, date_debut FROM FORMATION ORDER BY id DESC LIMIT 200",
		"orders":       "SELECT id, id_utilisateur, montant_total, statut, date_commande FROM COMMANDE ORDER BY id DESC LIMIT 200",
		"transactions": "SELECT id, id_acheteur, id_commande, montant_total, statut_paiement, date_transaction FROM `TRANSACTION` ORDER BY id DESC LIMIT 200",
		"reviews":      "SELECT id, id_auteur, id_cible, note, commentaire, date_creation FROM AVIS ORDER BY id DESC LIMIT 200",
	}
	query, ok := allowed[target]
	if !ok {
		http.Error(w, "target invalide", http.StatusBadRequest)
		return
	}

	rows, err := db.Conn.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	result := []map[string]any{}
	for rows.Next() {
		values := make([]any, len(columns))
		ptrs := make([]any, len(columns))
		for i := range values {
			ptrs[i] = &values[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		row := map[string]any{}
		for i, col := range columns {
			if b, ok := values[i].([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = values[i]
			}
		}
		result = append(result, row)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"items": result})
}
