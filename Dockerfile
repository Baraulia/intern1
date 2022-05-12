FROM golang:alpine AS builder

COPY . /lesson_1/
WORKDIR /lesson_1/

RUN go mod download
RUN GOOS=linux go build -o ./.bin/service ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /lesson_1/.bin/service .
COPY --from=0 /lesson_1/.env .
COPY --from=0 /lesson_1/countries.csv .

EXPOSE 8090

CMD ["./service"]

