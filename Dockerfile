FROM golang:1.8

RUN go get github.com/golang/dep/cmd/dep
RUN go get github.com/canthefason/go-watcher
RUN go install github.com/canthefason/go-watcher/cmd/watcher

