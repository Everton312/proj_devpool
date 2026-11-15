package users

import (
	"context"
	"errors"
	"uuid"

	"github.com/Everton312/proj-devpool/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Users struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Users {
	return &Users{db: db}
}

func (u *Users) GetAll(ctx context.Context) ([]models.Users, error) {
	query := `
		SELECT id, name, email, login, active, created_at FROM users
		ORDER BY created_at DESC
	`
	rows, err := u.db.Query(ctx, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []models.Users{}

	for rows.Next() {
		var user models.Users
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Login, &user.Active, &user.CreatedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (u Users) Add(ctx context.Context, newUser *models.Users) error {
	query := `
		INSERT INTO users (id, name, email, login, password_hash, active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at
	`
	return u.db.QueryRow(
		ctx,
		query,
		newUser.ID,
		newUser.Name,
		newUser.Email,
		newUser.Login,
		newUser.PasswordHash,
		newUser.Active,
	).Scan(&newUser.CreatedAt)
}

func (u *Users) GetById(ctx context.Context, id uuid.UUID) (*models.Users, error) {
	query := `
		SELECT id, name, email, login, password_hash, active, created_at FROM users WHERE id = $1
	
	`
	var user models.Users
	err := u.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Login, &user.PasswordHash, &user.Active, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil

}

func (u *Users) DeleteById(ctx context.Context, id uuid.UUID) error {
	cmdTag, err := u.db.Exec(ctx, "DELETE FROM users WHERE id = $1", id)

	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return models.ErrUserNotFound
	}

	return nil

}

func (u *Users) Update(ctx context.Context, user *models.Users) error {
	query := `
		UPDATE users SET name = $1, email = $2, login = $3, password_hash = $4, active = $5 WHERE id = $6
		RETURNING created_at
	`

	err := u.db.QueryRow(ctx, query, user.Name, user.Email, user.Login, user.PasswordHash, user.Active, user.ID).Scan(&user.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ErrUserNotFound
		}
		return err
	}
	return nil

}

func (u *Users) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	err := u.db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (u *Users) ExistsByEmailExcludingID(ctx context.Context, email string, id uuid.UUID) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND id <> $2)
	`
	err := u.db.QueryRow(ctx, query, email, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, err
}

func (u *Users) ExistsByLoginExcludingID(ctx context.Context, login string, id uuid.UUID) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(SELECT 1 FROM users WHERE login = $1 AND id <> $2)
	`
	err := u.db.QueryRow(ctx, query, login, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, err
}

func (u *Users) ExistsByLogin(ctx context.Context, login string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE login = $1)`

	err := u.db.QueryRow(ctx, query, login).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
