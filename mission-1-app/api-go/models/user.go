package models

type User struct {
	Id               int    `json:"id"`
	Prenom           string `json:"prenom"`
	Nom              string `json:"nom"`
	Password         string `json:"password,omitempty"` 
	Mail             string `json:"mail"`
	Adresse          string `json:"adresse"`
	Ville            string `json:"ville"`
	CodePostal       string `json:"code_postal"`
	DateNaissance    string `json:"date_naissance"`  
	DateInscription  string `json:"date_inscription"` 
	DateUpdatePassword string `json:"date_update_password"`
	Role             string `json:"role"`
	Statut           string `json:"statut"`
	IdLangue         int    `json:"id_langue"`
}

type GetUser struct {
	Id               int    `json:"id"`
	Prenom           string `json:"prenom"`
	Nom              string `json:"nom"`
	Mail             string `json:"mail"`
	Adresse          string `json:"adresse"`
	Ville            string `json:"ville"`
	CodePostal       string `json:"code_postal"`
	DateNaissance    string `json:"date_naissance"`  
	DateInscription  string `json:"date_inscription"`
	DateUpdatePassword string `json:"date_update_password"`
	Role             string `json:"role"`
	Statut           string `json:"statut"`
	IdLangue         int    `json:"id_langue"`
}

type UserStats struct {
	Points           int     `json:"total_points"`
	Niveau           string  `json:"niveau"`
	Co2Evite         float64 `json:"co2_total_evite_kg"`
	ObjetsRecycles   int     `json:"nb_objets_recycles"`
	ArgentEconomise  float64 `json:"ressources_economisees"`
}

type UserPlanningEntry struct {
	ID         int    `json:"id"`
	Kind       string `json:"kind"`
	Title      string `json:"title"`
	DateTime   string `json:"date_time"`
	Address    string `json:"address"`
	City       string `json:"city"`
	PostalCode string `json:"postal_code"`
}

type PlatformOverview struct {
	MembersCount         int     `json:"members_count"`
	PrestatairesCount    int     `json:"prestataires_count"`
	UpcycledObjectsCount int     `json:"upcycled_objects_count"`
	Co2SavedKg           float64 `json:"co2_saved_kg"`
	SitesCount           int     `json:"sites_count"`
	AnnoncesCount        int     `json:"annonces_count"`
	EventsCount          int     `json:"events_count"`
	FormationsCount      int     `json:"formations_count"`
}
