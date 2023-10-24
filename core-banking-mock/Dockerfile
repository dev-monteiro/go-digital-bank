FROM golang:1.20

RUN go install github.com/cosmtrek/air@latest

COPY commons /home/commons

WORKDIR /home/core-banking-mock

COPY core-banking-mock/go.mod .
COPY core-banking-mock/go.sum .
RUN go mod download

WORKDIR /home/core-banking-mock/src

RUN air init
CMD ["air"]