# syntax=docker/dockerfile:1

# Alpine is chosen for its small footprint
# compared to Ubuntu
FROM golang:1.18-alpine AS builder

#RUN apk add --no-cache git

# Set the Current Working Directory inside the container
#WORKDIR .

# We want to populate the module cache based on the go.{mod,sum} files.
#COPY go.mod .
#COPY go.sum .

#RUN go mod download

ADD *.go .

# Unit tests
#RUN go test -v

# Build the Go app
RUN go build main.go
#
#FROM scratch
#COPY --from=builder /main .

# Run the binary program produced by `go install`
CMD go run main.go

# This container exposes port 8080 to the outside world
EXPOSE 8080