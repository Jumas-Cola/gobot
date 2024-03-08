FROM golang:alpine AS builder

WORKDIR /build

ADD go.mod go.sum .

COPY ./src ./src

WORKDIR /build/src

RUN go build -o ./app/gobot

FROM alpine

WORKDIR /home/app

COPY --from=builder /build/src/app/gobot /home/app/gobot
COPY hamspam.db /home/hamspam.db
COPY .env /home/app/.env

CMD ["./gobot"]
