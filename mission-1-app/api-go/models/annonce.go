package models

import "time"

type Annonce struct {
	ID           int       `json:"id" db:"id"`
	IdVendeur    int       `json:"id_vendeur" db:"idvendeur"`
	IdAcheteur   *int      `json:"id_acheteur,omitempty" db:"idacheteur"`
	Titre        string    `json:"titre" db:"titre"`
	Description  string    `json:"description" db:"description"`
	Statut       string    `json:"statut" db:"statut"`
	EstValide    string    `json:"est_valide" db:"estvalide"`
	Prix         float64   `json:"prix" db:"prix"`
	EtatObjet    string    `json:"etat_objet" db:"etatobjet"`
	Adresse      string    `json:"adresse" db:"adresse"`
	Ville        string    `json:"ville" db:"ville"`
	CodePostal   string    `json:"code_postal" db:"codepostal"`
	DateCreation time.Time `json:"date_creation" db:"datecreation"`
	Type         string    `json:"type" db:"type"`
}
