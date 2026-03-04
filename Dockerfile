# Dockerfile.dev
FROM golang:1.23-alpine

RUN apk add --no-cache git ca-certificates postgresql-client

WORKDIR /app

# Cache dependencies as a separate layer
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your code
COPY . .

EXPOSE 3333

CMD ["sh", "-c", "set -o allexport && . .env.development && go run app/main.go"]