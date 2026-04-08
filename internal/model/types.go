package model

import "time"

type FuelPricesResponse struct {
	Items []Station `json:"items"`
}

type Station struct {
	FacilityNumber  int       `json:"facility_number"`
	LastUpdatedTime time.Time `json:"last_updated_time"`
	Prices          []Price   `json:"prices"`
}

type Price struct {
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
}
