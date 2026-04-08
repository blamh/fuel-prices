package store

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strconv"
	"time"

	_ "github.com/lib/pq"

	"fuel-prices/internal/config"
	"fuel-prices/internal/model"
)

type PriceStore struct {
	db *sql.DB
}

func NewPriceStore(db *sql.DB) *PriceStore {
	return &PriceStore{db: db}
}

func OpenPostgres(ctx context.Context, cfg config.DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.User,
		cfg.Password,
		cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return db, nil
}

func (s *PriceStore) SaveChangedPrices(ctx context.Context, facilityNumber int, lastUpdated time.Time, prices []model.Price) (int, error) {
	inserted := 0
	for _, p := range prices {
		changed, err := s.hasPriceChanged(ctx, facilityNumber, p)
		if err != nil {
			return inserted, err
		}
		if !changed {
			continue
		}

		if err := s.insertPrice(ctx, facilityNumber, p, lastUpdated); err != nil {
			return inserted, err
		}
		inserted++
	}
	return inserted, nil
}

func (s *PriceStore) hasPriceChanged(ctx context.Context, facilityNumber int, p model.Price) (bool, error) {
	const q = `
SELECT price
FROM fuel_price_history
WHERE facility_number = $1 AND product_name = $2
ORDER BY last_updated_time DESC, id DESC
LIMIT 1`

	var latest string
	err := s.db.QueryRowContext(ctx, q, facilityNumber, p.ProductName).Scan(&latest)
	if err == sql.ErrNoRows {
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("query latest price for %q: %w", p.ProductName, err)
	}

	latestValue, err := strconv.ParseFloat(latest, 64)
	if err != nil {
		return false, fmt.Errorf("parse latest price for %q: %w", p.ProductName, err)
	}

	return priceToCents(latestValue) != priceToCents(p.Price), nil
}

func (s *PriceStore) insertPrice(ctx context.Context, facilityNumber int, p model.Price, lastUpdated time.Time) error {
	const q = `
INSERT INTO fuel_price_history (facility_number, product_name, price, last_updated_time)
VALUES ($1, $2, $3, $4)`

	priceValue := fmt.Sprintf("%.2f", p.Price)
	_, err := s.db.ExecContext(ctx, q, facilityNumber, p.ProductName, priceValue, lastUpdated)
	if err != nil {
		return fmt.Errorf("insert price for %q: %w", p.ProductName, err)
	}

	return nil
}

func priceToCents(v float64) int64 {
	return int64(math.Round(v * 100))
}
