# Use the official Golang image as the base image
FROM golang:1.20

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main ./cmd/main.go

# Command to run the executable
CMD ["./main"]
