package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchStationSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"items":[{"facility_number":507,"last_updated_time":"2026-04-08T17:33:40Z","prices":[{"product_name":"Blyfri 95","price":16.09}]}]}`))
	}))
	defer ts.Close()

	client := NewClient(ts.URL)
	station, err := client.FetchStation(context.Background(), 507)
	if err != nil {
		t.Fatalf("FetchStation returned error: %v", err)
	}
	if station.FacilityNumber != 507 {
		t.Fatalf("expected facility 507, got %d", station.FacilityNumber)
	}
	if len(station.Prices) != 1 {
		t.Fatalf("expected 1 price, got %d", len(station.Prices))
	}
}

func TestFetchStationNotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"items":[{"facility_number":1,"last_updated_time":"2026-04-08T17:33:40Z","prices":[{"product_name":"Blyfri 95","price":16.09}]}]}`))
	}))
	defer ts.Close()

	client := NewClient(ts.URL)
	_, err := client.FetchStation(context.Background(), 507)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, ErrStationNotFound) {
		t.Fatalf("expected ErrStationNotFound, got %v", err)
	}
}

func TestFetchStationNon200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad gateway", http.StatusBadGateway)
	}))
	defer ts.Close()

	client := NewClient(ts.URL)
	_, err := client.FetchStation(context.Background(), 507)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestFindStationByFacility(t *testing.T) {
	itemsJSON := []byte(`{"items":[{"facility_number":11,"last_updated_time":"2026-04-08T17:33:40Z","prices":[{"product_name":"A","price":1.23}]},{"facility_number":22,"last_updated_time":"2026-04-08T17:33:40Z","prices":[{"product_name":"B","price":2.34}]}]}`)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(itemsJSON)
	}))
	defer ts.Close()

	client := NewClient(ts.URL)
	station, err := client.FetchStation(context.Background(), 22)
	if err != nil {
		t.Fatalf("FetchStation returned error: %v", err)
	}
	if station.FacilityNumber != 22 {
		t.Fatalf("expected facility 22, got %d", station.FacilityNumber)
	}
}
