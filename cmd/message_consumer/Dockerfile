FROM golang:1.24

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY . .
# COPY vendor/ ./vendor/

RUN go env -w GOFLAGS="-buildvcs=false"

CMD air
