package orders

import (
	"context"
	"encoding/json"
	er "errors"
	"net/http"
	"os"

	uuid "github.com/satori/go.uuid"
)

// mockgen -destination=../mocks/mock_orders/mock_service.go -package=mock_orders github.com/alanyeung95/delivery-order-api/pkg/orders Service

// Service interface
type Service interface {
	PlaceOrder(ctx context.Context, distance int) (*Order, error)
	ListOrder(ctx context.Context, limit int, offset int) ([]Order, error)
	GetDistance(ctx context.Context, req PlaceOrderRequest) (int, error)
	TakeOrder(ctx context.Context, id string) error
	GetOrderById(ctx context.Context, id string) (*Order, error)
}

type service struct {
	repository Repository
}

// NewService start the new service
func NewService(repository Repository) Service {
	return &service{repository}
}

func (s *service) PlaceOrder(ctx context.Context, distance int) (*Order, error) {
	id := uuid.NewV4().String()

	err := s.repository.Create(ctx, id, distance, OrderStatusUnassigned.String())
	if err != nil {
		return nil, err
	}

	order, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *service) GetDistance(ctx context.Context, req PlaceOrderRequest) (int, error) {
	client := http.Client{
		// Set timeout if necessary
		//Timeout: time.Duration(1000) * time.Millisecond,
	}

	url := "https://maps.googleapis.com/maps/api/distancematrix/json"
	googleDistanceReq, _ := http.NewRequest("POST", url, nil)

	q := googleDistanceReq.URL.Query()
	q.Add("origins", req.Origin[0]+","+req.Origin[1])
	q.Add("destinations", req.Destination[0]+","+req.Destination[1])
	q.Add("key", os.Getenv("GOOGLE_MAP_API_KEY"))

	googleDistanceReq.URL.RawQuery = q.Encode()

	resp, err := client.Do(googleDistanceReq)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	res := Response{}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return -1, err
	}

	if res.Status == Denied {
		return -1, nil // TODO
	}

	distance := res.Rows[0].Elements[0].Distance.Value
	if distance == nil || *distance < 0 {
		return -1, er.New("google api cannot return distance value")
	}

	return *distance, nil
}

func (s *service) ListOrder(ctx context.Context, limit int, offset int) ([]Order, error) {
	orders, err := s.repository.Find(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *service) TakeOrder(ctx context.Context, id string) error {
	return s.repository.Update(ctx, id)
}

func (s *service) GetOrderById(ctx context.Context, id string) (*Order, error) {
	return s.repository.FindByID(ctx, id)
}
