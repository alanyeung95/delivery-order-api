package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alanyeung95/delivery-order-api/pkg/mysql"
	"github.com/alanyeung95/delivery-order-api/pkg/orders"

	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_ADDRESSES"), os.Getenv("MYSQL_DATABASE")))

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

	addr := fmt.Sprintf(":%s", os.Getenv("API_PORT"))
	fmt.Println("Service is running on " + addr)
	http.ListenAndServe(addr, r)

}
