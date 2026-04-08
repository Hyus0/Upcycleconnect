package app

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type AdminUser struct {
	ID         string `json:"id"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	FullName   string `json:"fullName,omitempty"`
	Email      string `json:"email"`
	City       string `json:"city"`
	PostalCode string `json:"postalCode"`
	Role       string `json:"role"`
	Status     string `json:"status"`
	CreatedAt  string `json:"createdAt"`
}

type AdminPrestation struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Type        string  `json:"type"`
	Price       float64 `json:"price"`
	Status      string  `json:"status"`
	Provider    string  `json:"provider"`
	CreatedAt   string  `json:"createdAt"`
}

type AdminCategory struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ParentID    string `json:"parentId"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type AdminEvent struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Location    string `json:"location"`
	Date        string `json:"date"`
	Status      string `json:"status"`
	Capacity    int    `json:"capacity"`
	Description string `json:"description"`
}

type AdminFinanceRecord struct {
	ID       string  `json:"id"`
	Label    string  `json:"label"`
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
	Status   string  `json:"status"`
	DueDate  string  `json:"dueDate"`
	Source   string  `json:"source"`
}

type AdminNotification struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Channel     string `json:"channel"`
	Audience    string `json:"audience"`
	Status      string `json:"status"`
	ScheduledAt string `json:"scheduledAt"`
	Message     string `json:"message"`
}

type adminStoreData struct {
	Users          []AdminUser       `json:"users"`
	Prestations    []AdminPrestation `json:"prestations"`
	Categories     []AdminCategory   `json:"categories"`
	Events         []AdminEvent      `json:"events"`
	FinanceRecords []AdminFinanceRecord `json:"financeRecords"`
	Notifications  []AdminNotification  `json:"notifications"`
	UserSeq        int               `json:"userSeq"`
	PrestationSeq  int               `json:"prestationSeq"`
	CategorySeq    int               `json:"categorySeq"`
	EventSeq       int               `json:"eventSeq"`
	FinanceSeq     int               `json:"financeSeq"`
	NotificationSeq int              `json:"notificationSeq"`
	LastUpdatedAt  string            `json:"lastUpdatedAt"`
}

type adminStore struct {
	mu       sync.RWMutex
	filePath string
	data     adminStoreData
}

var store = newAdminStore()

func newAdminStore() *adminStore {
	filePath := os.Getenv("ADMIN_STORE_PATH")
	if strings.TrimSpace(filePath) == "" {
		filePath = filepath.Join("storage", "admin_data.json")
	}

	s := &adminStore{
		filePath: filePath,
		data:     defaultAdminStoreData(),
	}

	if err := s.load(); err != nil {
		_ = s.saveLocked()
	}

	return s
}

func defaultAdminStoreData() adminStoreData {
	return adminStoreData{
		Users: []AdminUser{
			{ID: "u1", FirstName: "Alice", LastName: "Martin", Email: "alice@upcycleconnect.local", City: "Paris", PostalCode: "75011", Role: "Particulier", Status: "active", CreatedAt: "2026-03-01"},
			{ID: "u2", FirstName: "Karim", LastName: "Benali", Email: "karim@atelier.local", City: "Lyon", PostalCode: "69002", Role: "Prestataire", Status: "active", CreatedAt: "2026-02-18"},
			{ID: "u3", FirstName: "Nina", LastName: "Roux", Email: "nina.admin@upcycleconnect.local", City: "Lille", PostalCode: "59000", Role: "Admin", Status: "inactive", CreatedAt: "2026-01-14"},
		},
		Prestations: []AdminPrestation{
			{ID: "p1", Title: "Diagnostic mobilier", Description: "Analyse d'etat et estimation de remise en valeur.", Type: "service", Price: 45, Status: "published", Provider: "Atelier Renouveau", CreatedAt: "2026-03-04"},
			{ID: "p2", Title: "Lot textile recycle", Description: "Selection textile revalorisee pour ateliers creatifs.", Type: "vente", Price: 120, Status: "draft", Provider: "ReTex", CreatedAt: "2026-03-08"},
		},
		Categories: []AdminCategory{
			{ID: "c1", Name: "Mobilier", ParentID: "", Description: "Meubles et accessoires", Status: "active"},
			{ID: "c2", Name: "Chaises", ParentID: "c1", Description: "Assises et tabourets", Status: "active"},
		},
		Events: []AdminEvent{
			{ID: "e1", Title: "Atelier bois", Location: "Paris 11", Date: "2026-04-20", Status: "planned", Capacity: 18, Description: "Session pratique autour de la remise en etat du bois."},
			{ID: "e2", Title: "Collecte textile", Location: "Lyon 2", Date: "2026-04-26", Status: "published", Capacity: 40, Description: "Journee de collecte et tri textile."},
		},
		FinanceRecords: []AdminFinanceRecord{
			{ID: "f1", Label: "Abonnement Pro Avril", Category: "subscription", Amount: 89, Status: "paid", DueDate: "2026-04-02", Source: "Stripe"},
			{ID: "f2", Label: "Atelier bois - session 20/04", Category: "event", Amount: 240, Status: "pending", DueDate: "2026-04-20", Source: "Reservation"},
			{ID: "f3", Label: "Diagnostic mobilier", Category: "service", Amount: 45, Status: "paid", DueDate: "2026-04-05", Source: "Catalogue"},
		},
		Notifications: []AdminNotification{
			{ID: "n1", Title: "Collecte textile maintenue", Channel: "email", Audience: "all", Status: "scheduled", ScheduledAt: "2026-04-18T10:00", Message: "La collecte textile de Lyon est maintenue ce weekend."},
			{ID: "n2", Title: "Nouveaux ateliers disponibles", Channel: "push", Audience: "particuliers", Status: "draft", ScheduledAt: "", Message: "Deux nouveaux ateliers viennent d'etre ajoutes au catalogue."},
		},
		UserSeq:       3,
		PrestationSeq: 2,
		CategorySeq:   2,
		EventSeq:      2,
		FinanceSeq:    3,
		NotificationSeq: 2,
		LastUpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func (s *adminStore) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	raw, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}
	if len(raw) == 0 {
		return errors.New("empty admin store")
	}

	var payload adminStoreData
	if err := json.Unmarshal(raw, &payload); err != nil {
		return err
	}

	s.data = payload
	s.ensureDefaultsLocked()
	return nil
}

