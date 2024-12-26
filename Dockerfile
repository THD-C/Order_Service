FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

COPY include/ ./include/

RUN go build -o main cmd/order_service/main.go

EXPOSE 50051

CMD ["./main"]