FROM golang:latest
RUN mkdir -p /go/src/github.com/nktsitas/checkout-techlab
ADD . /go/src/github.com/nktsitas/checkout-techlab
WORKDIR /go/src/github.com/nktsitas/checkout-techlab
RUN go build -o main .
CMD ["/go/src/github.com/nktsitas/checkout-techlab/main"]