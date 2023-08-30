# Get the primary IP address of the host machine
$HOST_IP = "192.168.3.85"

# Set the IP address in the environment variable
$env:VIDEO_SERVER_URL = "http://{0}:8080" -f $HOST_IP

# Launch your services with Docker Compose
docker-compose up -d
