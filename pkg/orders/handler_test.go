package orders_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	er "errors"

	"github.com/alanyeung95/delivery-order-api/pkg/mocks/mock_orders"
	"github.com/alanyeung95/delivery-order-api/pkg/orders"
	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var TestLoginHandler = Describe("Order handler", func() {
	var (
		mockCtrl    *gomock.Controller
		mockService *mock_orders.MockService
		handler     http.Handler
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockService = mock_orders.NewMockService(mockCtrl)
		handler = orders.NewHandler(mockService)
	})

	var _ = Describe("Place Order", func() {
		Context("Valid request", func() {
			It("should return valid response", func() {
				expectedResponse := &orders.Order{}
				mockService.EXPECT().PlaceOrder(gomock.Any(), gomock.Any()).Return(expectedResponse, nil)
				mockService.EXPECT().GetDistance(gomock.Any(), gomock.Any()).Return(1, nil)

				var jsonStr = []byte(`{
					"origin": ["22.316397", "114.264144" ],
					"destination": ["22.307588", "114.260881"]
				}`)
				r, err := http.NewRequest("POST", "/", bytes.NewBuffer(jsonStr))
				Expect(err).ShouldNot((HaveOccurred()))
				r.Header.Set("Content-Type", "application/json")

				w := httptest.NewRecorder()
				handler.ServeHTTP(w, r)
				resp := w.Result()
				Expect(resp.StatusCode).Should(Equal(http.StatusOK))
			})
		})
		Context("Invalid request", func() {
			It("should return error when params are missing", func() {
				expectedResponse := &orders.Order{}
				mockService.EXPECT().PlaceOrder(gomock.Any(), gomock.Any()).Return(expectedResponse, nil)
				mockService.EXPECT().GetDistance(gomock.Any(), gomock.Any()).Return(1, nil)

				var jsonStr = []byte(`{
					"origin": ["22.316397", "114.264144" ]
				}`)
				r, err := http.NewRequest("POST", "/", bytes.NewBuffer(jsonStr))
				Expect(err).ShouldNot((HaveOccurred()))
				r.Header.Set("Content-Type", "application/json")

				w := httptest.NewRecorder()
				handler.ServeHTTP(w, r)
				resp := w.Result()
				Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))
			})
			It("should return error when params are invalid (string)", func() {
				mockService.EXPECT().PlaceOrder(gomock.Any(), gomock.Any()).Return(&orders.Order{}, nil)
				mockService.EXPECT().GetDistance(gomock.Any(), gomock.Any()).Return(1, nil)

				var jsonStr = []byte(`{
					"origin": ["a", "b" ],
					"destination": ["a", "b"]
				}`)
				r, err := http.NewRequest("POST", "/", bytes.NewBuffer(jsonStr))
				Expect(err).ShouldNot((HaveOccurred()))
				r.Header.Set("Content-Type", "application/json")

				w := httptest.NewRecorder()
				handler.ServeHTTP(w, r)
				resp := w.Result()
				Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))
			})
		})
	})

	var _ = Describe("List Order", func() {
		Context("Valid request", func() {
			It("should return valid response", func() {
				mockService.EXPECT().ListOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return([]orders.Order{}, nil)

				r, err := http.NewRequest("GET", "/?page=1&limit=1", nil)
				Expect(err).ShouldNot((HaveOccurred()))
				w := httptest.NewRecorder()
				handler.ServeHTTP(w, r)
				resp := w.Result()
				Expect(resp.StatusCode).Should(Equal(http.StatusOK))
			})
		})
		Context("Invalid request", func() {
			It("should return error when page = 0", func() {
				mockService.EXPECT().ListOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return([]orders.Order{}, nil)

				r, err := http.NewRequest("GET", "/?page=0&limit=1", nil)
				Expect(err).ShouldNot((HaveOccurred()))
				w := httptest.NewRecorder()
				handler.ServeHTTP(w, r)
				resp := w.Result()
				Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))
			})
			It("should return error when limit = 0", func() {
				mockService.EXPECT().ListOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return([]orders.Order{}, nil)

				r, err := http.NewRequest("GET", "/?page=1&limit=0", nil)
				Expect(err).ShouldNot((HaveOccurred()))
				w := httptest.NewRecorder()
				handler.ServeHTTP(w, r)
				resp := w.Result()
				Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))
			})
			It("should return error when page and limit are not provided", func() {
				mockService.EXPECT().ListOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return([]orders.Order{}, nil)

				r, err := http.NewRequest("GET", "/", nil)
				Expect(err).ShouldNot((HaveOccurred()))
				w := httptest.NewRecorder()
				handler.ServeHTTP(w, r)
				resp := w.Result()
				Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))
			})
		})
	})

	var _ = Describe("Take Order", func() {
		Context("Valid request", func() {
			It("should return valid response", func() {
				mockService.EXPECT().TakeOrder(gomock.Any(), gomock.Any()).Return(nil)
				mockService.EXPECT().GetOrderById(gomock.Any(), gomock.Any()).Return(&orders.Order{}, nil)

				r, err := http.NewRequest("PATCH", "/testing-id", nil)
				Expect(err).ShouldNot((HaveOccurred()))
				w := httptest.NewRecorder()
				handler.ServeHTTP(w, r)
				resp := w.Result()
				Expect(resp.StatusCode).Should(Equal(http.StatusOK))
			})
			It("should return error when id is not found", func() {
				mockService.EXPECT().TakeOrder(gomock.Any(), gomock.Any()).Return(nil)
				mockService.EXPECT().GetOrderById(gomock.Any(), gomock.Any()).Return(nil, sql.ErrNoRows)

				r, err := http.NewRequest("PATCH", "/testing-id", nil)
				Expect(err).ShouldNot((HaveOccurred()))
				w := httptest.NewRecorder()
				handler.ServeHTTP(w, r)
				resp := w.Result()
				Expect(resp.StatusCode).Should(Equal(http.StatusNotFound))
			})
			It("should return error when order is taken", func() {
				mockService.EXPECT().TakeOrder(gomock.Any(), gomock.Any()).Return(er.New(""))
				mockService.EXPECT().GetOrderById(gomock.Any(), gomock.Any()).Return(&orders.Order{}, nil)

				r, err := http.NewRequest("PATCH", "/testing-id", nil)
				Expect(err).ShouldNot((HaveOccurred()))
				w := httptest.NewRecorder()
				handler.ServeHTTP(w, r)
				resp := w.Result()

				respModel := orders.TakeOrderResponse{}
				err = json.NewDecoder(resp.Body).Decode(&respModel)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(respModel.Status).Should(Equal(orders.OrderStatusTaken))
			})
		})

	})

})
