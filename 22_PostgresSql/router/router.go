package router

import (
	"database/sql"
	"net/http"
)

type StocksRoutes struct {
}

func NewRouter() *StocksRoutes {
	return &StocksRoutes{}
}

func Router(db *sql.DB) *http.ServeMux {
	router := http.NewServeMux()

	stockRoutes := NewStockRoutes(db)
	stockRoutes.InitalizeStocksRoutes(router)

	return router
}
