package application

import (
	"net/http"

	"github.com/Naveenchand06/go-redis-microservice/handler"
	"github.com/Naveenchand06/go-redis-microservice/repository/order"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/order", a.loadOrderRoutes)

	a.router = router
}

func (a *App) loadOrderRoutes(router chi.Router) {
	orderHandler := &handler.Order{
		Repo: &order.RedisRepo{
			Client: a.rdb,
		},
	}

	router.Post("/", orderHandler.Create)
	router.Get("/", orderHandler.List)
	router.Get("/{id}", orderHandler.GetByID)
	router.Put("/{id}", orderHandler.UpdateByID)
	router.Delete("/{id}", orderHandler.DeleteByID)
}

// {
//     "customer_id": "09480E02-0E91-43B6-9EA2-365D7FEAC015",
//     "line_items": [
//         {
//             "item_id": "99F97E8D-06AC-4074-9F14-A133BABA7A25",
//             "quantity": 6,
//             "price": 6000
//         }
//     ]
// }
