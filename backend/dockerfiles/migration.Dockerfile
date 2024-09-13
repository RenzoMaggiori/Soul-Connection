FROM golang:1.22-alpine

WORKDIR /app

COPY ./api ./api

COPY ./migration ./migration

COPY ./scripts ./scripts

COPY .env .env

WORKDIR /app/migration

RUN go mod tidy

RUN go build -o migration .

WORKDIR /app

CMD ["./scripts/migration.sh"]
