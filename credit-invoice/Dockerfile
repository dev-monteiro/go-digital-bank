FROM golang:1.20

RUN go install github.com/cosmtrek/air@latest

COPY commons /home/commons

WORKDIR /home/credit-invoice

COPY credit-invoice/go.mod .
COPY credit-invoice/go.sum .
RUN go mod download

WORKDIR /home/credit-invoice/src

RUN air init
CMD ["air"]