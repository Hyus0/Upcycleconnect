package models

import "time"

type Annonce struct {
	ID                 int        `json:"id" db:"id"`
	IdVendeur          int        `json:"id_vendeur" db:"id_vendeur"`
	IdAcheteur         *int       `json:"id_acheteur,omitempty" db:"id_acheteur"`
	IdCasier           *int       `json:"id_casier,omitempty" db:"id_casier"`
	IdCategorie        int        `json:"id_categorie" db:"id_categorie"`
	Titre              string     `json:"titre" db:"titre"`
	Description        string     `json:"description" db:"description"`
	TypeMateriau       string     `json:"type_materiau" db:"type_materiau"`
	PoidsEstimeKg      float64    `json:"poids_estime_kg" db:"poids_estime_kg"`
	Prix               float64    `json:"prix" db:"prix"`
	EtatObjet          string     `json:"etat_objet" db:"etat_objet"`
	Statut             string     `json:"statut" db:"statut"`
	EstValide          string     `json:"est_valide" db:"est_valide"`
	CodePinDepot       string     `json:"code_pin_depot" db:"code_pin_depot"`
	CodeBarreRetrait   string     `json:"code_barre_retrait" db:"code_barre_retrait"`
	DateCreation       time.Time  `json:"date_creation" db:"date_creation"`
	DateDepotEffective *time.Time `json:"date_depot_effective,omitempty" db:"date_depot_effective"`
	Type               string     `json:"type" db:"type"`
	Ville              string     `json:"ville" db:"ville"`
	CodePostal         string     `json:"code_postal" db:"code_postal"`
	Adresse            string     `json:"adresse" db:"adresse"`
}
