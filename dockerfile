# Use Golang base image
FROM golang:1.24-alpine

# Install make, git, curl(for sqlc install)
RUN apk add --no-cache make git curl

# Cài đặt goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest && \
    cp /go/bin/goose /usr/local/bin/goose

# Install sqlc
RUN curl -sSfL https://github.com/sqlc-dev/sqlc/releases/download/v1.27.0/sqlc_1.27.0_linux_amd64.tar.gz | tar -xvz && \
    mv sqlc /usr/local/bin/

# Set the current working directory inside the container
WORKDIR /app

# Copy the rest of the source code
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make migrate
RUN sqlc generate

EXPOSE 8000

CMD ["go", "run", "main.go"]