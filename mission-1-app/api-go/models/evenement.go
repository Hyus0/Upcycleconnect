package models

import "time"

type Evenement struct {
	ID            int       `json:"id"`
	Titre         string    `json:"titre"`
	Description   string    `json:"description"`
	Adresse       string    `json:"adresse"`
	Ville         string    `json:"ville"`
	CodePostal    string    `json:"code_postal"`
	DateCreation  time.Time `json:"date_creation"`
	DateEvenement time.Time `json:"date_evenement"`
	Type          string    `json:"type"`
}
