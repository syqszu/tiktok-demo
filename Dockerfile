# Start from the Go base image
FROM golang:latest

# Setup repo
WORKDIR /app
COPY . .

# Install dependencies
RUN go mod download

# Build the Go server
RUN go build -o main .

# Expose port 8080 (or whatever port your app runs on)
EXPOSE 8080

# Wait for MySQL to start
COPY ./scripts/wait-mysql.sh /wait-mysql.sh
RUN chmod +x /wait-mysql.sh

# Run the binary
CMD ["./wait-mysql.sh && ./main"]
