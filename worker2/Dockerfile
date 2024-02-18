FROM golang:1.22-rc-alpine
RUN apk update && apk add bash

WORKDIR /worker2
COPY go.mod go.sum /
RUN go mod download
COPY . .
RUN chmod +x wait-for-it.sh

RUN go build -o main .
CMD ["./wait-for-it.sh", "rabbitmq:5672", "--", "./main"]
