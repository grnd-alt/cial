FROM golang:1.23-alpine AS builder


WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /goapp

FROM alpine:3.18

# Copy the binary from the builder stage.
COPY --from=builder /goapp /goapp
COPY --from=builder /app/db/sql/migrations /app/db/sql/migrations

# Copy the .env file to the working directory.
# COPY .env /app/.env

# Set the working directory.
WORKDIR /app

# Expose the port the app runs on.
EXPOSE 8000

# Run the Go app.
CMD ["/goapp"]
