# Use the official Golang image as the base image
FROM golang:1.23

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download


# Build the Go app
RUN go build -o main ./cmd/main.go
RUN chmod +x ./main


# Command to run the executable
CMD ["./main"]
