package app

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
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

type adminStore struct {
	mu           sync.RWMutex
	users        []AdminUser
	prestations  []AdminPrestation
	categories   []AdminCategory
	events       []AdminEvent
	userSeq      int
	prestationSeq int
	categorySeq  int
	eventSeq     int
}

var store = newAdminStore()

func newAdminStore() *adminStore {
	return &adminStore{
		users: []AdminUser{
			{ID: "u1", FirstName: "Alice", LastName: "Martin", Email: "alice@upcycleconnect.local", City: "Paris", PostalCode: "75011", Role: "Particulier", Status: "active", CreatedAt: "2026-03-01"},
			{ID: "u2", FirstName: "Karim", LastName: "Benali", Email: "karim@atelier.local", City: "Lyon", PostalCode: "69002", Role: "Prestataire", Status: "active", CreatedAt: "2026-02-18"},
			{ID: "u3", FirstName: "Nina", LastName: "Roux", Email: "nina.admin@upcycleconnect.local", City: "Lille", PostalCode: "59000", Role: "Admin", Status: "inactive", CreatedAt: "2026-01-14"},
		},
		prestations: []AdminPrestation{
			{ID: "p1", Title: "Diagnostic mobilier", Description: "Analyse d'etat et estimation de remise en valeur.", Type: "service", Price: 45, Status: "published", Provider: "Atelier Renouveau", CreatedAt: "2026-03-04"},
			{ID: "p2", Title: "Lot textile recycle", Description: "Selection textile revalorisee pour ateliers creatifs.", Type: "vente", Price: 120, Status: "draft", Provider: "ReTex", CreatedAt: "2026-03-08"},
		},
		categories: []AdminCategory{
			{ID: "c1", Name: "Mobilier", ParentID: "", Description: "Meubles et accessoires", Status: "active"},
			{ID: "c2", Name: "Chaises", ParentID: "c1", Description: "Assises et tabourets", Status: "active"},
		},
		events: []AdminEvent{
			{ID: "e1", Title: "Atelier bois", Location: "Paris 11", Date: "2026-03-28", Status: "planned", Capacity: 18, Description: "Session pratique autour de la remise en etat du bois."},
			{ID: "e2", Title: "Collecte textile", Location: "Lyon 2", Date: "2026-04-02", Status: "published", Capacity: 40, Description: "Journee de collecte et tri textile."},
		},
		userSeq:       3,
		prestationSeq: 2,
		categorySeq:   2,
		eventSeq:      2,
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
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

func AdminMetrics(w http.ResponseWriter, r *http.Request) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	writeJSON(w, http.StatusOK, map[string]any{
		"metrics": map[string]int{
			"users":     len(store.users),
			"annonces":  len(store.prestations),
			"categories": len(store.categories),
			"events":    len(store.events),
		},
	})
}

func AdminListUsers(w http.ResponseWriter, r *http.Request) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	items := make([]AdminUser, len(store.users))
	copy(items, store.users)
	for i := range items {
		items[i].FullName = fullName(items[i].FirstName, items[i].LastName)
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func AdminCreateUser(w http.ResponseWriter, r *http.Request) {
	var payload AdminUser
	if err := decodeJSON(r, &payload); err != nil {
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	payload.ID = nextID("u", &store.userSeq)
	payload.CreatedAt = defaultString(payload.CreatedAt, "2026-04-02")
	payload.Status = defaultString(payload.Status, "active")
	payload.Role = defaultString(payload.Role, "Particulier")
	payload.FullName = fullName(payload.FirstName, payload.LastName)
	store.users = append(store.users, payload)

	writeJSON(w, http.StatusCreated, map[string]any{"created": payload})
}

func AdminUpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var payload AdminUser
	if err := decodeJSON(r, &payload); err != nil {
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	for index, user := range store.users {
		if user.ID != id {
			continue
		}

		payload.ID = id
		payload.CreatedAt = defaultString(payload.CreatedAt, user.CreatedAt)
		payload.Status = defaultString(payload.Status, user.Status)
		payload.Role = defaultString(payload.Role, user.Role)
		payload.FullName = fullName(payload.FirstName, payload.LastName)
		store.users[index] = payload
		writeJSON(w, http.StatusOK, map[string]any{"updated": payload})
		return
	}

	http.Error(w, "Utilisateur introuvable", http.StatusNotFound)
}

func AdminToggleUserStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	store.mu.Lock()
	defer store.mu.Unlock()

	for index, user := range store.users {
		if user.ID != id {
			continue
		}
		if user.Status == "active" {
			store.users[index].Status = "inactive"
		} else {
			store.users[index].Status = "active"
		}
		writeJSON(w, http.StatusOK, map[string]any{"updated": store.users[index]})
		return
	}

	http.Error(w, "Utilisateur introuvable", http.StatusNotFound)
}

func AdminDeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	store.mu.Lock()
	defer store.mu.Unlock()

	for index, user := range store.users {
		if user.ID != id {
			continue
		}
		store.users = append(store.users[:index], store.users[index+1:]...)
		writeJSON(w, http.StatusOK, map[string]any{"deleted": user.ID})
		return
	}

	http.Error(w, "Utilisateur introuvable", http.StatusNotFound)
}

func AdminListPrestations(w http.ResponseWriter, r *http.Request) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	filterType := strings.TrimSpace(r.URL.Query().Get("type"))
	items := make([]AdminPrestation, 0, len(store.prestations))
	for _, item := range store.prestations {
		if filterType != "" && item.Type != filterType {
			continue
		}
		items = append(items, item)
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func AdminCreatePrestation(w http.ResponseWriter, r *http.Request) {
	var payload AdminPrestation
	if err := decodeJSON(r, &payload); err != nil {
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	payload.ID = nextID("p", &store.prestationSeq)
	payload.Status = defaultString(payload.Status, "draft")
	payload.Provider = defaultString(payload.Provider, "Equipe locale")
	payload.CreatedAt = defaultString(payload.CreatedAt, "2026-04-02")
	store.prestations = append(store.prestations, payload)

	writeJSON(w, http.StatusCreated, map[string]any{"created": payload})
}

func AdminUpdatePrestation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var payload AdminPrestation
	if err := decodeJSON(r, &payload); err != nil {
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	for index, item := range store.prestations {
		if item.ID != id {
			continue
		}
		payload.ID = id
		payload.Provider = defaultString(payload.Provider, item.Provider)
		payload.CreatedAt = defaultString(payload.CreatedAt, item.CreatedAt)
		payload.Status = defaultString(payload.Status, item.Status)
		store.prestations[index] = payload
		writeJSON(w, http.StatusOK, map[string]any{"updated": payload})
		return
	}

	http.Error(w, "Prestation introuvable", http.StatusNotFound)
}

func AdminDeletePrestation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	store.mu.Lock()
	defer store.mu.Unlock()

	for index, item := range store.prestations {
		if item.ID != id {
			continue
		}
		store.prestations = append(store.prestations[:index], store.prestations[index+1:]...)
		writeJSON(w, http.StatusOK, map[string]any{"deleted": item.ID})
		return
	}

	http.Error(w, "Prestation introuvable", http.StatusNotFound)
}

func AdminListCategories(w http.ResponseWriter, r *http.Request) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	items := make([]AdminCategory, len(store.categories))
	copy(items, store.categories)
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func AdminCreateCategory(w http.ResponseWriter, r *http.Request) {
	var payload AdminCategory
	if err := decodeJSON(r, &payload); err != nil {
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	payload.ID = nextID("c", &store.categorySeq)
	payload.Status = defaultString(payload.Status, "active")
	store.categories = append(store.categories, payload)

	writeJSON(w, http.StatusCreated, map[string]any{"created": payload})
}

func AdminUpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var payload AdminCategory
	if err := decodeJSON(r, &payload); err != nil {
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	for index, item := range store.categories {
		if item.ID != id {
			continue
		}
		payload.ID = id
		payload.Status = defaultString(payload.Status, item.Status)
		store.categories[index] = payload
		writeJSON(w, http.StatusOK, map[string]any{"updated": payload})
		return
	}

	http.Error(w, "Categorie introuvable", http.StatusNotFound)
}

func AdminDeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	store.mu.Lock()
	defer store.mu.Unlock()

	for index, item := range store.categories {
		if item.ID != id {
			continue
		}
		store.categories = append(store.categories[:index], store.categories[index+1:]...)
		writeJSON(w, http.StatusOK, map[string]any{"deleted": item.ID})
		return
	}

	http.Error(w, "Categorie introuvable", http.StatusNotFound)
}

func AdminListEvents(w http.ResponseWriter, r *http.Request) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	items := make([]AdminEvent, len(store.events))
	copy(items, store.events)
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func AdminCreateEvent(w http.ResponseWriter, r *http.Request) {
	var payload AdminEvent
	if err := decodeJSON(r, &payload); err != nil {
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	payload.ID = nextID("e", &store.eventSeq)
	payload.Status = defaultString(payload.Status, "planned")
	store.events = append(store.events, payload)

	writeJSON(w, http.StatusCreated, map[string]any{"created": payload})
}

func AdminUpdateEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var payload AdminEvent
	if err := decodeJSON(r, &payload); err != nil {
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	for index, item := range store.events {
		if item.ID != id {
			continue
		}
		payload.ID = id
		payload.Status = defaultString(payload.Status, item.Status)
		store.events[index] = payload
		writeJSON(w, http.StatusOK, map[string]any{"updated": payload})
		return
	}

	http.Error(w, "Evenement introuvable", http.StatusNotFound)
}

func AdminDeleteEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	store.mu.Lock()
	defer store.mu.Unlock()

	for index, item := range store.events {
		if item.ID != id {
			continue
		}
		store.events = append(store.events[:index], store.events[index+1:]...)
		writeJSON(w, http.StatusOK, map[string]any{"deleted": item.ID})
		return
	}

	http.Error(w, "Evenement introuvable", http.StatusNotFound)
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
