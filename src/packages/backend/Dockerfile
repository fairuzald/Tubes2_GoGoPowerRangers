# Use official Golang image as base image for backend
FROM golang:alpine AS backend

# Set working directory for backend
WORKDIR /app/packages/backend

# Copy go.mod and go.sum
COPY packages/backend/go.mod packages/backend/go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the rest of the backend application
COPY packages/backend/ .

# Build the backend
RUN go build -o main .

# Expose the port used by the backend application
EXPOSE 8080

# Start the backend server
CMD ["./main"]
