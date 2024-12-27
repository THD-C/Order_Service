FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

# Install protobuf compiler and make
RUN apk add --no-cache protobuf make

# Install Go protobuf plugin
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY . .

RUN make proto

RUN go build -o main cmd/order_service/main.go

EXPOSE 50051

CMD ["./main"]