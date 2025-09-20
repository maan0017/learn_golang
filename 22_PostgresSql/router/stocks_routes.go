package router

import (
	"database/sql"
	"go/postgresql-demo/handlers"
	"net/http"
)

type StockRoutes struct {
	Handlers *handlers.StocksHandlers
}

func NewStockRoutes(db *sql.DB) *StockRoutes {
	stockHandlers := handlers.NewStocksHandlers(db)

	return &StockRoutes{
		Handlers: stockHandlers,
	}
}

func (r *StockRoutes) InitalizeStocksRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /api/sotcks", r.Handlers.GetAllStocksHandler)
	router.HandleFunc("GET /api/sotck/{id}", r.Handlers.GetStockByIdHander)
	router.HandleFunc("PUT /api/sotck/{id}", r.Handlers.UpdateStock)
	router.HandleFunc("DELETE /api/sotck/{id}", r.Handlers.DeleteStock)
	router.HandleFunc("POST /api/new-stock", r.Handlers.CreateNewStockHander)
}
