package handlers

import (
	"database/sql"
	"encoding/json"
	"go/postgresql-demo/models"
	"go/postgresql-demo/repositories"
	"log"
	"net/http"
)

type response struct {
	Id      uint64 `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

type StocksHandlers struct {
	StocksRepo *repositories.StocksRepo
}

func NewStocksHandlers(db *sql.DB) *StocksHandlers {
	stocksRepo := repositories.NewSotcksRepo(db)
	return &StocksHandlers{
		StocksRepo: stocksRepo,
	}
}

func (h *StocksHandlers) GetAllStocksHandler(w http.ResponseWriter, r *http.Request) {
	stocks := h.StocksRepo.GetAllStocks()

	json.NewEncoder(w).Encode(stocks)
}

func (h *StocksHandlers) GetStockByIdHander(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	stock := h.StocksRepo.GetStockById(id)

	json.NewEncoder(w).Encode(stock)
}

func (h *StocksHandlers) CreateNewStockHander(w http.ResponseWriter, r *http.Request) {
	var stock models.Stock

	if err := json.NewDecoder(r.Body).Decode(&stock); err != nil {
		log.Fatal("Unable to decode the request body: ", err)
		return
	}

	stockId := h.StocksRepo.CreateNewStock(&stock)

	res := response{
		Id:      *stockId,
		Message: "stock created successfully",
	}

	json.NewEncoder(w).Encode(res)
}

func (h *StocksHandlers) UpdateStock(w http.ResponseWriter, r *http.Request) {
	// id := r.PathValue("id")
}

func (h *StocksHandlers) DeleteStock(w http.ResponseWriter, r *http.Request) {}