func (s *adminStore) saveLocked() error {
	s.ensureDefaultsLocked()
	s.data.LastUpdatedAt = time.Now().UTC().Format(time.RFC3339)

	if err := os.MkdirAll(filepath.Dir(s.filePath), 0o755); err != nil {
		return err
	}

	raw, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, raw, 0o644)
}

func (s *adminStore) ensureDefaultsLocked() {
	if s.data.UserSeq < len(s.data.Users) {
		s.data.UserSeq = len(s.data.Users)
	}
	if s.data.PrestationSeq < len(s.data.Prestations) {
		s.data.PrestationSeq = len(s.data.Prestations)
	}
	if s.data.CategorySeq < len(s.data.Categories) {
		s.data.CategorySeq = len(s.data.Categories)
	}
	if s.data.EventSeq < len(s.data.Events) {
		s.data.EventSeq = len(s.data.Events)
	}
	if s.data.FinanceSeq < len(s.data.FinanceRecords) {
		s.data.FinanceSeq = len(s.data.FinanceRecords)
	}
	if s.data.NotificationSeq < len(s.data.Notifications) {
		s.data.NotificationSeq = len(s.data.Notifications)
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func decodeJSON(r *http.Request, dest any) error {
	return json.NewDecoder(r.Body).Decode(dest)
}

func fullName(firstName, lastName string) string {
	return strings.TrimSpace(firstName + " " + lastName)
}

func nextID(prefix string, seq *int) string {
	*seq += 1
	return prefix + strconv.Itoa(*seq)
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func isAllowedValue(value string, allowed ...string) bool {
	for _, item := range allowed {
		if value == item {
			return true
		}
	}
	return false
}

func parseISODate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", strings.TrimSpace(value))
}

func todayDate() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

func normalizeUser(payload AdminUser) AdminUser {
	payload.FirstName = strings.TrimSpace(payload.FirstName)
	payload.LastName = strings.TrimSpace(payload.LastName)
	payload.Email = strings.TrimSpace(payload.Email)
	payload.City = strings.TrimSpace(payload.City)
	payload.PostalCode = strings.TrimSpace(payload.PostalCode)
	payload.Role = defaultString(payload.Role, "Particulier")
	payload.Status = defaultString(payload.Status, "active")
	payload.CreatedAt = defaultString(payload.CreatedAt, time.Now().Format("2006-01-02"))
	payload.FullName = fullName(payload.FirstName, payload.LastName)
	return payload
}

func normalizePrestation(payload AdminPrestation) AdminPrestation {
	payload.Title = strings.TrimSpace(payload.Title)
	payload.Description = strings.TrimSpace(payload.Description)
	payload.Type = defaultString(payload.Type, "service")
	payload.Status = defaultString(payload.Status, "draft")
	payload.Provider = defaultString(payload.Provider, "Equipe locale")
	payload.CreatedAt = defaultString(payload.CreatedAt, time.Now().Format("2006-01-02"))
	return payload
}

func normalizeCategory(payload AdminCategory) AdminCategory {
	payload.Name = strings.TrimSpace(payload.Name)
	payload.ParentID = strings.TrimSpace(payload.ParentID)
	payload.Description = strings.TrimSpace(payload.Description)
	payload.Status = defaultString(payload.Status, "active")
	return payload
}

func normalizeEvent(payload AdminEvent) AdminEvent {
	payload.Title = strings.TrimSpace(payload.Title)
	payload.Location = strings.TrimSpace(payload.Location)
	payload.Date = strings.TrimSpace(payload.Date)
	payload.Status = defaultString(payload.Status, "planned")
	payload.Description = strings.TrimSpace(payload.Description)
	if payload.Capacity < 0 {
		payload.Capacity = 0
	}
	return payload
}

func normalizeNotification(payload AdminNotification) AdminNotification {
	payload.Title = strings.TrimSpace(payload.Title)
	payload.Channel = defaultString(strings.TrimSpace(payload.Channel), "email")
	payload.Audience = defaultString(strings.TrimSpace(payload.Audience), "all")
	payload.Status = defaultString(strings.TrimSpace(payload.Status), "draft")
	payload.ScheduledAt = strings.TrimSpace(payload.ScheduledAt)
	payload.Message = strings.TrimSpace(payload.Message)
	return payload
}

func validateUser(payload AdminUser) []string {
	var issues []string
	if len(payload.FirstName) < 2 {
		issues = append(issues, "firstName must be at least 2 characters")
	}
	if len(payload.LastName) < 2 {
		issues = append(issues, "lastName must be at least 2 characters")
	}
	if !strings.Contains(payload.Email, "@") || !strings.Contains(payload.Email, ".") {
		issues = append(issues, "email is invalid")
	}
	if payload.PostalCode != "" && len(payload.PostalCode) < 4 {
		issues = append(issues, "postalCode is invalid")
	}
	if !isAllowedValue(payload.Role, "Particulier", "Prestataire", "Admin") {
		issues = append(issues, "role is invalid")
	}
	if !isAllowedValue(payload.Status, "active", "inactive", "archived") {
		issues = append(issues, "status is invalid")
	}
	return issues
}

func validatePrestation(payload AdminPrestation) []string {
	var issues []string
	if len(payload.Title) < 3 {
		issues = append(issues, "title must be at least 3 characters")
	}
	if len(payload.Description) < 10 {
		issues = append(issues, "description must be at least 10 characters")
	}
	if payload.Price < 0 {
		issues = append(issues, "price must be >= 0")
	}
	if !isAllowedValue(payload.Type, "service", "vente", "don") {
		issues = append(issues, "type is invalid")
	}
	if !isAllowedValue(payload.Status, "draft", "published", "archived") {
		issues = append(issues, "status is invalid")
	}
	if payload.Status == "published" && strings.TrimSpace(payload.Provider) == "" {
		issues = append(issues, "published prestations must have a provider")
	}
	if payload.Type == "don" && payload.Price != 0 {
		issues = append(issues, "don prestations must have a price of 0")
	}
	if (payload.Type == "service" || payload.Type == "vente") && payload.Status == "published" && payload.Price <= 0 {
		issues = append(issues, "published prestations must have a price greater than 0")
	}
	return issues
}

func validateCategory(payload AdminCategory) []string {
	var issues []string
	if len(payload.Name) < 2 {
		issues = append(issues, "name must be at least 2 characters")
	}
	if !isAllowedValue(payload.Status, "active", "archived") {
		issues = append(issues, "status is invalid")
	}
	return issues
}

func validateEvent(payload AdminEvent) []string {
	var issues []string
	if len(payload.Title) < 3 {
		issues = append(issues, "title must be at least 3 characters")
	}
	if strings.TrimSpace(payload.Date) == "" {
		issues = append(issues, "date is required")
	} else {
		eventDate, err := parseISODate(payload.Date)
		if err != nil {
			issues = append(issues, "date format is invalid")
		} else if (payload.Status == "planned" || payload.Status == "published") && eventDate.Before(todayDate()) {
			issues = append(issues, "planned or published events cannot be dated in the past")
		}
	}
	if payload.Capacity < 0 {
		issues = append(issues, "capacity must be >= 0")
	}
	if (payload.Status == "planned" || payload.Status == "published") && strings.TrimSpace(payload.Location) == "" {
		issues = append(issues, "planned or published events require a location")
	}
	if payload.Status == "published" && payload.Capacity <= 0 {
		issues = append(issues, "published events require a capacity greater than 0")
	}
	if !isAllowedValue(payload.Status, "planned", "published", "archived") {
		issues = append(issues, "status is invalid")
	}
	return issues
}

func validateNotification(payload AdminNotification) []string {
	var issues []string
	if len(payload.Title) < 3 {
		issues = append(issues, "title must be at least 3 characters")
	}
	if len(payload.Message) < 10 {
		issues = append(issues, "message must be at least 10 characters")
	}
	if !isAllowedValue(payload.Channel, "email", "push", "sms") {
		issues = append(issues, "channel is invalid")
	}
	if !isAllowedValue(payload.Audience, "all", "particuliers", "prestataires", "admins") {
		issues = append(issues, "audience is invalid")
	}
	if !isAllowedValue(payload.Status, "draft", "scheduled", "sent") {
		issues = append(issues, "status is invalid")
	}
	if payload.Status == "scheduled" && payload.ScheduledAt == "" {
		issues = append(issues, "scheduled notifications require a scheduledAt value")
	}
	return issues
}

func emailExistsLocked(email, excludedID string) bool {
	for _, item := range store.data.Users {
		if strings.EqualFold(strings.TrimSpace(item.Email), strings.TrimSpace(email)) && item.ID != excludedID {
			return true
		}
	}
	return false
}

func categoryExistsLocked(id string) bool {
	for _, item := range store.data.Categories {
		if item.ID == id {
			return true
		}
	}
	return false
}

func buildModerationQueueLocked() []map[string]string {
	items := make([]map[string]string, 0)
	for _, item := range store.data.Prestations {
		if item.Status != "draft" {
			continue
		}
		items = append(items, map[string]string{
			"id":          item.ID,
			"type":        "prestation",
			"title":       item.Title,
			"description": item.Description,
			"owner":       item.Provider,
			"status":      item.Status,
		})
	}
	for _, item := range store.data.Events {
		if item.Status != "planned" {
			continue
		}
		items = append(items, map[string]string{
			"id":          item.ID,
			"type":        "event",
			"title":       item.Title,
			"description": item.Description,
			"owner":       item.Location,
			"status":      item.Status,
		})
	}
	return items
}

func activeAdminCountLocked() int {
	count := 0
	for _, item := range store.data.Users {
		if item.Role == "Admin" && item.Status == "active" {
			count += 1
		}
	}
	return count
}

func writeValidationError(w http.ResponseWriter, issues []string) {
	writeJSON(w, http.StatusBadRequest, map[string]any{"error": "validation_failed", "issues": issues})
}

func AdminMetrics(w http.ResponseWriter, r *http.Request) {
	if adminDBEnabled() {
		users, err := listAdminUsersFromDB()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "cannot_read_users")
			return
		}
		prestations, err := listAdminPrestationsFromDB("")
		if err != nil {
			writeError(w, http.StatusInternalServerError, "cannot_read_prestations")
			return
		}
		categories, err := listAdminCategoriesFromDB()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "cannot_read_categories")
			return
		}
		events, err := listAdminEventsFromDB()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "cannot_read_events")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"source": "database",
			"metrics": map[string]int{
				"users":      len(users),
				"annonces":   len(prestations),
				"categories": len(categories),
				"events":     len(events),
			},
		})
		return
	}

	store.mu.RLock()
	defer store.mu.RUnlock()

	writeJSON(w, http.StatusOK, map[string]any{
		"source": "api",
		"metrics": map[string]int{
			"users":      len(store.data.Users),
			"annonces":   len(store.data.Prestations),
			"categories": len(store.data.Categories),
			"events":     len(store.data.Events),
		},
		"updatedAt": store.data.LastUpdatedAt,
	})
}

