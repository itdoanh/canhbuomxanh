package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type userRecord struct {
	ID           uint64
	Email        string
	PasswordHash string
	Role         string
}

type repository struct {
	db *sql.DB
}

func newRepository(db *sql.DB) *repository {
	return &repository{db: db}
}

func (r *repository) findUserByEmail(ctx context.Context, email string) (*userRecord, error) {
	row := r.db.QueryRowContext(ctx,
		"SELECT id, email, password_hash, role FROM users WHERE email = ? AND status = 'active' LIMIT 1",
		email,
	)

	var user userRecord
	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("find user by email failed: %w", err)
	}

	return &user, nil
}

func (r *repository) createUser(ctx context.Context, req RegisterRequest, passwordHash string) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO users (email, password_hash, full_name, role) VALUES (?, ?, ?, ?)",
		req.Email,
		passwordHash,
		req.FullName,
		req.Role,
	)
	if err != nil {
		return fmt.Errorf("create user failed: %w", err)
	}
	return nil
}
