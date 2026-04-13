package repository

import (
	"database/sql"

	"github.com/aeva-eat/backend/internal/model"
)

type CatalogRepo struct {
	db *sql.DB
}

func NewCatalogRepo(db *sql.DB) *CatalogRepo {
	return &CatalogRepo{db: db}
}

func (r *CatalogRepo) ListCuisineTypes() ([]model.CuisineType, error) {
	rows, err := r.db.Query(`SELECT id, name FROM cuisine_types ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.CuisineType
	for rows.Next() {
		var item model.CuisineType
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *CatalogRepo) ListCategories() ([]model.Category, error) {
	rows, err := r.db.Query(`SELECT id, name FROM categories ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.Category
	for rows.Next() {
		var item model.Category
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
