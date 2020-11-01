package orders

import (
	"context"
	"database/sql"
	"net/http"
)

// Service interface
type Service interface {
	CreateOrder(ctx context.Context, r *http.Request, order *Order) (*Order, error)
}

type service struct {
	repository *sql.DB
}

// NewService start the new service
func NewService(repository *sql.DB) Service {
	return &service{repository}
}

func (s *service) CreateOrder(ctx context.Context, r *http.Request, account *Order) (*Order, error) {
	return nil, nil
}
