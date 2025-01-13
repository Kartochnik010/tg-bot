# Build stage
FROM golang:1.22.10 AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and source files
COPY go.mod go.sum ./
RUN go mod t

COPY . .

# Build the Go binary with static linking (CGO disabled)
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app

# Final stage using distroless (no need for glibc since binary is statically linked)
FROM gcr.io/distroless/static

# Set the working directory
WORKDIR /app

# Copy the binary and config file from the build stage
COPY --from=build /app/main .
COPY --from=build /app/config.json .
COPY --from=build /app/migrations .

# Expose the port that the app will run on
EXPOSE 8080


# Command to run the Go app
ENTRYPOINT ["./main"]
