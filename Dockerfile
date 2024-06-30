ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /run-app .

FROM debian:bookworm

# Copy the .env file into the image
COPY .env /usr/src/app/.env

# Copy the binary from the builder stage
COPY --from=builder /run-app /usr/local/bin/

# Ensure the working directory is set to where the .env file is located
WORKDIR /usr/src/app

CMD ["run-app"]
