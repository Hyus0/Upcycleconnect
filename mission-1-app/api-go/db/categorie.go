package db

import (
	"database/sql"
	"fmt"
	"upcycleconnect/api-go/models"
)

func GetAllCategories() ([]models.Category, error) {
	categories := []models.Category{}
	query := "SELECT id, nom, description, id_parent, statut, date_creation FROM CATEGORIE"
	rows, err := Conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetAllCategories query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.ID, &c.Nom, &c.Description, &c.ParentID, &c.Statut, &c.DateCreation)
		if err != nil {
			return nil, fmt.Errorf("GetAllCategories scan: %v", err)
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func GetCategory(categoryId int) (*models.Category, error) {
	var c models.Category
	query := "SELECT id, nom, description, id_parent, statut, date_creation FROM CATEGORIE WHERE id = ?"
	row := Conn.QueryRow(query, categoryId)

	err := row.Scan(&c.ID, &c.Nom, &c.Description, &c.ParentID, &c.Statut, &c.DateCreation)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GetCategory scan: %v", err)
	}

	return &c, nil
}

func CreateCategory(category models.Category) error {
	query := `INSERT INTO CATEGORIE (nom, description, id_parent, statut, date_creation) 
	          VALUES (?, ?, ?, ?, NOW())`

	_, err := Conn.Exec(query,
		category.Nom,
		category.Description,
		category.ParentID,
		category.Statut,
	)

	if err != nil {
		return fmt.Errorf("CreateCategory: %v", err)
	}
	return nil
}

func ModifyCategory(categoryId int, category models.Category) error {
	query := `UPDATE CATEGORIE SET 
				nom = ?, 
				description = ?, 
				id_parent = ?, 
				statut = ?
			  WHERE id = ?`
              
	result, err := Conn.Exec(query,
		category.Nom,
		category.Description,
		category.ParentID,
		category.Statut,
		categoryId,
	)

	if err != nil {
		return fmt.Errorf("ModifyCategory: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ModifyCategory (RowsAffected): %v", err)
	}

	if rowsAffected == 0 {
		fmt.Printf("Note: Catégorie %d non modifiée\n", categoryId)
	}

	return nil
}

func DeleteCategory(id int) error {
	if id <= 0 {
		return fmt.Errorf("DeleteCategory: L'ID doit être un entier positif")
	}
	result, err := Conn.Exec("DELETE FROM CATEGORIE WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("DeleteCategory: échec de la suppression : %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteCategory: erreur de RowsAffected : %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("DeleteCategory: aucune catégorie trouvée avec l'ID %d", id)
	}
	return nil
}