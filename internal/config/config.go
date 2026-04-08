package config

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"strings"
)

type Config struct {
	FacilityNumber int
	LogLevel       slog.Level
	SaveToDB       bool
	DB             DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

func Parse(args []string, getenv func(string) string) (Config, error) {
	fs := flag.NewFlagSet("fuel-prices", flag.ContinueOnError)
	fs.SetOutput(ioDiscard{})

	facilityNumber := fs.Int("facility-number", 0, "facility number to fetch")
	logLevel := fs.String("log-level", "info", "log level: debug, info, warn, error")
	saveToDB := fs.Bool("save-to-db", false, "save changed prices to postgres")

	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}

	if *facilityNumber <= 0 {
		return Config{}, errors.New("--facility-number must be a positive integer")
	}

	parsedLevel, err := ParseLogLevel(*logLevel)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		FacilityNumber: *facilityNumber,
		LogLevel:       parsedLevel,
		SaveToDB:       *saveToDB,
	}

	if cfg.SaveToDB {
		dbCfg, err := parseDBConfig(getenv)
		if err != nil {
			return Config{}, err
		}
		cfg.DB = dbCfg
	}

	return cfg, nil
}

func ParseLogLevel(v string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("invalid --log-level %q (allowed: debug, info, warn, error)", v)
	}
}

func parseDBConfig(getenv func(string) string) (DBConfig, error) {
	cfg := DBConfig{
		Host:     strings.TrimSpace(getenv("DB_HOST")),
		Port:     strings.TrimSpace(getenv("DB_PORT")),
		Name:     strings.TrimSpace(getenv("DB_NAME")),
		User:     strings.TrimSpace(getenv("DB_USER")),
		Password: strings.TrimSpace(getenv("DB_PASSWORD")),
		SSLMode:  strings.TrimSpace(getenv("DB_SSLMODE")),
	}

	missing := make([]string, 0, 6)
	if cfg.Host == "" {
		missing = append(missing, "DB_HOST")
	}
	if cfg.Port == "" {
		missing = append(missing, "DB_PORT")
	}
	if cfg.Name == "" {
		missing = append(missing, "DB_NAME")
	}
	if cfg.User == "" {
		missing = append(missing, "DB_USER")
	}
	if cfg.Password == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if cfg.SSLMode == "" {
		missing = append(missing, "DB_SSLMODE")
	}

	if len(missing) > 0 {
		return DBConfig{}, fmt.Errorf("--save-to-db requires env vars: %s", strings.Join(missing, ", "))
	}

	return cfg, nil
}

type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) {
	return len(p), nil
}