func AdminListUsers(w http.ResponseWriter, r *http.Request) {
	if adminDBEnabled() {
		items, err := listAdminUsersFromDB()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "cannot_read_users")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": items})
		return
	}

	store.mu.RLock()
	defer store.mu.RUnlock()

	items := make([]AdminUser, len(store.data.Users))
	copy(items, store.data.Users)
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func AdminGetUser(w http.ResponseWriter, r *http.Request) {
	if adminDBEnabled() {
		item, err := getAdminUserFromDB(r.PathValue("id"))
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid_user_id")
			return
		}
		if item == nil {
			writeError(w, http.StatusNotFound, "user_not_found")
			return
		}
		writeJSON(w, http.StatusOK, item)
		return
	}

	id := r.PathValue("id")

	store.mu.RLock()
	defer store.mu.RUnlock()

	for _, item := range store.data.Users {
		if item.ID == id {
			writeJSON(w, http.StatusOK, item)
			return
		}
	}

	writeError(w, http.StatusNotFound, "user_not_found")
}

func AdminCreateUser(w http.ResponseWriter, r *http.Request) {
	var payload AdminUser
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}

	payload = normalizeUser(payload)
	if issues := validateUser(payload); len(issues) > 0 {
		writeValidationError(w, issues)
		return
	}

	if adminDBEnabled() {
		item, err := createAdminUserInDB(payload)
		if err != nil {
			writeValidationError(w, []string{err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{"created": item})
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	if emailExistsLocked(payload.Email, "") {
		writeValidationError(w, []string{"email must be unique"})
		return
	}

	payload.ID = nextID("u", &store.data.UserSeq)
	store.data.Users = append(store.data.Users, payload)
	if err := store.saveLocked(); err != nil {
		writeError(w, http.StatusInternalServerError, "cannot_persist_admin_store")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"created": payload})
}

func AdminUpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var payload AdminUser
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}

	if adminDBEnabled() {
		item, err := updateAdminUserInDB(id, payload)
		if err != nil {
			writeValidationError(w, []string{err.Error()})
			return
		}
		if item == nil {
			writeError(w, http.StatusNotFound, "user_not_found")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"updated": item})
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	for index, item := range store.data.Users {
		if item.ID != id {
			continue
		}
		payload = normalizeUser(payload)
		payload.ID = id
		payload.CreatedAt = defaultString(payload.CreatedAt, item.CreatedAt)
		if issues := validateUser(payload); len(issues) > 0 {
			writeValidationError(w, issues)
			return
		}
		if emailExistsLocked(payload.Email, id) {
			writeValidationError(w, []string{"email must be unique"})
			return
		}
		if item.Role == "Admin" && item.Status == "active" && (payload.Role != "Admin" || payload.Status != "active") && activeAdminCountLocked() <= 1 {
			writeValidationError(w, []string{"at least one active admin must remain available"})
			return
		}
		store.data.Users[index] = payload
		if err := store.saveLocked(); err != nil {
			writeError(w, http.StatusInternalServerError, "cannot_persist_admin_store")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"updated": payload})
		return
	}

	writeError(w, http.StatusNotFound, "user_not_found")
}

func AdminToggleUserStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if adminDBEnabled() {
		item, err := toggleAdminUserStatusInDB(id)
		if err != nil {
			writeValidationError(w, []string{err.Error()})
			return
		}
		if item == nil {
			writeError(w, http.StatusNotFound, "user_not_found")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"updated": item})
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	for index, item := range store.data.Users {
		if item.ID != id {
			continue
		}
		if item.Status == "active" {
			if item.Role == "Admin" && activeAdminCountLocked() <= 1 {
				writeValidationError(w, []string{"the last active admin cannot be deactivated"})
				return
			}
			store.data.Users[index].Status = "inactive"
		} else {
			store.data.Users[index].Status = "active"
		}
		if err := store.saveLocked(); err != nil {
			writeError(w, http.StatusInternalServerError, "cannot_persist_admin_store")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"updated": store.data.Users[index]})
		return
	}

	writeError(w, http.StatusNotFound, "user_not_found")
}

func AdminDeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if adminDBEnabled() {
		if err := deleteAdminUserInDB(id); err != nil {
			writeValidationError(w, []string{err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"deleted": id})
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	for index, item := range store.data.Users {
		if item.ID != id {
			continue
		}
		if item.Role == "Admin" && item.Status == "active" && activeAdminCountLocked() <= 1 {
			writeValidationError(w, []string{"the last active admin cannot be deleted"})
			return
		}
		store.data.Users = append(store.data.Users[:index], store.data.Users[index+1:]...)
		if err := store.saveLocked(); err != nil {
			writeError(w, http.StatusInternalServerError, "cannot_persist_admin_store")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"deleted": item.ID})
		return
	}

	writeError(w, http.StatusNotFound, "user_not_found")
}

