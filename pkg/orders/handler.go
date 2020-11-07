package orders

import (
	"database/sql"
	"encoding/json"
	er "errors"
	"net/http"
	"strconv"

	"github.com/alanyeung95/delivery-order-api/pkg/errors"
	"github.com/go-chi/chi"
	kithttp "github.com/go-kit/kit/transport/http"

	"gopkg.in/go-playground/validator.v9"
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
		kithttp.DefaultErrorEncoder(ctx, errors.NewBadRequestError(err), w)
		return
	}

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		kithttp.DefaultErrorEncoder(ctx, errors.NewBadRequestError(err), w)
		return
	}

	distance, err := h.svc.GetDistance(ctx, req)
	if err != nil {
		kithttp.DefaultErrorEncoder(ctx, errors.NewServerError(er.New(err.Error())), w)
		return
	}

	newOrder, err := h.svc.PlaceOrder(ctx, distance)
	if err != nil {
		kithttp.DefaultErrorEncoder(ctx, errors.NewServerError(er.New(err.Error())), w)
		return
	}

	kithttp.EncodeJSONResponse(ctx, w, newOrder)
}

func (h *handlers) handleListOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqestMap := r.URL.Query()

	page, err := strconv.Atoi(reqestMap.Get("page"))
	if err != nil {
		kithttp.DefaultErrorEncoder(ctx, errors.NewBadRequestError(err), w)
		return
	} else if page <= 0 {
		kithttp.DefaultErrorEncoder(ctx, errors.NewBadRequestError(er.New("page number must starts with 1")), w)
		return
	}

	limit, err := strconv.Atoi(reqestMap.Get("limit"))
	if err != nil {
		kithttp.DefaultErrorEncoder(ctx, errors.NewBadRequestError(err), w)
		return
	} else if limit <= 0 {
		kithttp.DefaultErrorEncoder(ctx, errors.NewBadRequestError(er.New("limte must larger than 0")), w)
		return
	}

	orderList, err := h.svc.ListOrder(ctx, limit, limit*(page-1))
	if err != nil {
		kithttp.DefaultErrorEncoder(ctx, err, w)
		return
	}
	kithttp.EncodeJSONResponse(ctx, w, orderList)
}

func (h *handlers) handleTakeOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	targetOrder, err := h.svc.GetOrderById(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			kithttp.DefaultErrorEncoder(ctx, errors.NewResourceNotFoundError(er.New("Cannot find order: "+id)), w)
		} else {
			kithttp.DefaultErrorEncoder(ctx, errors.NewServerError(er.New(err.Error())), w)
		}
		return
	}

	if targetOrder.Status == OrderStatusTaken {
		kithttp.EncodeJSONResponse(ctx, w, TakeOrderResponse{Status: OrderStatusTaken})
		return
	}

	err = h.svc.TakeOrder(ctx, id)
	if err != nil {
		kithttp.EncodeJSONResponse(ctx, w, TakeOrderResponse{Status: OrderStatusTaken})
		return
	}
	kithttp.EncodeJSONResponse(ctx, w, TakeOrderResponse{Status: OrderStatusSuccess})
}
