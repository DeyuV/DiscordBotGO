# Use an official Golang runtime as the base image
FROM golang:1.18 AS build

# Install the tzdata package to include timezone data in the container
RUN apt-get update && apt-get install -y --no-install-recommends tzdata && rm -rf /var/lib/apt/lists/*

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Set the working directory inside the container
WORKDIR /app/cmd/monolith

# Build the Go application inside the container
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Use a lightweight Alpine Linux as the base image for the final container
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built Go executable from the previous stage
COPY --from=build /app/cmd/monolith/app .

# Expose the port on which the Go application listens
EXPOSE 8080

# Set the command to run the Go application when the container starts
CMD ["./app"]