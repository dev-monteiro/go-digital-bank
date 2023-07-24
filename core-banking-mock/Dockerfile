FROM golang:1.20
WORKDIR /home

COPY go.mod .
RUN go mod download

COPY . .
RUN go build .
CMD ["./core-banking-mock"]