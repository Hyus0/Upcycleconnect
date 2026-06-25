package models

type Langue struct {
	ID        int    `json:"id"`
	Code      string `json:"code"`
	NomLangue string `json:"nom_langue"`
}