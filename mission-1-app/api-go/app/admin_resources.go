package app

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"upcycleconnect/api-go/db"
)

type AdminResourceField struct {
	Name     string   `json:"name"`
	Label    string   `json:"label"`
	Type     string   `json:"type"`
	Required bool     `json:"required,omitempty"`
	ReadOnly bool     `json:"readOnly,omitempty"`
	Options  []string `json:"options,omitempty"`
}

type AdminResource struct {
	Key         string               `json:"key"`
	Label       string               `json:"label"`
	Table       string               `json:"-"`
	PrimaryKeys []string             `json:"primaryKeys"`
	Fields      []AdminResourceField `json:"fields"`
	ReadOnly    bool                 `json:"readOnly,omitempty"`
}

func adminResourceDefinitions() map[string]AdminResource {
	return map[string]AdminResource{
		"formations": resource("formations", "Formations", "FORMATION", []string{"id"}, []AdminResourceField{
			field("id", "ID", "number", false, true),
			enumField("type", "Type", []string{"Atelier", "Cours", "Webinaire"}),
			field("titre", "Titre", "text", true, false),
			field("description", "Description", "textarea", false, false),
			field("id_formateur", "ID formateur", "number", true, false),
			field("capacite_max", "Capacite max", "number", false, false),
			enumField("est_valide", "Validation", []string{"En attente", "Valide", "Refuse"}),
			field("date_debut", "Date debut", "datetime", false, false),
			field("date_fin", "Date fin", "datetime", false, false),
			enumField("statut", "Statut", []string{"Ouvert", "Complet", "Termine", "Annule"}),
			field("prix_unitaire", "Prix", "decimal", false, false),
			field("adresse", "Adresse", "text", false, false),
			field("ville", "Ville", "text", false, false),
			field("code_postal", "Code postal", "text", false, false),
		}),
		"projects": resource("projects", "Projets upcycling", "PROJET_UPCYCLING", []string{"id"}, []AdminResourceField{
			field("id", "ID", "number", false, true),
			field("id_createur", "ID createur", "number", true, false),
			field("image_url", "Image", "text", true, false),
			field("titre", "Titre", "text", true, false),
			field("description_courte", "Description courte", "textarea", false, false),
			field("score_impact", "Score impact", "decimal", false, false),
			field("nb_vues", "Vues", "number", false, false),
			field("nb_likes", "Likes", "number", false, false),
			field("co2_evite_kg", "CO2 evite kg", "decimal", false, false),
			field("visible_public", "Visible public", "boolean", false, false),
		}),
		"project_steps": resource("project_steps", "Etapes projet", "ETAPE", []string{"id"}, []AdminResourceField{
			field("id", "ID", "number", false, true),
			field("id_projet", "ID projet", "number", true, false),
			field("numero_ordre", "Ordre", "number", true, false),
			field("titre", "Titre", "text", false, false),
			field("description", "Description", "textarea", false, false),
			field("image_url", "Image", "text", false, false),
		}),
		"tips": resource("tips", "Conseils / tips", "TIPS", []string{"id"}, []AdminResourceField{
			field("id", "ID", "number", false, true),
			field("id_createur", "ID createur", "number", true, false),
			field("titre", "Titre", "text", true, false),
			field("description", "Description", "textarea", true, false),
			field("video_url", "Video", "text", false, false),
			enumField("role_cible", "Role cible", []string{"Particulier", "Prestataire", "Salarie"}),
			field("actif", "Actif", "boolean", false, false),
		}),
		"forums": resource("forums", "Discussions forum", "FORUM", []string{"id"}, []AdminResourceField{
			field("id", "ID", "number", false, true),
			field("id_utilisateur", "ID auteur", "number", true, false),
			field("id_salon", "ID salon", "number", false, false),
			field("titre", "Titre", "text", false, false),
			field("sujet", "Sujet", "textarea", false, false),
			field("ouvert", "Ouvert", "boolean", false, false),
		}),
		"forum_salons": resource("forum_salons", "Salons forum", "FORUM_SALON", []string{"id"}, []AdminResourceField{
			field("id", "ID", "number", false, true),
			field("nom", "Nom", "text", true, false),
			field("description", "Description", "textarea", false, false),
		}),
		"forum_messages": resource("forum_messages", "Messages forum", "FORUM_MESSAGE", []string{"id"}, []AdminResourceField{
			field("id", "ID", "number", false, true),
			field("id_utilisateur", "ID auteur", "number", true, false),
			field("id_forum", "ID discussion", "number", true, false),
			field("contenu", "Contenu", "textarea", false, false),
		}),
		"sites": resource("sites", "Sites logistiques", "SITE", []string{"id"}, []AdminResourceField{
			field("id", "ID", "number", false, true),
			field("nom", "Nom", "text", false, false),
			field("ville", "Ville", "text", false, false),
			field("code_postal", "Code postal", "text", false, false),
			field("adresse", "Adresse", "text", false, false),
			field("telephone", "Telephone", "text", false, false),
			enumField("type", "Type", []string{"Decheterie", "Point de collecte", "Association"}),
			field("actif", "Actif", "boolean", false, false),
		}),
		"containers": resource("containers", "Conteneurs", "CONTENEUR", []string{"id"}, []AdminResourceField{
			field("id", "ID", "number", false, true),
			field("id_site", "ID site", "number", true, false),
			enumField("type_dechet", "Type dechet", []string{"Verre", "Plastique", "Metal", "Papier", "Electronique"}),
			enumField("statut", "Statut", []string{"Operationnel", "Plein", "Maintenance"}),
			field("capacite_max_kg", "Capacite kg", "decimal", false, false),
			field("niveau_remplissage", "Remplissage", "decimal", false, false),
		}),
		"lockers": resource("lockers", "Casiers", "CASIER", []string{"id"}, []AdminResourceField{
			field("id", "ID", "number", false, true),
			field("id_conteneur", "ID conteneur", "number", true, false),
			field("numero_casier", "Numero", "text", true, false),
			enumField("taille", "Taille", []string{"Petit", "Moyen", "Grand"}),
			enumField("statut", "Statut", []string{"Libre", "Reserve", "Occupe", "Maintenance"}),
		}),
		"subscription_types": resource("subscription_types", "Types abonnement", "TYPE_ABONNEMENT", []string{"id"}, []AdminResourceField{
			field("id", "ID", "number", false, true),
			field("nom", "Nom", "text", true, false),
			field("description", "Description", "textarea", false, false),
			field("prix_ht", "Prix HT", "decimal", true, false),
			field("duree_mois", "Duree mois", "number", true, false),
		}),
		"user_subscriptions": resource("user_subscriptions", "Abonnements utilisateurs", "ABONNEMENT", []string{"id"}, []AdminResourceField{
			field("id", "ID", "number", false, true),
			field("id_acheteur", "ID utilisateur", "number", true, false),
			field("id_type_abonnement", "ID type abonnement", "number", true, false),
			field("date_debut", "Date debut", "datetime", false, false),
			field("date_fin", "Date fin", "datetime", false, false),
			enumField("statut", "Statut", []string{"Actif", "Expire", "Resilie"}),
			field("stripe_subscription_id", "Stripe subscription", "text", false, false),
		}),
		"orders": resource("orders", "Commandes", "COMMANDE", []string{"id"}, []AdminResourceField{
			field("id", "ID", "number", false, true),
			field("id_utilisateur", "ID utilisateur", "number", true, false),
			field("montant_total", "Montant total", "decimal", true, false),
			enumField("statut", "Statut", []string{"En attente", "Payee", "Annulee"}),
		}),
		"cart_items": resource("cart_items", "Panier", "PANIER_ITEM", []string{"id"}, []AdminResourceField{
			field("id", "ID", "number", false, true),
			field("id_utilisateur", "ID utilisateur", "number", true, false),
			enumField("type_item", "Type", []string{"Formation", "Annonce", "Abonnement", "Evenement"}),
			field("reference_id", "Reference", "number", true, false),
			field("prix_unitaire", "Prix", "decimal", true, false),
		}),
		"invoices": resource("invoices", "Factures", "FACTURE", []string{"id"}, []AdminResourceField{
			field("id", "ID", "number", false, true),
			field("id_transaction", "ID transaction", "number", true, false),
			field("numero_facture", "Numero", "text", false, false),
			enumField("type_paiement", "Paiement", []string{"Carte", "Espece", "Virement"}),
		}),
		"dm_conversations": resource("dm_conversations", "Conversations", "DM_CONVERSATION", []string{"id"}, []AdminResourceField{
			field("id", "ID", "number", false, true),
			field("id_user_one", "Utilisateur 1", "number", true, false),
			field("id_user_two", "Utilisateur 2", "number", true, false),
			field("id_annonce", "ID annonce", "number", false, false),
			field("initiator_id", "Initiateur", "number", true, false),
		}),
		"dm_messages": resource("dm_messages", "Messages prives", "DM_MESSAGE", []string{"id"}, []AdminResourceField{
			field("id", "ID", "number", false, true),
			field("id_conversation", "Conversation", "number", true, false),
			field("id_sender", "Expediteur", "number", true, false),
			field("contenu", "Contenu", "textarea", true, false),
			field("lu", "Lu", "boolean", false, false),
		}),
		"dm_offers": resource("dm_offers", "Offres DM", "DM_OFFER", []string{"id"}, []AdminResourceField{
			field("id", "ID", "number", false, true),
			field("id_conversation", "Conversation", "number", true, false),
			field("id_annonce", "Annonce", "number", false, false),
			field("id_buyer", "Acheteur", "number", true, false),
			field("id_seller", "Vendeur", "number", true, false),
			field("amount", "Montant", "decimal", true, false),
			enumField("status", "Statut", []string{"En attente", "Acceptee", "Refusee", "Annulee"}),
		}),
		"languages": resource("languages", "Langues", "LANGUE", []string{"id"}, []AdminResourceField{
			field("id", "ID", "number", false, true),
			field("code", "Code", "text", false, false),
			field("nom_langue", "Nom", "text", false, false),
		}),
		"translations": resource("translations", "Traductions", "TRADUCTION", []string{"id"}, []AdminResourceField{
			field("id", "ID", "number", false, true),
			field("id_langue", "ID langue", "number", false, false),
			field("cle_traduction", "Cle", "text", false, false),
			field("text_traduit", "Texte", "textarea", false, false),
		}),
		"comments": resource("comments", "Commentaires", "COMMENTAIRE", []string{"id"}, []AdminResourceField{
			field("id", "ID", "number", false, true),
			field("id_utilisateur", "ID utilisateur", "number", true, false),
			field("description", "Description", "textarea", true, false),
		}),
		"reviews": resource("reviews", "Avis", "AVIS", []string{"id"}, []AdminResourceField{
			field("id", "ID", "number", false, true),
			field("id_auteur", "ID auteur", "number", true, false),
			field("id_cible", "ID cible", "number", true, false),
			field("note", "Note", "number", true, false),
			field("commentaire", "Commentaire", "textarea", false, false),
		}),
		"favorites": resource("favorites", "Favoris annonces", "FAVORIS", []string{"id_utilisateur", "id_annonce"}, []AdminResourceField{
			field("id_utilisateur", "ID utilisateur", "number", true, false),
			field("id_annonce", "ID annonce", "number", true, false),
		}),
		"formation_registrations": resource("formation_registrations", "Inscriptions formations", "FORMATION_INSCRIPTION", []string{"id_utilisateur", "id_formation"}, []AdminResourceField{
			field("id_utilisateur", "ID utilisateur", "number", true, false),
			field("id_formation", "ID formation", "number", true, false),
		}),
		"event_registrations": resource("event_registrations", "Inscriptions evenements", "EVENEMENT_INSCRIPTION", []string{"id_utilisateur", "id_evenement"}, []AdminResourceField{
			field("id_utilisateur", "ID utilisateur", "number", true, false),
			field("id_evenement", "ID evenement", "number", true, false),
		}),
		"project_registrations": resource("project_registrations", "Inscriptions projets", "PROJET_INSCRIPTION", []string{"id_utilisateur", "id_projet"}, []AdminResourceField{
			field("id_utilisateur", "ID utilisateur", "number", true, false),
			field("id_projet", "ID projet", "number", true, false),
		}),
		"project_likes": resource("project_likes", "Likes projets", "PROJET_LIKE", []string{"id_utilisateur", "id_projet"}, []AdminResourceField{
			field("id_utilisateur", "ID utilisateur", "number", true, false),
			field("id_projet", "ID projet", "number", true, false),
		}),
		"project_views": resource("project_views", "Vues projets", "PROJET_VUE", []string{"id"}, []AdminResourceField{
			field("id", "ID", "number", false, true),
			field("id_projet", "ID projet", "number", true, false),
			field("id_utilisateur", "ID utilisateur", "number", false, false),
			field("ip_adresse", "IP", "text", true, false),
		}),
		"follows": resource("follows", "Abonnements sociaux", "ABONNEMENT_UTILISATEUR", []string{"id_abonne", "id_suivi"}, []AdminResourceField{
			field("id_abonne", "ID abonne", "number", true, false),
			field("id_suivi", "ID suivi", "number", true, false),
		}),
		"message_reports": resource("message_reports", "Signalements forum", "MESSAGE_SIGNALEMENT", []string{"id_message", "id_utilisateur"}, []AdminResourceField{
			field("id_message", "ID message", "number", true, false),
			field("id_utilisateur", "ID utilisateur", "number", true, false),
			field("motif", "Motif", "text", false, false),
		}),
		"order_lines": resource("order_lines", "Lignes de commande", "LIGNE_COMMANDE", []string{"id"}, []AdminResourceField{
			field("id", "ID", "number", false, true),
			field("id_commande", "ID commande", "number", true, false),
			field("id_vendeur", "ID vendeur", "number", false, false),
			enumField("type_item", "Type", []string{"Formation", "Annonce", "Abonnement", "Evenement"}),
			field("reference_id", "Reference", "number", true, false),
			field("prix_unitaire", "Prix", "decimal", true, false),
			field("commission_upc", "Commission", "decimal", false, false),
		}),
		"transactions": resource("transactions", "Transactions", "TRANSACTION", []string{"id"}, []AdminResourceField{
			field("id", "ID", "number", false, true),
			field("id_acheteur", "ID acheteur", "number", true, false),
			field("id_commande", "ID commande", "number", true, false),
			field("montant_total", "Montant total", "decimal", false, false),
			enumField("statut_paiement", "Statut paiement", []string{"En attente", "Valide", "Echoue"}),
			field("stripe_payment_id", "Stripe payment", "text", false, false),
		}),
		"dm_sales": resource("dm_sales", "Ventes DM", "DM_SALE", []string{"id"}, []AdminResourceField{
			field("id", "ID", "number", false, true),
			field("id_offer", "ID offre", "number", true, false),
			field("id_conversation", "ID conversation", "number", true, false),
			field("id_annonce", "ID annonce", "number", false, false),
			field("id_buyer", "ID acheteur", "number", true, false),
			field("id_seller", "ID vendeur", "number", true, false),
			field("amount", "Montant", "decimal", true, false),
			enumField("status", "Statut", []string{"Offre acceptee", "Payee", "Recue", "Evaluee"}),
			field("received_at", "Reception", "datetime", false, false),
			field("reviewed_at", "Evaluation", "datetime", false, false),
		}),
		"upcycling_scores": resource("upcycling_scores", "Scores upcycling", "UPCYCLING_SCORE", []string{"id"}, []AdminResourceField{
			field("id", "ID", "number", false, true),
			field("id_utilisateur", "ID utilisateur", "number", true, false),
			field("ressources_economisees", "Ressources", "decimal", false, false),
			field("co2_total_evite_kg", "CO2 total", "decimal", false, false),
			field("nb_objets_recycles", "Objets recycles", "number", false, false),
			field("total_points", "Points", "number", false, false),
			field("niveau", "Niveau", "text", false, false),
		}),
	}
}

