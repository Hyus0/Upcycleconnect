package models

type User struct {
	Id              int
	Prenom          string
	Nom             string
	Password        string
	Mail            string
	Adresse         string
	Ville           string
	Code_postal     string
	Date_naissance  string 
	Date_inscription string
	Role            string
	Id_langue       int
}

type GetUser struct {
	Id              int
	Prenom          string
	Nom             string
	Mail            string
	Adresse         string
	Ville           string
	Code_postal     string
	Date_naissance  string
	Date_inscription string
	Role            string
	Id_langue       int
}