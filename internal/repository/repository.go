package repository

import "github.com/medivh13/koalatest/internal/models"

type Repository interface {
	Register(data *models.Customers) error
	GetToken(phoneOrEmail, pass string) ([]*models.GetTokens, error)
}
