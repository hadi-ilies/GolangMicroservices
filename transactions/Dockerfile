FROM golang:latest

WORKDIR /go/src/golangmicroservices/transactions

ADD . /go/src/golangmicroservices

RUN go mod vendor

CMD go run cmd/main.go
