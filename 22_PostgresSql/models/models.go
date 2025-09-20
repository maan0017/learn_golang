package models

type Stock struct {
	StockId uint64 `json:"stockId"`
	Name    string `json:"name"`
	Price   uint64 `json:"price"`
	Company string `json:"company"`
}
