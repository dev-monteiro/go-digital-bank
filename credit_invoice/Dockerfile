FROM golang:1.20 AS base

COPY commons /home/commons

WORKDIR /home/credit_invoice
COPY credit_invoice/go.mod .
COPY credit_invoice/go.sum .
RUN go mod download

WORKDIR /home/credit_invoice/src
COPY credit_invoice/src .
RUN go build -o credit_invoice


FROM base AS dev
RUN go install github.com/cosmtrek/air@latest
WORKDIR /home/credit_invoice/src
RUN air init
CMD ["air"]

FROM golang:1.20 AS test
COPY --from=base /home/credit_invoice/src/. .
CMD ["./credit_invoice"]