package app

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"upcycleconnect/api-go/db"
)

func adminDBEnabled() bool {
	return db.Conn != nil
}

func parseAdminIntID(id string) (int, error) {
	value, err := strconv.Atoi(strings.TrimSpace(id))
	if err != nil || value <= 0 {
		return 0, errors.New("invalid_id")
	}
	return value, nil
}

func userStatusFromDB(value string) string {
	switch strings.TrimSpace(value) {
	case "Actif":
		return "active"
	case "Inactif":
		return "inactive"
	case "Banni":
		return "archived"
	default:
		return "inactive"
	}
}

func userStatusToDB(value string) string {
	switch strings.TrimSpace(value) {
	case "active":
		return "Actif"
	case "archived":
		return "Banni"
	default:
		return "Inactif"
	}
}

func categoryStatusFromDB(value string) string {
	if strings.EqualFold(strings.TrimSpace(value), "active") {
		return "active"
	}
	return "archived"
}

func categoryStatusToDB(value string) string {
	if strings.EqualFold(strings.TrimSpace(value), "active") {
		return "active"
	}
	return "inactive"
}

func annonceTypeFromDB(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "don":
		return "don"
	case "service":
		return "service"
	default:
		return "vente"
	}
}

func annonceTypeToDB(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "don":
		return "Don"
	case "service":
		return "Service"
	default:
		return "Vente"
	}
}

func annonceStatusFromDB(statut, estValide string) string {
	if strings.EqualFold(strings.TrimSpace(statut), "Annule") {
		return "archived"
	}
	if strings.EqualFold(strings.TrimSpace(estValide), "Valide") {
		return "published"
	}
	return "draft"
}

func annonceStatusToDB(value string) (string, string) {
	switch strings.TrimSpace(value) {
	case "published":
		return "Disponible", "Valide"
	case "archived":
		return "Annule", "Refuse"
	default:
		return "Disponible", "En attente"
	}
}

func normalizeDateTime(value string) string {
	if strings.TrimSpace(value) == "" {
		return ""
	}

	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed.Format("2006-01-02 15:04:05")
		}
	}
	return value
}

func createdDate(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format("2006-01-02")
}

func dateOnly(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format("2006-01-02")
}

func activeAdminCountDB(excludedID int) (int, error) {
	var count int
	err := db.Conn.QueryRow(`
		SELECT COUNT(*)
		FROM UTILISATEUR
		WHERE role = 'Admin' AND statut = 'Actif' AND id <> ?
	`, excludedID).Scan(&count)
	return count, err
}

