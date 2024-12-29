package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"os"
	"time"
)

const (
	address     = "localhost:50051"
	serviceName = "order_service"
)

func checkHealth(address, serviceName string) int {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return 1
	}
	defer conn.Close()

	client := grpc_health_v1.NewHealthClient(conn)
	request := &grpc_health_v1.HealthCheckRequest{Service: serviceName}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := client.Check(ctx, request)
	if err != nil {
		log.Printf("Health check failed: %v", err)
		return 1
	}

	if response.Status == grpc_health_v1.HealthCheckResponse_SERVING {
		fmt.Println("SERVING")
		return 0
	} else {
		fmt.Println("NOT_SERVING")
		return 1
	}
}

func main() {
	os.Exit(checkHealth(address, serviceName))
}