func resource(key, label, table string, primaryKeys []string, fields []AdminResourceField) AdminResource {
	return AdminResource{Key: key, Label: label, Table: table, PrimaryKeys: primaryKeys, Fields: fields}
}

func field(name, label, fieldType string, required, readOnly bool) AdminResourceField {
	return AdminResourceField{Name: name, Label: label, Type: fieldType, Required: required, ReadOnly: readOnly}
}

func enumField(name, label string, options []string) AdminResourceField {
	return AdminResourceField{Name: name, Label: label, Type: "select", Options: options}
}

func AdminListResources(w http.ResponseWriter, r *http.Request) {
	defs := adminResourceDefinitions()
	items := make([]AdminResource, 0, len(defs))
	for _, item := range defs {
		items = append(items, item)
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func AdminListResourceRows(w http.ResponseWriter, r *http.Request) {
	def, ok := adminResourceDefinitions()[r.PathValue("resource")]
	if !ok {
		writeError(w, http.StatusNotFound, "resource_not_found")
		return
	}
	rows, err := listAdminResourceRows(def)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "cannot_read_resource")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"resource": def, "items": rows})
}

func AdminCreateResourceRow(w http.ResponseWriter, r *http.Request) {
	def, ok := adminResourceDefinitions()[r.PathValue("resource")]
	if !ok {
		writeError(w, http.StatusNotFound, "resource_not_found")
		return
	}
	var payload map[string]any
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}
	item, err := createAdminResourceRow(def, payload)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"created": item})
}

