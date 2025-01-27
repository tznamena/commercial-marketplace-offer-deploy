# Use the offical golang image to create a binary.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.19-buster as builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.

COPY go.mod go.sum ./

RUN go mod download
COPY *.go ./

# Copy local code to the container image.
COPY . ./

RUN go build -o ./bin/ ./cmd/operator
RUN go build -o ./bin/ ./cmd/apiserver
COPY ./build/scripts/entrypoint.sh ./bin


# Use the official Debian slim image for a lean production container.
# https://hub.docker.com/_/debian
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM debian:buster-slim

USER root

RUN apt-get clean && apt-get update
RUN apt install curl -y
RUN curl -sL https://aka.ms/InstallAzureCLIDeb | bash

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/bin/operator /operator
COPY --from=builder /app/bin/apiserver /apiserver
COPY --from=builder /app/bin/entrypoint.sh /entrypoint.sh

# Run the web service on container startup.
CMD ["/entrypoint.sh"]
