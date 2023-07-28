FROM golang:1.20
WORKDIR /home

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build .
CMD ["./producer-mock"]