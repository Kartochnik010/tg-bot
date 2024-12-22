# Build stage
FROM golang:1.22.2 AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and source files
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the Go binary
RUN GOOS=linux go build -o main ./cmd/app

# Final stage
FROM debian:bullseye-slim

# Set the working directory
WORKDIR /app

# Copy the binary and config file from the build stage
COPY --from=build /app/main .
COPY --from=build /app/config.json .

# Expose the port that the app will run on
EXPOSE 8080

# Command to run the Go app
ENTRYPOINT ["./main"]