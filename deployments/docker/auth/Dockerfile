# ---------- build ----------
FROM golang:1.27-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/auth ./cmd/auth

# ---------- runtime ----------
FROM alpine:3.22
WORKDIR /app

COPY --from=build /out/auth ./auth
COPY internal/auth/migrations ./internal/auth/migrations

EXPOSE 50051
ENTRYPOINT ["./auth"]