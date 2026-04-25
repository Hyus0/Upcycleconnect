package db

import (
	"database/sql"
	"fmt"
	"upcycleconnect/api-go/models"
	"math/rand"
	"time"
)

func ReserverUnCasier(annonceID int, siteID int) (string, error) {
	query := `
		SELECT c.id 
		FROM CASIER c
		JOIN CONTENEUR co ON c.id_conteneur = co.id
		WHERE co.id_site = ? 
		AND co.statut = 'Operationnel' 
		AND c.statut = 'Libre'
		LIMIT 1`

	var casierID int
	err := Conn.QueryRow(query, siteID).Scan(&casierID)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("plus de casiers disponibles sur ce site")
	}
	if err != nil {
		return "", err
	}

	rand.Seed(time.Now().UnixNano())
	pin := fmt.Sprintf("%04d", rand.Intn(10000))

	tx, err := Conn.Begin()
	if err != nil { return "", err }

	_, err = tx.Exec("UPDATE CASIER SET statut = 'Reserve' WHERE id = ?", casierID)
	if err != nil { tx.Rollback(); return "", err }

	_, err = tx.Exec(`
		UPDATE ANNONCE SET 
			id_casier = ?, 
			code_pin_depot = ?, 
			statut = 'En attente de depot' 
		WHERE id = ?`, casierID, pin, annonceID)
	if err != nil { tx.Rollback(); return "", err }

	err = tx.Commit()
	if err != nil { return "", err }

	return pin, nil
}

func GetAllSites() ([]models.Site, error) {
	query := "SELECT id, nom, ville, code_postal, adresse, telephone, type, actif FROM SITE WHERE actif = 1"
	rows, err := Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sites []models.Site
	for rows.Next() {
		var s models.Site
		if err := rows.Scan(&s.ID, &s.Nom, &s.Ville, &s.CodePostal, &s.Adresse, &s.Telephone, &s.Type, &s.Actif); err != nil {
			return nil, err
		}
		sites = append(sites, s)
	}
	return sites, nil
}

func GetConteneursBySite(siteID int) ([]models.Conteneur, error) {
	query := `SELECT id, id_site, type_dechet, statut, capacite_max_kg, niveau_remplissage 
	          FROM CONTENEUR WHERE id_site = ?`
	rows, err := Conn.Query(query, siteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conteneurs []models.Conteneur
	for rows.Next() {
		var c models.Conteneur
		if err := rows.Scan(&c.ID, &c.IdSite, &c.TypeDechet, &c.Statut, &c.CapaciteMaxKg, &c.NiveauRemplissage); err != nil {
			return nil, err
		}
		conteneurs = append(conteneurs, c)
	}
	return conteneurs, nil
}