func AdminUpdateResourceRow(w http.ResponseWriter, r *http.Request) {
	def, ok := adminResourceDefinitions()[r.PathValue("resource")]
	if !ok {
		writeError(w, http.StatusNotFound, "resource_not_found")
		return
	}
	var payload map[string]any
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}
	item, err := updateAdminResourceRow(def, r.PathValue("key"), payload)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"updated": item})
}

func AdminDeleteResourceRow(w http.ResponseWriter, r *http.Request) {
	def, ok := adminResourceDefinitions()[r.PathValue("resource")]
	if !ok {
		writeError(w, http.StatusNotFound, "resource_not_found")
		return
	}
	if err := deleteAdminResourceRow(def, r.PathValue("key")); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"deleted": r.PathValue("key")})
}

func listAdminResourceRows(def AdminResource) ([]map[string]any, error) {
	columns := writableAndReadableColumns(def)
	query := fmt.Sprintf("SELECT %s FROM %s ORDER BY %s DESC LIMIT 500", quoteList(columns), quoteIdent(def.Table), quoteIdent(def.PrimaryKeys[0]))
	rows, err := db.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanResourceRows(def, columns, rows)
}

func createAdminResourceRow(def AdminResource, payload map[string]any) (map[string]any, error) {
	columns, values, err := mutationColumns(def, payload, false)
	if err != nil {
		return nil, err
	}
	placeholders := make([]string, len(columns))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", quoteIdent(def.Table), quoteList(columns), strings.Join(placeholders, ", "))
	res, err := db.Conn.Exec(query, values...)
	if err != nil {
		return nil, err
	}
	if len(def.PrimaryKeys) == 1 && hasReadOnlyID(def) {
		id, _ := res.LastInsertId()
		return getAdminResourceRow(def, encodeResourceKey([]string{strconv.FormatInt(id, 10)}))
	}
	return map[string]any{"ok": true}, nil
}

