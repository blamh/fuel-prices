package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"fuel-prices/internal/model"
)

const DefaultBaseURL = "https://mobility-prices.ok.dk/api/v1/fuel-prices"

var ErrStationNotFound = errors.New("facility number not found")

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimSpace(baseURL),
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

func (c *Client) FetchStation(ctx context.Context, facilityNumber int) (model.Station, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL, nil)
	if err != nil {
		return model.Station{}, fmt.Errorf("build request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return model.Station{}, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Station{}, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	var payload model.FuelPricesResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return model.Station{}, fmt.Errorf("decode response: %w", err)
	}

	station, err := FindStationByFacility(payload.Items, facilityNumber)
	if err != nil {
		return model.Station{}, err
	}

	if len(station.Prices) == 0 {
		return model.Station{}, fmt.Errorf("facility %d has no prices", facilityNumber)
	}

	return station, nil
}

func FindStationByFacility(items []model.Station, facilityNumber int) (model.Station, error) {
	for _, item := range items {
		if item.FacilityNumber == facilityNumber {
			return item, nil
		}
	}
	return model.Station{}, fmt.Errorf("%w: %d", ErrStationNotFound, facilityNumber)
}
