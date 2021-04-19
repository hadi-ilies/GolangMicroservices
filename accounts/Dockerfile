FROM golang:latest

WORKDIR /go/src/golangmicroservices/accounts

ADD . /go/src/golangmicroservices

RUN go mod vendor
RUN go get  -t -v ./...

#RUN go install github.com/canthefason/go-watcher/cmd/watcher

#ENTRYPOINT  watcher -run golangmicroservices/accounts/cmd  -watch golangmicroservices/accounts
CMD go run cmd/main.go