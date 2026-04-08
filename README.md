# Fuel Prices Log CLI

A Go command-line application that fetches fuel prices from the public OK API for a single facility and stores the data in a Postgres database. The app is designed for one-shot execution and is suitable for CronJob usage.

## Features

- Fetches price data from the OK public endpoint.
- Filters to a single station using --facility-number.
- Outputs products only as JSON to stdout.
- Supports configurable log verbosity with --log-level.
- Optionally stores changed prices in PostgreSQL with --save-to-db.
- Inserts new history rows only when a product price changed since the previous entry.

## Prerequisites

- Go 1.22+ (or project-defined Go version once go.mod is added)
- PostgreSQL 14+
- Network access to the OK API

## Installation

1. Clone the repository.
2. Build the binary.
3. Provision PostgreSQL database and table manually.

Example build:

```bash
go build -o fuel-prices ./cmd/fuel-prices
```

## Database Setup (Manual)

The application does not create databases or tables programmatically.

Create database (example):

```sql
CREATE DATABASE fuel_prices;
```

Connect to the database and create table/index:

```sql
CREATE TABLE IF NOT EXISTS fuel_price_history (
    id BIGSERIAL PRIMARY KEY,
    facility_number INTEGER NOT NULL,
    product_name TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL,
    last_updated_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_fuel_price_history_lookup
    ON fuel_price_history (facility_number, product_name, last_updated_time DESC);
```

## Configuration

### CLI Flags

- --facility-number (required): facility number to fetch prices for (example: 507)
- --log-level (optional): debug, info, warn, error (default: info)
- --save-to-db (optional): when set, persist changed prices to PostgreSQL

### Environment Variables

These are required only when --save-to-db is enabled:

- DB_HOST
- DB_PORT
- DB_NAME
- DB_USER
- DB_PASSWORD
- DB_SSLMODE

## Usage

Fetch and print products JSON only:

```bash
./fuel-prices --facility-number 507
```

Enable debug logs:

```bash
./fuel-prices --facility-number 507 --log-level debug
```

Fetch and save changed prices to database:

```bash
DB_HOST=localhost \
DB_PORT=5432 \
DB_NAME=fuel_prices \
DB_USER=postgres \
DB_PASSWORD=postgres \
DB_SSLMODE=disable \
./fuel-prices --facility-number 507 --save-to-db
```

## Output

On success, stdout contains JSON with product entries from the selected facility.

Logs are written to stderr.

## Exit Behavior

The application runs one fetch cycle and exits.

Recommended for scheduling with CronJob/Kubernetes CronJob.

## Project Structure

Planned structure:

- cmd/fuel-prices/main.go
- internal/api/client.go
- internal/config/config.go
- internal/model/types.go
- internal/output/json.go
- internal/store/postgres.go
- docs/api-spec.md
- plan.md

## Development

Run formatting, tests, and build:

```bash
gofmt -w .
go test ./...
go build ./cmd/fuel-prices
```
