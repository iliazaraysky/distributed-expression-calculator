FROM golang:1.22-rc-alpine
WORKDIR /worker2
COPY . .
RUN go build -o main .
CMD ["./main"]
