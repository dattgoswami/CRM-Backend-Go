# Starting from golang base image
FROM golang:1.19

# Location of our directory where we'd execute commands
WORKDIR /go/src/github.com/CRM-Backend-Go

# Copy the source from this project to the filesystem of the container
COPY . .

# Download all the dependencies
RUN go get -u github.com/gorilla/mux
RUN go get -u github.com/lib/pq

# Build the binary exe for the Go app
RUN go run main.go 