func AdminListPrestations(w http.ResponseWriter, r *http.Request) {
	filterType := strings.TrimSpace(r.URL.Query().Get("type"))

	if adminDBEnabled() {
		items, err := listAdminPrestationsFromDB(filterType)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "cannot_read_prestations")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": items})
		return
	}

	store.mu.RLock()
	defer store.mu.RUnlock()

	items := make([]AdminPrestation, 0, len(store.data.Prestations))
	for _, item := range store.data.Prestations {
		if filterType != "" && item.Type != filterType {
			continue
		}
		items = append(items, item)
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func AdminGetPrestation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if adminDBEnabled() {
		item, err := getAdminPrestationFromDB(id)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid_prestation_id")
			return
		}
		if item == nil {
			writeError(w, http.StatusNotFound, "prestation_not_found")
			return
		}
		writeJSON(w, http.StatusOK, item)
		return
	}

	store.mu.RLock()
	defer store.mu.RUnlock()

	for _, item := range store.data.Prestations {
		if item.ID == id {
			writeJSON(w, http.StatusOK, item)
			return
		}
	}

	writeError(w, http.StatusNotFound, "prestation_not_found")
}

func AdminCreatePrestation(w http.ResponseWriter, r *http.Request) {
	var payload AdminPrestation
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}

	payload = normalizePrestation(payload)
	if issues := validatePrestation(payload); len(issues) > 0 {
		writeValidationError(w, issues)
		return
	}

	if adminDBEnabled() {
		item, err := createAdminPrestationInDB(payload)
		if err != nil {
			writeValidationError(w, []string{err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{"created": item})
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	payload.ID = nextID("p", &store.data.PrestationSeq)
	store.data.Prestations = append(store.data.Prestations, payload)
	if err := store.saveLocked(); err != nil {
		writeError(w, http.StatusInternalServerError, "cannot_persist_admin_store")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"created": payload})
}

func AdminUpdatePrestation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var payload AdminPrestation
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}

	if adminDBEnabled() {
		item, err := updateAdminPrestationInDB(id, payload)
		if err != nil {
			writeValidationError(w, []string{err.Error()})
			return
		}
		if item == nil {
			writeError(w, http.StatusNotFound, "prestation_not_found")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"updated": item})
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	for index, item := range store.data.Prestations {
		if item.ID != id {
			continue
		}
		payload = normalizePrestation(payload)
		payload.ID = id
		payload.CreatedAt = defaultString(payload.CreatedAt, item.CreatedAt)
		if issues := validatePrestation(payload); len(issues) > 0 {
			writeValidationError(w, issues)
			return
		}
		store.data.Prestations[index] = payload
		if err := store.saveLocked(); err != nil {
			writeError(w, http.StatusInternalServerError, "cannot_persist_admin_store")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"updated": payload})
		return
	}

	writeError(w, http.StatusNotFound, "prestation_not_found")
}