func listAdminUsersFromDB() ([]AdminUser, error) {
	rows, err := db.Conn.Query(`
		SELECT id, prenom, nom, mail, ville, code_postal, role, statut, date_inscription
		FROM UTILISATEUR
		ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []AdminUser{}
	for rows.Next() {
		var item AdminUser
		var id int
		var createdAt time.Time
		var statut string
		if err := rows.Scan(&id, &item.FirstName, &item.LastName, &item.Email, &item.City, &item.PostalCode, &item.Role, &statut, &createdAt); err != nil {
			return nil, err
		}
		item.ID = strconv.Itoa(id)
		item.FullName = fullName(item.FirstName, item.LastName)
		item.Status = userStatusFromDB(statut)
		item.CreatedAt = createdDate(createdAt)
		items = append(items, item)
	}
	return items, rows.Err()
}

func getAdminUserFromDB(id string) (*AdminUser, error) {
	intID, err := parseAdminIntID(id)
	if err != nil {
		return nil, err
	}

	var item AdminUser
	var intIDValue int
	var createdAt time.Time
	var statut string
	err = db.Conn.QueryRow(`
		SELECT id, prenom, nom, mail, ville, code_postal, role, statut, date_inscription
		FROM UTILISATEUR
		WHERE id = ?
	`, intID).Scan(&intIDValue, &item.FirstName, &item.LastName, &item.Email, &item.City, &item.PostalCode, &item.Role, &statut, &createdAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	item.ID = strconv.Itoa(intIDValue)
	if err != nil {
		return nil, err
	}
	item.FullName = fullName(item.FirstName, item.LastName)
	item.Status = userStatusFromDB(statut)
	item.CreatedAt = createdDate(createdAt)
	return &item, nil
}

func createAdminUserInDB(payload AdminUser) (*AdminUser, error) {
	payload = normalizeUser(payload)
	if payload.Email == "" {
		return nil, errors.New("email must be unique")
	}

	var existing int
	if err := db.Conn.QueryRow(`SELECT COUNT(*) FROM UTILISATEUR WHERE mail = ?`, payload.Email).Scan(&existing); err != nil {
		return nil, err
	}
	if existing > 0 {
		return nil, errors.New("email must be unique")
	}

	tempPassword := "AdminTemp1!"
	res, err := db.Conn.Exec(`
		INSERT INTO UTILISATEUR (prenom, nom, password, mail, adresse, ville, code_postal, date_naissance, role, statut, id_langue, date_inscription)
		VALUES (?, ?, ?, ?, '', ?, ?, NULL, ?, ?, 1, NOW())
	`, payload.FirstName, payload.LastName, tempPassword, payload.Email, payload.City, payload.PostalCode, payload.Role, userStatusToDB(payload.Status))
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return getAdminUserFromDB(strconv.FormatInt(id, 10))
}

func updateAdminUserInDB(id string, payload AdminUser) (*AdminUser, error) {
	intID, err := parseAdminIntID(id)
	if err != nil {
		return nil, err
	}
	current, err := getAdminUserFromDB(id)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, nil
	}

	payload = normalizeUser(payload)
	var existing int
	if err := db.Conn.QueryRow(`SELECT COUNT(*) FROM UTILISATEUR WHERE mail = ? AND id <> ?`, payload.Email, intID).Scan(&existing); err != nil {
		return nil, err
	}
	if existing > 0 {
		return nil, errors.New("email must be unique")
	}

	if current.Role == "Admin" && current.Status == "active" && (payload.Role != "Admin" || payload.Status != "active") {
		count, err := activeAdminCountDB(intID)
		if err != nil {
			return nil, err
		}
		if count == 0 {
			return nil, errors.New("at least one active admin must remain available")
		}
	}

	_, err = db.Conn.Exec(`
		UPDATE UTILISATEUR
		SET prenom = ?, nom = ?, mail = ?, ville = ?, code_postal = ?, role = ?, statut = ?
		WHERE id = ?
	`, payload.FirstName, payload.LastName, payload.Email, payload.City, payload.PostalCode, payload.Role, userStatusToDB(payload.Status), intID)
	if err != nil {
		return nil, err
	}
	return getAdminUserFromDB(id)
}

func toggleAdminUserStatusInDB(id string) (*AdminUser, error) {
	current, err := getAdminUserFromDB(id)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, nil
	}
	nextStatus := "active"
	if current.Status == "active" {
		intID, _ := parseAdminIntID(id)
		if current.Role == "Admin" {
			count, err := activeAdminCountDB(intID)
			if err != nil {
				return nil, err
			}
			if count == 0 {
				return nil, errors.New("the last active admin cannot be deactivated")
			}
		}
		nextStatus = "inactive"
	}
	_, err = db.Conn.Exec(`UPDATE UTILISATEUR SET statut = ? WHERE id = ?`, userStatusToDB(nextStatus), id)
	if err != nil {
		return nil, err
	}
	return getAdminUserFromDB(id)
}

func deleteAdminUserInDB(id string) error {
	current, err := getAdminUserFromDB(id)
	if err != nil {
		return err
	}
	if current == nil {
		return nil
	}
	if current.Role == "Admin" && current.Status == "active" {
		intID, _ := parseAdminIntID(id)
		count, err := activeAdminCountDB(intID)
		if err != nil {
			return err
		}
		if count == 0 {
			return errors.New("the last active admin cannot be deleted")
		}
	}
	_, err = db.Conn.Exec(`DELETE FROM UTILISATEUR WHERE id = ?`, id)
	return err
}

func listAdminPrestationsFromDB(filterType string) ([]AdminPrestation, error) {
	query := `
		SELECT id, titre, description, type, prix, statut, est_valide, provider, date_creation
		FROM ANNONCE
	`
	args := []any{}
	if filterType != "" {
		query += ` WHERE type = ?`
		args = append(args, annonceTypeToDB(filterType))
	}
	query += ` ORDER BY id DESC`

	rows, err := db.Conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []AdminPrestation{}
	for rows.Next() {
		var item AdminPrestation
		var id int
		var dbType, statut, estValide string
		var createdAt time.Time
		if err := rows.Scan(&id, &item.Title, &item.Description, &dbType, &item.Price, &statut, &estValide, &item.Provider, &createdAt); err != nil {
			return nil, err
		}
		item.ID = strconv.Itoa(id)
		item.Type = annonceTypeFromDB(dbType)
		item.Status = annonceStatusFromDB(statut, estValide)
		item.Provider = defaultString(item.Provider, "Non assigne")
		item.CreatedAt = createdDate(createdAt)
		items = append(items, item)
	}
	return items, rows.Err()
}

func getAdminPrestationFromDB(id string) (*AdminPrestation, error) {
	var item AdminPrestation
	var intIDValue int
	var dbType, statut, estValide string
	var createdAt time.Time
	err := db.Conn.QueryRow(`
		SELECT id, titre, description, type, prix, statut, est_valide, provider, date_creation
		FROM ANNONCE WHERE id = ?
	`, id).Scan(&intIDValue, &item.Title, &item.Description, &dbType, &item.Price, &statut, &estValide, &item.Provider, &createdAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	item.ID = strconv.Itoa(intIDValue)
	if err != nil {
		return nil, err
	}
	item.Type = annonceTypeFromDB(dbType)
	item.Status = annonceStatusFromDB(statut, estValide)
	item.Provider = defaultString(item.Provider, "Non assigne")
	item.CreatedAt = createdDate(createdAt)
	return &item, nil
}

func createAdminPrestationInDB(payload AdminPrestation) (*AdminPrestation, error) {
	payload = normalizePrestation(payload)
	statut, estValide := annonceStatusToDB(payload.Status)
	res, err := db.Conn.Exec(`
		INSERT INTO ANNONCE (id_vendeur, id_acheteur, titre, description, statut, est_valide, prix, etat_objet, adresse, ville, code_postal, provider, type, date_creation)
		VALUES (1, NULL, ?, ?, ?, ?, ?, 'Bon etat', '', '', '', ?, ?, NOW())
	`, payload.Title, payload.Description, statut, estValide, payload.Price, payload.Provider, annonceTypeToDB(payload.Type))
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return getAdminPrestationFromDB(strconv.FormatInt(id, 10))
}

func updateAdminPrestationInDB(id string, payload AdminPrestation) (*AdminPrestation, error) {
	payload = normalizePrestation(payload)
	statut, estValide := annonceStatusToDB(payload.Status)
	_, err := db.Conn.Exec(`
		UPDATE ANNONCE
		SET titre = ?, description = ?, prix = ?, provider = ?, type = ?, statut = ?, est_valide = ?
		WHERE id = ?
	`, payload.Title, payload.Description, payload.Price, payload.Provider, annonceTypeToDB(payload.Type), statut, estValide, id)
	if err != nil {
		return nil, err
	}
	return getAdminPrestationFromDB(id)
}

func deleteAdminPrestationInDB(id string) error {
	_, err := db.Conn.Exec(`DELETE FROM ANNONCE WHERE id = ?`, id)
	return err
}

func listAdminCategoriesFromDB() ([]AdminCategory, error) {
	rows, err := db.Conn.Query(`SELECT id, nom, id_parent, description, statut FROM CATEGORIE ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []AdminCategory{}
	for rows.Next() {
		var item AdminCategory
		var id int
		var parentID sql.NullInt64
		var statut string
		if err := rows.Scan(&id, &item.Name, &parentID, &item.Description, &statut); err != nil {
			return nil, err
		}
		item.ID = strconv.Itoa(id)
		if parentID.Valid {
			item.ParentID = strconv.FormatInt(parentID.Int64, 10)
		}
		item.Status = categoryStatusFromDB(statut)
		items = append(items, item)
	}
	return items, rows.Err()
}

func getAdminCategoryFromDB(id string) (*AdminCategory, error) {
	var item AdminCategory
	var intIDValue int
	var parentID sql.NullInt64
	var statut string
	err := db.Conn.QueryRow(`SELECT id, nom, id_parent, description, statut FROM CATEGORIE WHERE id = ?`, id).Scan(&intIDValue, &item.Name, &parentID, &item.Description, &statut)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	item.ID = strconv.Itoa(intIDValue)
	if err != nil {
		return nil, err
	}
	if parentID.Valid {
		item.ParentID = strconv.FormatInt(parentID.Int64, 10)
	}
	item.Status = categoryStatusFromDB(statut)
	return &item, nil
}

func createAdminCategoryInDB(payload AdminCategory) (*AdminCategory, error) {
	payload = normalizeCategory(payload)
	var parent any
	if payload.ParentID != "" {
		parentID, err := parseAdminIntID(payload.ParentID)
		if err != nil {
			return nil, err
		}
		parent = parentID
	}
	res, err := db.Conn.Exec(`
		INSERT INTO CATEGORIE (nom, description, id_parent, statut, date_creation)
		VALUES (?, ?, ?, ?, NOW())
	`, payload.Name, payload.Description, parent, categoryStatusToDB(payload.Status))
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return getAdminCategoryFromDB(strconv.FormatInt(id, 10))
}

func updateAdminCategoryInDB(id string, payload AdminCategory) (*AdminCategory, error) {
	payload = normalizeCategory(payload)
	var parent any
	if payload.ParentID != "" {
		parentID, err := parseAdminIntID(payload.ParentID)
		if err != nil {
			return nil, err
		}
		parent = parentID
	}
	_, err := db.Conn.Exec(`
		UPDATE CATEGORIE
		SET nom = ?, description = ?, id_parent = ?, statut = ?
		WHERE id = ?
	`, payload.Name, payload.Description, parent, categoryStatusToDB(payload.Status), id)
	if err != nil {
		return nil, err
	}
	return getAdminCategoryFromDB(id)
}

func deleteAdminCategoryInDB(id string) error {
	var count int
	if err := db.Conn.QueryRow(`SELECT COUNT(*) FROM CATEGORIE WHERE id_parent = ?`, id).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return errors.New("category has children and cannot be deleted")
	}
	_, err := db.Conn.Exec(`DELETE FROM CATEGORIE WHERE id = ?`, id)
	return err
}

