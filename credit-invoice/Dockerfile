FROM golang:1.20
WORKDIR /home

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -C src -o credit-invoice
RUN mv ./src/credit-invoice .
CMD ["./credit-invoice"]