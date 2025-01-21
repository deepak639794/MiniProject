# 1. Use an official Golang base image
FROM golang:1.23.3

# 2. Set the working directory in the container
WORKDIR /app

# 3. Copy go.mod and go.sum (caches dependencies)
COPY go.mod go.sum ./

# 4. Download Go module dependencies
RUN go mod download

# 5. Copy the rest of your application code
COPY . .

# 6. Build your Go application
# RUN go build -o main .
RUN go build -o main . && chmod +x main
# 7. List the contents of the /app directory to ensure 'main' binary is there
RUN ls -l /app

# 8. Command to run the app when the container starts
CMD ["main.go"]