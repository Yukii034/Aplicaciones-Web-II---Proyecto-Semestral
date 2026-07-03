# ---- Etapa 1: builder ----
FROM golang:1.26-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/proyecto-semestral ./cmd/api

# ---- Etapa 2: runner (imagen final minima) ----
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
RUN adduser -D -u 10001 appuser
WORKDIR /app
COPY --from=builder /bin/proyecto-semestral /app/proyecto-semestral
USER appuser
EXPOSE 8080
ENTRYPOINT ["/app/proyecto-semestral"]