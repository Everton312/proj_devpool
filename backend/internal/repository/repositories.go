package repository

import (
	"context"

	"github.com/Everton312/proj-devpool/internal/models"
	"github.com/Everton312/proj-devpool/internal/repository/users"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	User interface {
		GetAll(ctx context.Context) ([]models.Users, error)
		Add(ctx context.Context, newUser *models.Users) error
	}
}

func New(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		User: users.New(db),
	}
}
