FROM golang:latest

RUN mkdir -p /go/src/golangmicroservices

ADD . /go/src/golangmicroservices

RUN go get  -t -v ./...
RUN go get  github.com/canthefason/go-watcher
RUN go install github.com/canthefason/go-watcher/cmd/watcher


ENTRYPOINT  watcher -run golangmicroservices/ads/cmd  -watch golangmicroservices/ads
