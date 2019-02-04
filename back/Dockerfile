FROM golang:stretch

WORKDIR /typhoon

# Go watcher dependency for dev purposes
RUN go get github.com/canthefason/go-watcher
RUN go install github.com/canthefason/go-watcher/cmd/watcher

# Dependencies install
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -v -o typhoon-back
CMD ["./typhoon-back"]
