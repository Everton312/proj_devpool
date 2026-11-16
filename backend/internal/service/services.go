package service

import (
	"context"
	"errors"
	"fmt"
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

func (u Services) GetById(ctx context.Context, id uuid.UUID) (*models.Users, error) {
	user, err := u.repos.User.GetById(ctx, id)
	return user, err
}

func (u Services) DeleteById(ctx context.Context, id uuid.UUID) error {
	err := u.repos.User.DeleteById(ctx, id)
	return err
}

func (u Services) Update(ctx context.Context, id uuid.UUID, input models.UpdateUsersInput) (*models.Users, error) {
	user, err := u.repos.User.GetById(ctx, id)

	if err != nil {
		return nil, err
	}

	name := strings.TrimSpace(input.Name)
	if len(name) < 3 {
		return nil, fmt.Errorf("%w: o nome deve ter pelo menos 3 caracteres", models.ErrValidation)
	}

	login := strings.TrimSpace(input.Login)
	if len(login) < 3 {
		return nil, fmt.Errorf("%w: o login deve ter pelo menos 3 caracteres", models.ErrValidation)
	}

	email := strings.ToLower(strings.TrimSpace(input.Email))
	if !isValidEmail(email) {
		return nil, fmt.Errorf("%w: e-mail inválido", models.ErrValidation)
	}
	emailExists, err := u.repos.User.ExistsByEmailExcludingID(ctx, email, id)
	if err != nil {
		return nil, err
	}
	if emailExists {
		return nil, fmt.Errorf("%w: email já cadastrado", models.ErrValidation)
	}

	loginExists, err := u.repos.User.ExistsByLoginExcludingID(ctx, login, id)
	if err != nil {
		return nil, err
	}
	if loginExists {
		return nil, fmt.Errorf("%w: login já cadastrado", models.ErrValidation)
	}

	if input.NewPassword != "" {
		if input.CurrentPassword == "" {
			return nil, fmt.Errorf("%w: a senha atual é obrigatória para definir uma nova senha", models.ErrValidation)
		}

		err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.CurrentPassword))
		if err != nil {
			return nil, fmt.Errorf("%w: senha atual incorreta", models.ErrInvalidPassword)
		}

		if len(input.NewPassword) < 8 {
			return nil, fmt.Errorf("%w: a senha deve ter pelo menos 8 caracteres", models.ErrInvalidPassword)
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = string(hash)
	}

	if input.Active != nil {
		user.Active = *input.Active
	}
	user.Name = name
	user.Email = email
	user.Login = login

	if err := u.repos.User.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil

}

func (s *Services) Add(ctx context.Context, input models.CreateUsersInput) (*models.Users, error) {
	name := strings.TrimSpace(input.Name)
	if len(name) < 3 {
		return nil, errors.New("o nome deve ter pelo menos 3 caracteres")
	}

	login := strings.TrimSpace(input.Login)
	if len(login) < 3 {
		return nil, errors.New("o login deve ter pelo menos 3 caracteres")
	}

	email := strings.ToLower(strings.TrimSpace(input.Email))
	if !isValidEmail(email) {
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

func isValidEmail(email string) bool {
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return false
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[0] == "" || !strings.Contains(parts[1], ".") {
		return false
	}

	domainParts := strings.Split(parts[1], ".")
	for _, p := range domainParts {
		if p == "" {
			return false
		}
	}

	return true
}
