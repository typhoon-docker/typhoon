FROM golang:stretch

RUN go get github.com/canthefason/go-watcher
RUN go install github.com/canthefason/go-watcher/cmd/watcher

WORKDIR /typhoon
COPY . /typhoon

RUN go build -v -o typhoon-back
CMD ["./typhoon-back"]
