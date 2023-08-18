FROM golang:1.20

WORKDIR /home

RUN go install github.com/cosmtrek/air@latest

COPY go.mod .
COPY go.sum .
RUN go mod download

WORKDIR /home/src

RUN air init
CMD ["air"]