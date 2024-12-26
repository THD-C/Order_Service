FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

# Install protobuf compiler
RUN apk add --no-cache protobuf

# Install Go protobuf plugin
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY . .

COPY include/ ./include/

RUN go build -o main cmd/order_service/main.go

EXPOSE 50051

CMD ["./main"]