# ---------- build ----------
FROM golang:1.27-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/account ./cmd/account

# ---------- runtime ----------
FROM alpine:3.22
WORKDIR /app

COPY --from=build /out/account ./account
COPY internal/account/migrations ./internal/account/migrations

EXPOSE 50052
ENTRYPOINT ["./account"]