package repository

import (
	"database/sql"

	"github.com/aeva-eat/backend/internal/model"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(username, displayName, passwordHash string) (*model.User, error) {
	u := &model.User{}
	var name *string
	if displayName != "" {
		name = &displayName
	}
	err := r.db.QueryRow(
		`INSERT INTO users (username, password_hash, display_name)
		 VALUES ($1, $2, $3)
		 RETURNING id, username, display_name, avatar_url, created_at, updated_at`,
		username, passwordHash, name,
	).Scan(&u.ID, &u.Username, &u.DisplayName, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) GetByUsername(username string) (*model.User, error) {
	u := &model.User{}
	err := r.db.QueryRow(
		`SELECT id, username, display_name, password_hash, avatar_url, created_at, updated_at
		 FROM users WHERE username = $1`,
		username,
	).Scan(&u.ID, &u.Username, &u.DisplayName, &u.PasswordHash, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) GetByID(id int) (*model.User, error) {
	u := &model.User{}
	err := r.db.QueryRow(
		`SELECT id, username, display_name, avatar_url, created_at, updated_at
		 FROM users WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.Username, &u.DisplayName, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}
