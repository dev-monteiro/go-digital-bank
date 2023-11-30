FROM golang:1.20 AS base

COPY commons /home/commons

WORKDIR /home/core_banking_mock
COPY core_banking_mock/go.mod .
COPY core_banking_mock/go.sum .
RUN go mod download

WORKDIR /home/core_banking_mock/src
COPY core_banking_mock/src .
RUN go build -o core_banking_mock


FROM base AS dev
RUN go install github.com/cosmtrek/air@latest
WORKDIR /home/core_banking_mock/src
RUN air init
CMD ["air"]

FROM golang:1.20 AS test
COPY --from=base /home/core_banking_mock/src/. .
CMD ["./core_banking_mock"]