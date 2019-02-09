FROM golang:stretch

# Setup
RUN mkdir /typhoon
WORKDIR /typhoon

# Install git + the latest Docker CE binaries
RUN apt-get update && apt-get install -y --no-install-recommends \
  git \
  python python-dev python-distribute python-pip \
  apt-transport-https \
  ca-certificates \
  curl \
  gnupg2 \
  software-properties-common && \
  curl -fsSL https://download.docker.com/linux/$(. /etc/os-release; echo "$ID")/gpg > /tmp/dkey; apt-key add /tmp/dkey && \
  add-apt-repository \
      "deb [arch=amd64] https://download.docker.com/linux/$(. /etc/os-release; echo "$ID") \
      $(lsb_release -cs) \
      stable" && \
  apt-get update && \
  apt-get -y install docker-ce \
  && rm -rf /var/lib/apt/lists/* \
  && pip install docker-compose --no-cache-dir

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
