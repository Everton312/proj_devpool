package repository

import (
	"context"
	"uuid"

	"github.com/Everton312/proj-devpool/internal/models"
	"github.com/Everton312/proj-devpool/internal/repository/users"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	User interface {
		GetAll(ctx context.Context) ([]models.Users, error)
		GetById(ctx context.Context, id uuid.UUID) (*models.Users, error)
		Add(ctx context.Context, newUser *models.Users) error
		DeleteById(ctx context.Context, id uuid.UUID) error
		Update(ctx context.Context, user *models.Users) error
		ExistsByEmail(ctx context.Context, email string) (bool, error)
		ExistsByEmailExcludingID(ctx context.Context, email string, id uuid.UUID) (bool, error)
		ExistsByLogin(ctx context.Context, login string) (bool, error)
		ExistsByLoginExcludingID(ctx context.Context, login string, id uuid.UUID) (bool, error)
	}
}

func New(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		User: users.New(db),
	}
}
