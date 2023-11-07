FROM golang:1.20

RUN go install github.com/cosmtrek/air@latest

COPY commons /home/commons

WORKDIR /home/core_banking_mock

COPY core_banking_mock/go.mod .
COPY core_banking_mock/go.sum .
RUN go mod download

WORKDIR /home/core_banking_mock/src

RUN air init
CMD ["air"]