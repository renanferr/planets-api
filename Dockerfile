FROM golang:1.15.2-alpine3.12

WORKDIR /app

RUN apk add --no-cache build-base

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./out/app ./cmd/server

EXPOSE 8080

CMD ["./out/app"]