package models

import (
	"errors"
)

var ErrUserNotFound = errors.New("usuário não encontrado")
var ErrValidation = errors.New("erro de validação")
var ErrInvalidPassword = errors.New("senha inválida")
