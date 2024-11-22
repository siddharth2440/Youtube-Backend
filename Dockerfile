# Start from the official Golang image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

RUN go install github.com/air-verse/air@latest
# Copy the go.mod and go.sum files to download dependencies
COPY go.mod ./
COPY go.sum ./

# Download all dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Expose the port that the application will run on
EXPOSE 8000

# Set the entrypoint to Air for live reloading
CMD ["air", "-c", ".air.toml"]