func AdminDeletePrestation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if adminDBEnabled() {
		if err := deleteAdminPrestationInDB(id); err != nil {
			writeValidationError(w, []string{err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"deleted": id})
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	for index, item := range store.data.Prestations {
		if item.ID != id {
			continue
		}
		store.data.Prestations = append(store.data.Prestations[:index], store.data.Prestations[index+1:]...)
		if err := store.saveLocked(); err != nil {
			writeError(w, http.StatusInternalServerError, "cannot_persist_admin_store")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"deleted": item.ID})
		return
	}

	writeError(w, http.StatusNotFound, "prestation_not_found")
}

func AdminListCategories(w http.ResponseWriter, r *http.Request) {
	if adminDBEnabled() {
		items, err := listAdminCategoriesFromDB()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "cannot_read_categories")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": items})
		return
	}
	store.mu.RLock()
	defer store.mu.RUnlock()

	items := make([]AdminCategory, len(store.data.Categories))
	copy(items, store.data.Categories)
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func AdminGetCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if adminDBEnabled() {
		item, err := getAdminCategoryFromDB(id)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid_category_id")
			return
		}
		if item == nil {
			writeError(w, http.StatusNotFound, "category_not_found")
			return
		}
		writeJSON(w, http.StatusOK, item)
		return
	}

	store.mu.RLock()
	defer store.mu.RUnlock()

	for _, item := range store.data.Categories {
		if item.ID == id {
			writeJSON(w, http.StatusOK, item)
			return
		}
	}

	writeError(w, http.StatusNotFound, "category_not_found")
}

func AdminCreateCategory(w http.ResponseWriter, r *http.Request) {
	var payload AdminCategory
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}

	payload = normalizeCategory(payload)
	if issues := validateCategory(payload); len(issues) > 0 {
		writeValidationError(w, issues)
		return
	}

	if adminDBEnabled() {
		item, err := createAdminCategoryInDB(payload)
		if err != nil {
			writeValidationError(w, []string{err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{"created": item})
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	if payload.ParentID != "" && !categoryExistsLocked(payload.ParentID) {
		writeValidationError(w, []string{"parentId references an unknown category"})
		return
	}

	payload.ID = nextID("c", &store.data.CategorySeq)
	store.data.Categories = append(store.data.Categories, payload)
	if err := store.saveLocked(); err != nil {
		writeError(w, http.StatusInternalServerError, "cannot_persist_admin_store")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"created": payload})
}

func AdminUpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var payload AdminCategory
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}

	if adminDBEnabled() {
		item, err := updateAdminCategoryInDB(id, payload)
		if err != nil {
			writeValidationError(w, []string{err.Error()})
			return
		}
		if item == nil {
			writeError(w, http.StatusNotFound, "category_not_found")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"updated": item})
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	for index, item := range store.data.Categories {
		if item.ID != id {
			continue
		}
		payload = normalizeCategory(payload)
		payload.ID = id
		if issues := validateCategory(payload); len(issues) > 0 {
			writeValidationError(w, issues)
			return
		}
		if payload.ParentID == id {
			writeValidationError(w, []string{"parentId cannot reference the current category"})
			return
		}
		if payload.ParentID != "" && !categoryExistsLocked(payload.ParentID) {
			writeValidationError(w, []string{"parentId references an unknown category"})
			return
		}
		if payload.Status == "" {
			payload.Status = item.Status
		}
		store.data.Categories[index] = payload
		if err := store.saveLocked(); err != nil {
			writeError(w, http.StatusInternalServerError, "cannot_persist_admin_store")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"updated": payload})
		return
	}

	writeError(w, http.StatusNotFound, "category_not_found")
}

func AdminDeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if adminDBEnabled() {
		if err := deleteAdminCategoryInDB(id); err != nil {
			writeValidationError(w, []string{err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"deleted": id})
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	for _, item := range store.data.Categories {
		if item.ParentID == id {
			writeValidationError(w, []string{"category has children and cannot be deleted"})
			return
		}
	}

	for index, item := range store.data.Categories {
		if item.ID != id {
			continue
		}
		store.data.Categories = append(store.data.Categories[:index], store.data.Categories[index+1:]...)
		if err := store.saveLocked(); err != nil {
			writeError(w, http.StatusInternalServerError, "cannot_persist_admin_store")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"deleted": item.ID})
		return
	}

	writeError(w, http.StatusNotFound, "category_not_found")
}

func AdminListEvents(w http.ResponseWriter, r *http.Request) {
	if adminDBEnabled() {
		items, err := listAdminEventsFromDB()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "cannot_read_events")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": items})
		return
	}
	store.mu.RLock()
	defer store.mu.RUnlock()

	items := make([]AdminEvent, len(store.data.Events))
	copy(items, store.data.Events)
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func AdminGetEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if adminDBEnabled() {
		item, err := getAdminEventFromDB(id)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid_event_id")
			return
		}
		if item == nil {
			writeError(w, http.StatusNotFound, "event_not_found")
			return
		}
		writeJSON(w, http.StatusOK, item)
		return
	}

	store.mu.RLock()
	defer store.mu.RUnlock()

	for _, item := range store.data.Events {
		if item.ID == id {
			writeJSON(w, http.StatusOK, item)
			return
		}
	}

	writeError(w, http.StatusNotFound, "event_not_found")
}

