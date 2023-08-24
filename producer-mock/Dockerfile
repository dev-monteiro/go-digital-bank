FROM golang:1.20

RUN go install github.com/cosmtrek/air@latest

COPY commons /home/commons

WORKDIR /home/producer-mock

COPY producer-mock/go.mod .
COPY producer-mock/go.sum .
RUN go mod download

WORKDIR /home/producer-mock/src

RUN air init
CMD ["air"]