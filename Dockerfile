FROM golang:1.26-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o fuel-prices ./cmd/fuel-prices

FROM alpine:3

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=build /app/fuel-prices .

ENTRYPOINT ["./fuel-prices"]
