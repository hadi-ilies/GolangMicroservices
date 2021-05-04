FROM golang:latest

WORKDIR /go/src/golangmicroservices/ads

ADD . /go/src/golangmicroservices

RUN go mod vendor

CMD go run cmd/main.go

EXPOSE 8082
