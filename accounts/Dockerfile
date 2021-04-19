FROM golang:latest

RUN mkdir -p /go/src/golangmicroservices

WORKDIR /go/src/golangmicroservices/accounts

ADD . /go/src/golangmicroservices

RUN go get  -t -v ./...
RUN go get  github.com/canthefason/go-watcher
RUN go mod vendor
RUN go install github.com/canthefason/go-watcher/cmd/watcher

ENTRYPOINT  watcher -run golangmicroservices/accounts/cmd  -watch golangmicroservices/accounts