func updateAdminResourceRow(def AdminResource, encodedKey string, payload map[string]any) (map[string]any, error) {
	columns, values, err := mutationColumns(def, payload, true)
	if err != nil {
		return nil, err
	}
	assignments := make([]string, 0, len(columns))
	for _, column := range columns {
		assignments = append(assignments, fmt.Sprintf("%s = ?", quoteIdent(column)))
	}
	keyValues, err := decodeResourceKey(def, encodedKey)
	if err != nil {
		return nil, err
	}
	where, keyArgs := resourceWhere(def, keyValues)
	values = append(values, keyArgs...)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", quoteIdent(def.Table), strings.Join(assignments, ", "), where)
	if _, err := db.Conn.Exec(query, values...); err != nil {
		return nil, err
	}
	return getAdminResourceRow(def, encodedKey)
}

func deleteAdminResourceRow(def AdminResource, encodedKey string) error {
	keyValues, err := decodeResourceKey(def, encodedKey)
	if err != nil {
		return err
	}
	where, args := resourceWhere(def, keyValues)
	_, err = db.Conn.Exec(fmt.Sprintf("DELETE FROM %s WHERE %s", quoteIdent(def.Table), where), args...)
	return err
}

func getAdminResourceRow(def AdminResource, encodedKey string) (map[string]any, error) {
	keyValues, err := decodeResourceKey(def, encodedKey)
	if err != nil {
		return nil, err
	}
	columns := writableAndReadableColumns(def)
	where, args := resourceWhere(def, keyValues)
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s", quoteList(columns), quoteIdent(def.Table), where)
	rows, err := db.Conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items, err := scanResourceRows(def, columns, rows)
	if err != nil || len(items) == 0 {
		return nil, err
	}
	return items[0], nil
}

