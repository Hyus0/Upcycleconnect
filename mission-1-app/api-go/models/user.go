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