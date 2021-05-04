FROM golang:latest

WORKDIR /go/src/golangmicroservices/accounts

ADD . /go/src/golangmicroservices

RUN go mod vendor

CMD go run cmd/main.go

EXPOSE 8081
