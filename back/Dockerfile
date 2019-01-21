FROM golang:stretch

WORKDIR /typhoon
COPY . /typhoon
RUN go build -v -o typhoon-back

CMD ["./typhoon-back"]
