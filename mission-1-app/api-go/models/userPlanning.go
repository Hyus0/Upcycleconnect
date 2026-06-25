package models

type UserPlanning struct {
	Formations []GetFormation `json:"formations"`
	Evenements []Evenement    `json:"evenements"`
}