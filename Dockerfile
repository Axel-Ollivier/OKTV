# ---- Build stage ----
FROM golang:1.22-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git ca-certificates && update-ca-certificates

# Cache deps first
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bot .

# ---- Run stage (distroless, non-root) ----
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /bot /bot

# Environment comes from docker run / compose (.env)
USER nonroot:nonroot
ENTRYPOINT ["/bot"]
