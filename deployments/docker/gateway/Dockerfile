# ---------- build ----------
FROM golang:1.27-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/gateway ./cmd/gateway

# ---------- runtime ----------
FROM alpine:3.22
WORKDIR /app

COPY --from=build /out/gateway ./gateway
COPY internal/gateway/migrations ./internal/gateway/migrations

EXPOSE 50050
ENTRYPOINT ["./gateway"]