package service

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"uuid"

	"github.com/Everton312/proj-devpool/internal/models"
	"github.com/Everton312/proj-devpool/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type Services struct {
	repos *repository.Repositories
}

func New(repos *repository.Repositories) *Services {
	return &Services{
		repos: repos,
	}
}

func (u Services) GetAll(ctx context.Context) ([]models.Users, error) {
	users, err := u.repos.User.GetAll(ctx)
	return users, err
}

func (s *Services) Add(ctx context.Context, input models.CreateUsersInput) (*models.Users, error) {
	//tratando inputs
	name := strings.TrimSpace(input.Name)
	if len(name) < 3 {
		return nil, errors.New("o nome deve ter pelo menos 3 caracteres")
	}

	login := strings.TrimSpace(input.Login)
	if len(login) < 3 {
		return nil, errors.New("o login deve ter pelo menos 3 caracteres")
	}

	email := strings.TrimSpace(input.Email)
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, errors.New("e-mail inválido")
	}

	if len(input.Password) < 8 {
		return nil, errors.New("a senha deve ter pelo menos 8 caracteres")
	}

	//verifica se já tá cadastrado o email e login usando as funcs do repository
	emailExists, err := s.repos.User.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if emailExists {
		return nil, errors.New("este e-mail já está cadastrado")
	}

	loginExists, err := s.repos.User.ExistsByLogin(ctx, login)
	if err != nil {
		return nil, err
	}
	if loginExists {
		return nil, errors.New("este login já está cadastrado")
	}

	//BCrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	active := true
	if input.Active != nil {
		active = *input.Active
	}

	newUser := &models.Users{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		Login:        login,
		PasswordHash: string(hash),
		Active:       active,
	}

	if err := s.repos.User.Add(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}
