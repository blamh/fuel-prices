package config

import (
	"log/slog"
	"strings"
	"testing"
)

func TestParseWithoutSaveToDBDoesNotRequireEnv(t *testing.T) {
	cfg, err := Parse([]string{"--facility-number", "507"}, func(string) string { return "" })
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if cfg.FacilityNumber != 507 {
		t.Fatalf("expected facility number 507, got %d", cfg.FacilityNumber)
	}
	if cfg.SaveToDB {
		t.Fatalf("expected SaveToDB false")
	}
	if cfg.LogLevel != slog.LevelInfo {
		t.Fatalf("expected default log level info, got %v", cfg.LogLevel)
	}
}

func TestParseSaveToDBRequiresEnv(t *testing.T) {
	_, err := Parse([]string{"--facility-number", "507", "--save-to-db"}, func(string) string { return "" })
	if err == nil {
		t.Fatalf("expected error for missing env vars")
	}
	if !strings.Contains(err.Error(), "DB_HOST") {
		t.Fatalf("expected missing env var list in error, got %v", err)
	}
}

func TestParseSaveToDBWithEnvSucceeds(t *testing.T) {
	env := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"DB_NAME":     "fuel_prices",
		"DB_USER":     "postgres",
		"DB_PASSWORD": "secret",
		"DB_SSLMODE":  "disable",
	}

	cfg, err := Parse([]string{"--facility-number", "507", "--save-to-db"}, func(k string) string { return env[k] })
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if !cfg.SaveToDB {
		t.Fatalf("expected SaveToDB true")
	}
	if cfg.DB.Name != "fuel_prices" {
		t.Fatalf("expected DB name fuel_prices, got %s", cfg.DB.Name)
	}
}

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    slog.Level
		wantErr bool
	}{
		{name: "debug", value: "debug", want: slog.LevelDebug},
		{name: "info", value: "info", want: slog.LevelInfo},
		{name: "warn", value: "warn", want: slog.LevelWarn},
		{name: "error", value: "error", want: slog.LevelError},
		{name: "invalid", value: "trace", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLogLevel(tt.value)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}
