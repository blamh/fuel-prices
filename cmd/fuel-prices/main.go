package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"time"

	"fuel-prices/internal/api"
	"fuel-prices/internal/config"
	"fuel-prices/internal/store"
)

func main() {
	logger := newLogger(slog.LevelInfo)

	cfg, err := config.Parse(os.Args[1:], os.Getenv)
	if err != nil {
		logger.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	// Rebuild logger now that the configured level is known.
	logger = newLogger(cfg.LogLevel)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := api.NewClient(api.DefaultBaseURL)
	station, err := client.FetchStation(ctx, cfg.FacilityNumber)
	if err != nil {
		logger.Error("failed to fetch station prices", "error", err, "facility_number", cfg.FacilityNumber)
		os.Exit(1)
	}

	if cfg.SaveToDB {
		db, err := store.OpenPostgres(ctx, cfg.DB)
		if err != nil {
			logger.Error("failed to connect to postgres", "error", err)
			os.Exit(1)
		}
		defer db.Close()

		priceStore := store.NewPriceStore(db)
		inserted, err := priceStore.SaveChangedPrices(ctx, station.FacilityNumber, station.LastUpdatedTime, station.Prices)
		if err != nil {
			logger.Error("failed to persist prices", "error", err)
			os.Exit(1)
		}
		logger.Info("price persistence completed", "inserted_rows", inserted, "facility_number", station.FacilityNumber)
	}

	productsJSON, err := json.Marshal(station.Prices)
	if err != nil {
		logger.Error("failed to serialize products for logging", "error", err, "facility_number", station.FacilityNumber)
		os.Exit(1)
	}
	logger.Info("pulled fuel products for facility", "facility_number", station.FacilityNumber, "products", string(productsJSON))
}

func newLogger(level slog.Level) *slog.Logger {
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	return slog.New(h)
}
