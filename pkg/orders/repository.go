package orders

import "context"

// Repository is the order repo
type Repository interface {
	FindByID(ctx context.Context, id string) (*Order, error)
	Find(ctx context.Context, limit int, offset int) ([]Order, error)
	Create(ctx context.Context, id string, distance int, status string) error
	Update(ctx context.Context, id string) error
}
