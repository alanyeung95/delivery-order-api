package orders

import (
	"net/http"

	"github.com/go-chi/chi"
)

// NewHandler return handler that serves the order service
func NewHandler(srv Service) http.Handler {
	//h := handlers{srv}
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	return r
}
