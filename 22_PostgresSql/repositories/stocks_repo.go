package repositories

import (
	"database/sql"
	"go/postgresql-demo/models"
	"log"
)

type StocksRepo struct {
	Db *sql.DB
}

func NewSotcksRepo(db *sql.DB) *StocksRepo {
	return &StocksRepo{
		Db: db,
	}
}

func (s *StocksRepo) GetAllStocks() []models.Stock {
	return []models.Stock{}
}

func (s *StocksRepo) GetStockById(id string) models.Stock {
	return models.Stock{}
}

func (s *StocksRepo) CreateNewStock(stock *models.Stock) *uint64 {
	query := `INSERT INTO stocks(name,price,company) VALUES ($1,$2,$3) RETURNING stockId`

	var id uint64
	err := s.Db.QueryRow(query, stock.Name, stock.Price, stock.Company).Scan(&id)

	if err != nil {
		log.Fatal(err)
	}

	return &id
}

func (s *StocksRepo) UpdateStock() {}
func (s *StocksRepo) DeleteStock() {}
