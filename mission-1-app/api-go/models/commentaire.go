package models

type Commentaire struct {
    ID          int    `json:"id"`
    Description string `json:"description"`
    Auteur      string `json:"author"`
}