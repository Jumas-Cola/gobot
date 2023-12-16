FROM golang:alpine AS builder

WORKDIR /build

ADD go.mod .

COPY . .

RUN go build -o ./app/go-telegram-antispam

FROM alpine

WORKDIR /home/app

COPY --from=builder /build/app/go-telegram-antispam /home/app/go-telegram-antispam
COPY hamspam.db /home/app/hamspam.db
COPY .env /home/app/.env


CMD ["./go-telegram-antispam"]
