FROM golang:latest

WORKDIR /go/src/golangmicroservices/auths

ADD . /go/src/golangmicroservices

RUN go mod vendor

CMD go run cmd/main.go

EXPOSE 8084