func scanResourceRows(def AdminResource, columns []string, rows *sql.Rows) ([]map[string]any, error) {
	items := []map[string]any{}
	for rows.Next() {
		raw := make([]sql.NullString, len(columns))
		dest := make([]any, len(columns))
		for i := range raw {
			dest[i] = &raw[i]
		}
		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}
		item := map[string]any{}
		keyParts := make([]string, len(def.PrimaryKeys))
		for i, column := range columns {
			value := ""
			if raw[i].Valid {
				value = raw[i].String
			}
			item[column] = value
			for k, keyColumn := range def.PrimaryKeys {
				if column == keyColumn {
					keyParts[k] = value
				}
			}
		}
		item["__key"] = encodeResourceKey(keyParts)
		items = append(items, item)
	}
	return items, rows.Err()
}

func mutationColumns(def AdminResource, payload map[string]any, update bool) ([]string, []any, error) {
	columns := []string{}
	values := []any{}
	for _, f := range def.Fields {
		if f.ReadOnly {
			continue
		}
		value, exists := payload[f.Name]
		if !exists {
			if f.Required && !update {
				return nil, nil, fmt.Errorf("%s_required", f.Name)
			}
			continue
		}
		normalized := normalizeResourceValue(f, value)
		if f.Required && normalized == nil {
			return nil, nil, fmt.Errorf("%s_required", f.Name)
		}
		columns = append(columns, f.Name)
		values = append(values, normalized)
	}
	if len(columns) == 0 {
		return nil, nil, fmt.Errorf("empty_payload")
	}
	return columns, values, nil
}

