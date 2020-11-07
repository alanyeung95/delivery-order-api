package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/alanyeung95/delivery-order-api/pkg/mysql"
	"github.com/alanyeung95/delivery-order-api/pkg/orders"

	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "order_service:password@tcp(mysql.network:3306)/orders")
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	//db.SetMaxOpenConns(10)
	//db.SetMaxIdleConns(10)
	orderRepository, err := mysql.NewOrderRepository(db)
	if err != nil {
		panic(err)
	}
	orderSrv := orders.NewService(orderRepository)

	r := chi.NewRouter()

	// Route - API
	r.Route("/", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("welcome"))
		})
		r.Mount("/orders", orders.NewHandler(orderSrv))
	})

	addr := fmt.Sprintf(":%d", 8080)
	fmt.Println("Service is running on " + addr)
	http.ListenAndServe(addr, r)

}
