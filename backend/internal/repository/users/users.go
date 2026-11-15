package users

import (
	"context"

	"github.com/Everton312/proj-devpool/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Users struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Users {
	return &Users{db: db}
}

func (u *Users) GetAll(ctx context.Context) ([]models.Users, error) {
	return nil, nil
}

func (u Users) Add(ctx context.Context, newUser *models.Users) error {
	return nil
}