func listAdminEventsFromDB() ([]AdminEvent, error) {
	rows, err := db.Conn.Query(`
		SELECT id, titre, CONCAT_WS(', ', adresse, ville, code_postal) AS location, date_evenement, statut, capacite_max, description
		FROM EVENEMENT
		ORDER BY date_evenement ASC, id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []AdminEvent{}
	for rows.Next() {
		var item AdminEvent
		var id int
		var eventDate time.Time
		if err := rows.Scan(&id, &item.Title, &item.Location, &eventDate, &item.Status, &item.Capacity, &item.Description); err != nil {
			return nil, err
		}
		item.ID = strconv.Itoa(id)
		item.Date = dateOnly(eventDate)
		items = append(items, item)
	}
	return items, rows.Err()
}

func getAdminEventFromDB(id string) (*AdminEvent, error) {
	var item AdminEvent
	var intIDValue int
	var eventDate time.Time
	err := db.Conn.QueryRow(`
		SELECT id, titre, CONCAT_WS(', ', adresse, ville, code_postal) AS location, date_evenement, statut, capacite_max, description
		FROM EVENEMENT
		WHERE id = ?
	`, id).Scan(&intIDValue, &item.Title, &item.Location, &eventDate, &item.Status, &item.Capacity, &item.Description)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	item.ID = strconv.Itoa(intIDValue)
	if err != nil {
		return nil, err
	}
	item.Date = dateOnly(eventDate)
	return &item, nil
}

func splitLocation(location string) (string, string, string) {
	parts := strings.Split(location, ",")
	address, city, postal := "", "", ""
	if len(parts) > 0 {
		address = strings.TrimSpace(parts[0])
	}
	if len(parts) > 1 {
		city = strings.TrimSpace(parts[1])
	}
	if len(parts) > 2 {
		postal = strings.TrimSpace(parts[2])
	}
	return address, city, postal
}

func createAdminEventInDB(payload AdminEvent) (*AdminEvent, error) {
	payload = normalizeEvent(payload)
	address, city, postal := splitLocation(payload.Location)
	res, err := db.Conn.Exec(`
		INSERT INTO EVENEMENT (titre, description, adresse, ville, code_postal, date_evenement, capacite_max, statut, type)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'Atelier')
	`, payload.Title, payload.Description, address, city, postal, normalizeDateTime(payload.Date), payload.Capacity, payload.Status)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return getAdminEventFromDB(strconv.FormatInt(id, 10))
}

func updateAdminEventInDB(id string, payload AdminEvent) (*AdminEvent, error) {
	payload = normalizeEvent(payload)
	address, city, postal := splitLocation(payload.Location)
	_, err := db.Conn.Exec(`
		UPDATE EVENEMENT
		SET titre = ?, description = ?, adresse = ?, ville = ?, code_postal = ?, date_evenement = ?, capacite_max = ?, statut = ?
		WHERE id = ?
	`, payload.Title, payload.Description, address, city, postal, normalizeDateTime(payload.Date), payload.Capacity, payload.Status, id)
	if err != nil {
		return nil, err
	}
	return getAdminEventFromDB(id)
}

func deleteAdminEventInDB(id string) error {
	_, err := db.Conn.Exec(`DELETE FROM EVENEMENT WHERE id = ?`, id)
	return err
}

func financeOverviewFromDB() (map[string]any, []AdminFinanceRecord, error) {
	items := []AdminFinanceRecord{}

	rows, err := db.Conn.Query(`
		SELECT id, type, montant_ht, statut_paiement, date_transaction
		FROM ` + "`TRANSACTION`" + `
		ORDER BY date_transaction DESC, id DESC
	`)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	paidTotal := 0.0
	pendingTotal := 0.0
	for rows.Next() {
		var id int
		var category, status string
		var amount float64
		var dueDate time.Time
		if err := rows.Scan(&id, &category, &amount, &status, &dueDate); err != nil {
			return nil, nil, err
		}
		recordStatus := "failed"
		switch status {
		case "Valide":
			recordStatus = "paid"
			paidTotal += amount
		case "En attente":
			recordStatus = "pending"
			pendingTotal += amount
		}
		items = append(items, AdminFinanceRecord{
			ID:       strconv.Itoa(id),
			Label:    fmt.Sprintf("Transaction #%d", id),
			Category: strings.ToLower(category),
			Amount:   amount,
			Status:   recordStatus,
			DueDate:  createdDate(dueDate),
			Source:   "transaction",
		})
	}

	var activeSubscriptions int
	if err := db.Conn.QueryRow(`SELECT COUNT(*) FROM ABONNEMENT WHERE statut = 'Actif'`).Scan(&activeSubscriptions); err != nil {
		return nil, nil, err
	}

	return map[string]any{
		"paidTotal":           paidTotal,
		"pendingTotal":        pendingTotal,
		"activeSubscriptions": activeSubscriptions,
	}, items, nil
}

func listAdminNotificationsFromDB() ([]AdminNotification, error) {
	rows, err := db.Conn.Query(`
		SELECT id, titre, canal, audience, statut, scheduled_at, message
		FROM NOTIFICATION
		ORDER BY date_envoi DESC, id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []AdminNotification{}
	for rows.Next() {
		var item AdminNotification
		var id int
		var scheduledAt sql.NullTime
		if err := rows.Scan(&id, &item.Title, &item.Channel, &item.Audience, &item.Status, &scheduledAt, &item.Message); err != nil {
			return nil, err
		}
		item.ID = strconv.Itoa(id)
		if scheduledAt.Valid {
			item.ScheduledAt = scheduledAt.Time.Format("2006-01-02T15:04")
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func createAdminNotificationInDB(payload AdminNotification) (*AdminNotification, error) {
	payload = normalizeNotification(payload)
	var userID int
	if err := db.Conn.QueryRow(`SELECT id FROM UTILISATEUR ORDER BY id ASC LIMIT 1`).Scan(&userID); err != nil {
		return nil, errors.New("notification requires at least one user in database")
	}
	var scheduled any
	if payload.ScheduledAt != "" {
		scheduled = normalizeDateTime(payload.ScheduledAt)
	}
	res, err := db.Conn.Exec(`
		INSERT INTO NOTIFICATION (id_utilisateur, type, canal, audience, statut, titre, message, scheduled_at, lu, date_envoi)
		VALUES (?, 'Message', ?, ?, ?, ?, ?, ?, 0, NOW())
	`, userID, payload.Channel, payload.Audience, payload.Status, payload.Title, payload.Message, scheduled)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	items, err := listAdminNotificationsFromDB()
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.ID == strconv.FormatInt(id, 10) {
			copyItem := item
			return &copyItem, nil
		}
	}
	return nil, nil
}

func updateAdminNotificationStatusInDB(id, status string) (*AdminNotification, error) {
	_, err := db.Conn.Exec(`UPDATE NOTIFICATION SET statut = ? WHERE id = ?`, status, id)
	if err != nil {
		return nil, err
	}
	items, err := listAdminNotificationsFromDB()
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.ID == id {
			copyItem := item
			return &copyItem, nil
		}
	}
	return nil, nil
}

func deleteAdminNotificationInDB(id string) error {
	_, err := db.Conn.Exec(`DELETE FROM NOTIFICATION WHERE id = ?`, id)
	return err
}

func moderationQueueFromDB() ([]map[string]string, error) {
	items := []map[string]string{}

	rows, err := db.Conn.Query(`
		SELECT id, titre, description, provider
		FROM ANNONCE
		WHERE est_valide = 'En attente'
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var title, description, owner sql.NullString
		if err := rows.Scan(&id, &title, &description, &owner); err != nil {
			return nil, err
		}
		items = append(items, map[string]string{
			"id":          strconv.Itoa(id),
			"type":        "prestations",
			"title":       title.String,
			"description": description.String,
			"owner":       defaultString(owner.String, "Catalogue"),
			"status":      "draft",
		})
	}

	eventRows, err := db.Conn.Query(`
		SELECT id, titre, description, CONCAT_WS(', ', adresse, ville)
		FROM EVENEMENT
		WHERE statut = 'planned'
	`)
	if err != nil {
		return nil, err
	}
	defer eventRows.Close()
	for eventRows.Next() {
		var id int
		var title, description, owner sql.NullString
		if err := eventRows.Scan(&id, &title, &description, &owner); err != nil {
			return nil, err
		}
		items = append(items, map[string]string{
			"id":          strconv.Itoa(id),
			"type":        "events",
			"title":       title.String,
			"description": description.String,
			"owner":       owner.String,
			"status":      "planned",
		})
	}

	return items, nil
}

func publishAdminPrestationInDB(id string) (*AdminPrestation, error) {
	_, err := db.Conn.Exec(`UPDATE ANNONCE SET est_valide = 'Valide', statut = 'Disponible' WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return getAdminPrestationFromDB(id)
}

func archiveAdminPrestationInDB(id string) (*AdminPrestation, error) {
	_, err := db.Conn.Exec(`UPDATE ANNONCE SET statut = 'Annule', est_valide = 'Refuse' WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return getAdminPrestationFromDB(id)
}

func publishAdminEventInDB(id string) (*AdminEvent, error) {
	_, err := db.Conn.Exec(`UPDATE EVENEMENT SET statut = 'published' WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return getAdminEventFromDB(id)
}

func archiveAdminEventInDB(id string) (*AdminEvent, error) {
	_, err := db.Conn.Exec(`UPDATE EVENEMENT SET statut = 'archived' WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return getAdminEventFromDB(id)
}
