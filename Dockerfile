# Start from the Go base image
FROM golang:latest

# Setup repo
WORKDIR /app
COPY . .

# Install ffmpeg
RUN apt-get update && \
    apt-get install -y ffmpeg && \
    rm -rf /var/lib/apt/lists/*
    
# Install dependencies
RUN go mod download

# Build the Go server
RUN go build -o main .

# Expose port 8080 (or whatever port your app runs on)
EXPOSE 8080

# Run the binary
CMD ["./main"]