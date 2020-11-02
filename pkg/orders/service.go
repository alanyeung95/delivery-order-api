package orders

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	uuid "github.com/satori/go.uuid"
)

// Service interface
type Service interface {
	PlaceOrder(ctx context.Context, distance int) (*Order, error)
	ListOrder(ctx context.Context) ([]Order, error)
	GetDistance(ctx context.Context, req PlaceOrderRequest) (int, error)
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

	_, err := s.repository.Exec("INSERT INTO delivery_order(id,distance,status) values(?,?,?)", id, distance, Unassigned)
	if err != nil {
		return nil, err
	}

	var order Order
	// TODO: db naming
	err = s.repository.QueryRow("SELECT id, distance, status FROM delivery_order where id = ?", id).Scan(&order.ID, &order.Distance, &order.Status)
	if err != nil {
		fmt.Println(err)
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

	return res.Rows[0].Elements[0].Distance.Value, nil
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