func AdminCreateEvent(w http.ResponseWriter, r *http.Request) {
	var payload AdminEvent
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}

	payload = normalizeEvent(payload)
	if issues := validateEvent(payload); len(issues) > 0 {
		writeValidationError(w, issues)
		return
	}

	if adminDBEnabled() {
		item, err := createAdminEventInDB(payload)
		if err != nil {
			writeValidationError(w, []string{err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{"created": item})
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	payload.ID = nextID("e", &store.data.EventSeq)
	store.data.Events = append(store.data.Events, payload)
	if err := store.saveLocked(); err != nil {
		writeError(w, http.StatusInternalServerError, "cannot_persist_admin_store")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"created": payload})
}

func AdminUpdateEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var payload AdminEvent
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}

	if adminDBEnabled() {
		item, err := updateAdminEventInDB(id, payload)
		if err != nil {
			writeValidationError(w, []string{err.Error()})
			return
		}
		if item == nil {
			writeError(w, http.StatusNotFound, "event_not_found")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"updated": item})
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	for index, item := range store.data.Events {
		if item.ID != id {
			continue
		}
		payload = normalizeEvent(payload)
		payload.ID = id
		if issues := validateEvent(payload); len(issues) > 0 {
			writeValidationError(w, issues)
			return
		}
		payload.Status = defaultString(payload.Status, item.Status)
		store.data.Events[index] = payload
		if err := store.saveLocked(); err != nil {
			writeError(w, http.StatusInternalServerError, "cannot_persist_admin_store")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"updated": payload})
		return
	}

	writeError(w, http.StatusNotFound, "event_not_found")
}

func AdminDeleteEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if adminDBEnabled() {
		if err := deleteAdminEventInDB(id); err != nil {
			writeValidationError(w, []string{err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"deleted": id})
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	for index, item := range store.data.Events {
		if item.ID != id {
			continue
		}
		store.data.Events = append(store.data.Events[:index], store.data.Events[index+1:]...)
		if err := store.saveLocked(); err != nil {
			writeError(w, http.StatusInternalServerError, "cannot_persist_admin_store")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"deleted": item.ID})
		return
	}

	writeError(w, http.StatusNotFound, "event_not_found")
}

func AdminModerationQueue(w http.ResponseWriter, r *http.Request) {
	if adminDBEnabled() {
		items, err := moderationQueueFromDB()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "cannot_read_moderation")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": items})
		return
	}
	store.mu.RLock()
	defer store.mu.RUnlock()

	writeJSON(w, http.StatusOK, map[string]any{"items": buildModerationQueueLocked()})
}

func AdminPublishPrestation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if adminDBEnabled() {
		item, err := publishAdminPrestationInDB(id)
		if err != nil {
			writeValidationError(w, []string{err.Error()})
			return
		}
		if item == nil {
			writeError(w, http.StatusNotFound, "prestation_not_found")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"updated": item})
		return
	}
	store.mu.Lock()
	defer store.mu.Unlock()

	for index, item := range store.data.Prestations {
		if item.ID != id {
			continue
		}
		if item.Status == "published" {
			writeJSON(w, http.StatusOK, map[string]any{"updated": item})
			return
		}
		item.Status = "published"
		if issues := validatePrestation(item); len(issues) > 0 {
			writeValidationError(w, issues)
			return
		}
		store.data.Prestations[index] = item
		if err := store.saveLocked(); err != nil {
			writeError(w, http.StatusInternalServerError, "cannot_persist_admin_store")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"updated": item})
		return
	}

	writeError(w, http.StatusNotFound, "prestation_not_found")
}

func AdminArchivePrestation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if adminDBEnabled() {
		item, err := archiveAdminPrestationInDB(id)
		if err != nil {
			writeValidationError(w, []string{err.Error()})
			return
		}
		if item == nil {
			writeError(w, http.StatusNotFound, "prestation_not_found")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"updated": item})
		return
	}
	store.mu.Lock()
	defer store.mu.Unlock()

	for index, item := range store.data.Prestations {
		if item.ID != id {
			continue
		}
		item.Status = "archived"
		store.data.Prestations[index] = item
		if err := store.saveLocked(); err != nil {
			writeError(w, http.StatusInternalServerError, "cannot_persist_admin_store")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"updated": item})
		return
	}

	writeError(w, http.StatusNotFound, "prestation_not_found")
}

func AdminPublishEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if adminDBEnabled() {
		item, err := publishAdminEventInDB(id)
		if err != nil {
			writeValidationError(w, []string{err.Error()})
			return
		}
		if item == nil {
			writeError(w, http.StatusNotFound, "event_not_found")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"updated": item})
		return
	}
	store.mu.Lock()
	defer store.mu.Unlock()

	for index, item := range store.data.Events {
		if item.ID != id {
			continue
		}
		item.Status = "published"
		if issues := validateEvent(item); len(issues) > 0 {
			writeValidationError(w, issues)
			return
		}
		store.data.Events[index] = item
		if err := store.saveLocked(); err != nil {
			writeError(w, http.StatusInternalServerError, "cannot_persist_admin_store")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"updated": item})
		return
	}

	writeError(w, http.StatusNotFound, "event_not_found")
}

func AdminArchiveEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if adminDBEnabled() {
		item, err := archiveAdminEventInDB(id)
		if err != nil {
			writeValidationError(w, []string{err.Error()})
			return
		}
		if item == nil {
			writeError(w, http.StatusNotFound, "event_not_found")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"updated": item})
		return
	}
	store.mu.Lock()
	defer store.mu.Unlock()

	for index, item := range store.data.Events {
		if item.ID != id {
			continue
		}
		item.Status = "archived"
		store.data.Events[index] = item
		if err := store.saveLocked(); err != nil {
			writeError(w, http.StatusInternalServerError, "cannot_persist_admin_store")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"updated": item})
		return
	}

	writeError(w, http.StatusNotFound, "event_not_found")
}

func AdminFinanceOverview(w http.ResponseWriter, r *http.Request) {
	if adminDBEnabled() {
		summary, items, err := financeOverviewFromDB()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "cannot_read_finance")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"summary": summary, "items": items})
		return
	}
	store.mu.RLock()
	defer store.mu.RUnlock()

	paidTotal := 0.0
	pendingTotal := 0.0
	activeSubscriptions := 0
	for _, item := range store.data.FinanceRecords {
		if item.Status == "paid" {
			paidTotal += item.Amount
		}
		if item.Status == "pending" {
			pendingTotal += item.Amount
		}
		if item.Category == "subscription" && item.Status == "paid" {
			activeSubscriptions += 1
		}
	}

	items := make([]AdminFinanceRecord, len(store.data.FinanceRecords))
	copy(items, store.data.FinanceRecords)
	writeJSON(w, http.StatusOK, map[string]any{
		"summary": map[string]any{
			"paidTotal":          paidTotal,
			"pendingTotal":       pendingTotal,
			"activeSubscriptions": activeSubscriptions,
		},
		"items": items,
	})
}

func AdminListNotifications(w http.ResponseWriter, r *http.Request) {
	if adminDBEnabled() {
		items, err := listAdminNotificationsFromDB()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "cannot_read_notifications")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": items})
		return
	}
	store.mu.RLock()
	defer store.mu.RUnlock()

	items := make([]AdminNotification, len(store.data.Notifications))
	copy(items, store.data.Notifications)
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func AdminCreateNotification(w http.ResponseWriter, r *http.Request) {
	var payload AdminNotification
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}

	payload = normalizeNotification(payload)
	if issues := validateNotification(payload); len(issues) > 0 {
		writeValidationError(w, issues)
		return
	}

	if adminDBEnabled() {
		item, err := createAdminNotificationInDB(payload)
		if err != nil {
			writeValidationError(w, []string{err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{"created": item})
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	payload.ID = nextID("n", &store.data.NotificationSeq)
	store.data.Notifications = append(store.data.Notifications, payload)
	if err := store.saveLocked(); err != nil {
		writeError(w, http.StatusInternalServerError, "cannot_persist_admin_store")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"created": payload})
}

func AdminUpdateNotificationStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var payload struct {
		Status string `json:"status"`
	}
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}

	if adminDBEnabled() {
		item, err := updateAdminNotificationStatusInDB(id, strings.TrimSpace(payload.Status))
		if err != nil {
			writeValidationError(w, []string{err.Error()})
			return
		}
		if item == nil {
			writeError(w, http.StatusNotFound, "notification_not_found")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"updated": item})
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	for index, item := range store.data.Notifications {
		if item.ID != id {
			continue
		}
		item.Status = strings.TrimSpace(payload.Status)
		item = normalizeNotification(item)
		if issues := validateNotification(item); len(issues) > 0 {
			writeValidationError(w, issues)
			return
		}
		store.data.Notifications[index] = item
		if err := store.saveLocked(); err != nil {
			writeError(w, http.StatusInternalServerError, "cannot_persist_admin_store")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"updated": item})
		return
	}

	writeError(w, http.StatusNotFound, "notification_not_found")
}

func AdminDeleteNotification(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if adminDBEnabled() {
		if err := deleteAdminNotificationInDB(id); err != nil {
			writeValidationError(w, []string{err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"deleted": id})
		return
	}
	store.mu.Lock()
	defer store.mu.Unlock()

	for index, item := range store.data.Notifications {
		if item.ID != id {
			continue
		}
		store.data.Notifications = append(store.data.Notifications[:index], store.data.Notifications[index+1:]...)
		if err := store.saveLocked(); err != nil {
			writeError(w, http.StatusInternalServerError, "cannot_persist_admin_store")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"deleted": item.ID})
		return
	}

	writeError(w, http.StatusNotFound, "notification_not_found")
}
