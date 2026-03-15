package models

type Category struct {
    ID           int       `json:"id"`
    Nom          string    `json:"nom"`
    Description  string    `json:"description"`
    ParentID     *int      `json:"parent_id"` 
    Statut       string    `json:"statut"` 
    DateCreation string    `json:"date_creation"`
}