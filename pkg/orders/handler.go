package orders

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alanyeung95/GoProjectDemo/pkg/errors"
	"github.com/go-chi/chi"
	kithttp "github.com/go-kit/kit/transport/http"
)

// NewHandler return handler that serves the order service
func NewHandler(svc Service) http.Handler {
	h := handlers{svc}
	r := chi.NewRouter()
	r.Post("/", h.handlePlaceOrder)
	r.Get("/", h.handleListOrders)
	r.Patch("/{id}", h.handleTakeOrder)
	return r
}

type handlers struct {
	svc Service
}

func (h *handlers) handlePlaceOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req PlaceOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		kithttp.DefaultErrorEncoder(ctx, errors.NewBadRequest(err), w)
	}

	distance, err := h.svc.GetDistance(ctx, req)
	if err != nil {
		kithttp.DefaultErrorEncoder(ctx, errors.NewBadRequest(err), w)
	}

	fmt.Println(distance)
	newOrder, err := h.svc.PlaceOrder(ctx, distance)

	kithttp.EncodeJSONResponse(ctx, w, newOrder)
}

func (h *handlers) handleListOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orderList, err := h.svc.ListOrder(ctx)
	if err != nil {
		kithttp.DefaultErrorEncoder(ctx, err, w)
		return
	}
	kithttp.EncodeJSONResponse(ctx, w, orderList)
}

func (h *handlers) handleTakeOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	success, err := h.svc.TakeOrder(ctx, id)
	if err != nil {
		kithttp.DefaultErrorEncoder(ctx, err, w)
		return
	}
	kithttp.EncodeJSONResponse(ctx, w, success)

}
