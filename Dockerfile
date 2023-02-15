# Use an official Golang runtime as a parent image
FROM golang:1.17.5-alpine3.15 AS build-env

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Build the executable
RUN go build -o client ./cmd/client/main.go

# Use a smaller base image for the final container
FROM alpine:3.15

# Copy the built executable from the previous stage
COPY --from=build-env /app/client /usr/local/bin/client

# Set the command to run when the container starts
CMD ["client"]
