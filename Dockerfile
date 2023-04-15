# syntax=docker/dockerfile:1

FROM golang:1.20-alpine

WORKDIR /app

COPY . ./

RUN go mod download

RUN go build

EXPOSE 8000

CMD ["./bore-score-api"]
