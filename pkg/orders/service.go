package orders

import (
	"context"
	"database/sql"
	"encoding/json"
	er "errors"
	"net/http"
	"os"

	uuid "github.com/satori/go.uuid"
)

// Service interface
type Service interface {
	PlaceOrder(ctx context.Context, distance int) (*Order, error)
	ListOrder(ctx context.Context) ([]Order, error)
	GetDistance(ctx context.Context, req PlaceOrderRequest) (int, error)
	TakeOrder(ctx context.Context, id string) error
	GetOrderById(ctx context.Context, id string) (*Order, error)
}

type service struct {
	repository *sql.DB
}

// NewService start the new service
func NewService(repository *sql.DB) Service {
	return &service{repository}
}

func (s *service) PlaceOrder(ctx context.Context, distance int) (*Order, error) {
	id := uuid.NewV4().String()

	_, err := s.repository.Exec("INSERT INTO delivery_order(id,distance,status) values(?,?,?)", id, distance, OrderStatusUnassigned)
	if err != nil {
		return nil, err
	}

	var order Order
	// TODO: db naming
	err = s.repository.QueryRow("SELECT id, distance, status FROM delivery_order where id = ?", id).Scan(&order.ID, &order.Distance, &order.Status)
	if err != nil {
		return nil, err
	}

	return &order, nil
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

func (s *service) ListOrder(ctx context.Context) ([]Order, error) {
	// Execute the query
	results, err := s.repository.Query("SELECT * FROM delivery_order")
	if err != nil {
		return nil, err
	}
	var orderList []Order

	for results.Next() {
		var order Order
		// for each row, scan the result into our tag composite object
		err = results.Scan(&order.ID, &order.Distance, &order.Status)
		if err != nil {
			return nil, err
		}
		orderList = append(orderList, order)
	}

	if orderList == nil {
		return []Order{}, nil
	}

	return orderList, nil
}

func (s *service) TakeOrder(ctx context.Context, id string) error {
	result, err := s.repository.Exec("Update delivery_order set status=? where id=? and status=?", "TAKEN", id, "UNASSIGNED")
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if count != 1 || err != nil {
		return err
	}

	return nil
}

func (s *service) GetOrderById(ctx context.Context, id string) (*Order, error) {
	var order Order
	// TODO: db naming
	err := s.repository.QueryRow("SELECT id, distance, status FROM delivery_order where id = ?", id).Scan(&order.ID, &order.Distance, &order.Status)
	if err != nil {
		return nil, err
	}

	return &order, nil
}
