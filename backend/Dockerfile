FROM golang:1.22-rc-alpine
WORKDIR /backend
COPY go.mod go.sum /
RUN go mod download
COPY . .
RUN go build -o main .
CMD ["./main"]
