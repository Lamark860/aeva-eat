package repository

import (
	"database/sql"

	"github.com/aeva-eat/backend/internal/model"
)

type InviteRepo struct {
	db *sql.DB
}

func NewInviteRepo(db *sql.DB) *InviteRepo {
	return &InviteRepo{db: db}
}

func (r *InviteRepo) Create(code string, createdBy int, expiresAt *string) (*model.Invite, error) {
	inv := &model.Invite{}
	var err error
	if expiresAt != nil {
		err = r.db.QueryRow(
			`INSERT INTO invites (code, created_by, expires_at)
			 VALUES ($1, $2, $3::timestamptz)
			 RETURNING id, code, created_by, used_by, used_at, expires_at, created_at`,
			code, createdBy, *expiresAt,
		).Scan(&inv.ID, &inv.Code, &inv.CreatedBy, &inv.UsedBy, &inv.UsedAt, &inv.ExpiresAt, &inv.CreatedAt)
	} else {
		err = r.db.QueryRow(
			`INSERT INTO invites (code, created_by)
			 VALUES ($1, $2)
			 RETURNING id, code, created_by, used_by, used_at, expires_at, created_at`,
			code, createdBy,
		).Scan(&inv.ID, &inv.Code, &inv.CreatedBy, &inv.UsedBy, &inv.UsedAt, &inv.ExpiresAt, &inv.CreatedAt)
	}
	if err != nil {
		return nil, err
	}
	return inv, nil
}

func (r *InviteRepo) GetByCode(code string) (*model.Invite, error) {
	inv := &model.Invite{}
	err := r.db.QueryRow(
		`SELECT i.id, i.code, i.created_by, i.used_by, i.used_at, i.expires_at, i.created_at,
		        u.username
		 FROM invites i
		 JOIN users u ON u.id = i.created_by
		 WHERE i.code = $1`,
		code,
	).Scan(&inv.ID, &inv.Code, &inv.CreatedBy, &inv.UsedBy, &inv.UsedAt, &inv.ExpiresAt, &inv.CreatedAt, &inv.CreatorName)
	if err != nil {
		return nil, err
	}
	return inv, nil
}

func (r *InviteRepo) MarkUsed(code string, userID int) error {
	_, err := r.db.Exec(
		`UPDATE invites SET used_by = $1, used_at = now() WHERE code = $2`,
		userID, code,
	)
	return err
}

func (r *InviteRepo) ListByCreator(creatorID int) ([]model.Invite, error) {
	rows, err := r.db.Query(
		`SELECT i.id, i.code, i.created_by, i.used_by, i.used_at, i.expires_at, i.created_at,
		        COALESCE(u2.username, '')
		 FROM invites i
		 LEFT JOIN users u2 ON u2.id = i.used_by
		 WHERE i.created_by = $1
		 ORDER BY i.created_at DESC`,
		creatorID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invites []model.Invite
	for rows.Next() {
		var inv model.Invite
		var usedByName string
		if err := rows.Scan(&inv.ID, &inv.Code, &inv.CreatedBy, &inv.UsedBy, &inv.UsedAt, &inv.ExpiresAt, &inv.CreatedAt, &usedByName); err != nil {
			return nil, err
		}
		if usedByName != "" {
			inv.CreatorName = usedByName
		}
		invites = append(invites, inv)
	}
	return invites, nil
}

func (r *InviteRepo) ListAll() ([]model.Invite, error) {
	rows, err := r.db.Query(
		`SELECT i.id, i.code, i.created_by, i.used_by, i.used_at, i.expires_at, i.created_at,
		        u.username
		 FROM invites i
		 JOIN users u ON u.id = i.created_by
		 ORDER BY i.created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invites []model.Invite
	for rows.Next() {
		var inv model.Invite
		if err := rows.Scan(&inv.ID, &inv.Code, &inv.CreatedBy, &inv.UsedBy, &inv.UsedAt, &inv.ExpiresAt, &inv.CreatedAt, &inv.CreatorName); err != nil {
			return nil, err
		}
		invites = append(invites, inv)
	}
	return invites, nil
}

func (r *InviteRepo) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM invites WHERE id = $1`, id)
	return err
}
