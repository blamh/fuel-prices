# Project Guidelines

## Architecture

One-shot CLI designed for CronJob use: parse config → fetch API → (optionally) persist to Postgres → write JSON to stdout → exit.

| Package | Purpose |
|---|---|
| `cmd/fuel-prices` | Entrypoint; orchestrates all phases |
| `internal/api` | Fetches from OK public API; filters one station by `facility_number` |
| `internal/config` | CLI flag and env var parsing; injects `getenv` for testability |
| `internal/model` | Shared data types (`FuelPricesResponse`, `Station`, `Price`) |
| `internal/output` | Writes products-only JSON to `stdout` |
| `internal/store` | Postgres: change-aware inserts only (compares against latest stored row) |

See [README.md](../README.md) for full usage, env vars, and manual DB setup instructions.

## Build and Test

```bash
go build -o fuel-prices ./cmd/fuel-prices
go test ./...
gofmt -w .
go mod tidy
```

Container build and push is handled by the [GitHub Actions workflow](workflows/build-and-push.yml) — targets `linux/amd64` via `ghcr.io`.

## Security and Dependencies

- Prefer Go standard library first.
- Only add external dependencies that are trusted, actively maintained, and necessary.
- Implement trivial utility functionality inside the project instead of adding third-party packages.

## Conventions

**Output split:** products JSON → `stdout`; all logs (structured `log/slog`) → `stderr`. Never mix them.

**Dependency policy:** prefer stdlib. Only add external packages that are trusted and actively maintained. Implement trivial helpers in-project rather than adding a dependency.

**Errors:** wrap with `fmt.Errorf("context: %w", err)`. Check with `errors.Is()`. Return errors to callers; only call `os.Exit` in `main`.

**Logging:** use the constructed `logger` instance, never the package-level `slog.*` functions. Level default is `info`.

**Config injection:** `config.Parse(args, getenv)` takes a `func(string) string` for env lookup — use this pattern in new packages to keep them testable without global state.

**Price comparison:** always use `priceToCents()` (converts to `int64`) when comparing floats to avoid precision drift. Store as `NUMERIC(10,2)`.

## Pitfalls

- DB env vars (`DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`, `DB_SSLMODE`) are only validated when `--save-to-db` is set. Do not add unconditional checks.
- The `fuel_price_history` table must be provisioned manually before using `--save-to-db`. The app does not create it. Schema is in [README.md](../README.md#database-setup-manual).
- Global context timeout is 30 seconds (API fetch has its own 20-second timeout). Keep DB operations lightweight.

## Commit and PR title format

- Use semantic format: `type(scope): subject`.
- Allowed `type` examples: `feat`, `fix`, `refactor`, `test`, `docs`, `chore`.
- Keep subject imperative and specific (example: `fix(reports): align report class names with DI keys`).
- PR titles follow the same semantic format.
- After completing any change to files in this repo, always end your response with a suggested `git commit` message in a code block, following the Conventional Commits format (e.g. `fix:`, `feat:`, `chore:`, `docs:`). Use third-person singular present tense (e.g. `fix: Adds readinessProbe tcpSocket`, not `Add` or `Added`). The message must not exceed 50 characters.
