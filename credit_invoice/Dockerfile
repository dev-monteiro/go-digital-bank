FROM golang:1.20

RUN go install github.com/cosmtrek/air@latest

COPY commons /home/commons

WORKDIR /home/credit_invoice

COPY credit_invoice/go.mod .
COPY credit_invoice/go.sum .
RUN go mod download

WORKDIR /home/credit_invoice/src

RUN air init
CMD ["air"]