func normalizeResourceValue(field AdminResourceField, value any) any {
	if value == nil {
		return nil
	}
	text := strings.TrimSpace(fmt.Sprint(value))
	if text == "" {
		return nil
	}
	if field.Type == "boolean" {
		if text == "true" || text == "1" || strings.EqualFold(text, "on") {
			return 1
		}
		return 0
	}
	return text
}

func writableAndReadableColumns(def AdminResource) []string {
	columns := make([]string, 0, len(def.Fields))
	for _, f := range def.Fields {
		columns = append(columns, f.Name)
	}
	return columns
}

func hasReadOnlyID(def AdminResource) bool {
	for _, f := range def.Fields {
		if f.Name == def.PrimaryKeys[0] && f.ReadOnly {
			return true
		}
	}
	return false
}

func resourceWhere(def AdminResource, values []string) (string, []any) {
	parts := make([]string, len(def.PrimaryKeys))
	args := make([]any, len(def.PrimaryKeys))
	for i, column := range def.PrimaryKeys {
		parts[i] = fmt.Sprintf("%s = ?", quoteIdent(column))
		args[i] = values[i]
	}
	return strings.Join(parts, " AND "), args
}

func encodeResourceKey(parts []string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(strings.Join(parts, "\x1f")))
}

func decodeResourceKey(def AdminResource, key string) ([]string, error) {
	raw, err := base64.RawURLEncoding.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("invalid_key")
	}
	parts := strings.Split(string(raw), "\x1f")
	if len(parts) != len(def.PrimaryKeys) {
		return nil, fmt.Errorf("invalid_key")
	}
	return parts, nil
}

func quoteIdent(identifier string) string {
	return "`" + strings.ReplaceAll(identifier, "`", "``") + "`"
}

func quoteList(columns []string) string {
	quoted := make([]string, len(columns))
	for i, column := range columns {
		quoted[i] = quoteIdent(column)
	}
	return strings.Join(quoted, ", ")
}
