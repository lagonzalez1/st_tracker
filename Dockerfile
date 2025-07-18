# Dockerfile.dev
FROM golang:1.23-alpine

# (Optional) if you pull deps from private repos:
RUN apk add --no-cache git

RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Cache your module dependencies
COPY go.mod go.sum ./


# Copy the rest of your code
COPY . .

# Expose your app’s port
EXPOSE 3333

# Run your app on container start
CMD ["sh", "-c", "set -o allexport && . .env.development && go run app/main.go"]