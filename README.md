# Order Service

This Go microservice is responsible for managing and continuously updating orders based on price changes. Key features
include:

- Order Management: Provides functions to create, update, and delete orders.
- Continuous Updates: Continuously updates orders based on real-time price changes.
- Distributed Tracing: Integrated with OpenTelemetry to capture and send trace data to Grafana Tempo (host: Tempo, port:
  4317), a scalable, open-source distributed tracing backend.
- Prometheus Metrics: Exposes Prometheus-compatible metrics on port 8111 for real-time monitoring and performance
  tracking.
- gRPC Server: Operates a gRPC server on port 50051 for fast and reliable client-server communication.

## Requirements

- Go 1.23+
- Docker (optional)
- Docker Compose (optional)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/THD-C/Order_Service.git
   cd order_service
   ```

2. Install Go dependencies:
   ```bash
   go mod tidy
   ```

## Configuration

Set environment variables in the `.env` file or directly in your system environment:

- `APPLICATION_ADDR`              - Application address (default: `0.0.0.0`)
- `APPLICATION_PORT`              - Application port (default: `50051`)
- `PROMETHEUS_PORT`               - Prometheus port (default: `8111`)
- `DB_MANAGER_ADDRESS`            - Database manager address
- `COINGECKO_SERVICE_ADDRESS`     - Coingecko service address
- `DB_MANAGER_TIMEOUT`            - Database manager timeout (in seconds, default: `30`)
- `COINGECKO_SERVICE_TIMEOUT`     - Coingecko service timeout (in seconds, default: `30`)
- `COINGECKO_POLLING_FREQUENCY`   - Coingecko polling frequency (in seconds, default: `60`)
- `PENDING_ORDER_CHECK_FREQUENCY` - Pending orders check frequency (in seconds, default: `60`)

## Running

1. Start the Go service:
   ```bash
   go build cmd/order_service/main.go
   ```

2. (Optional) Start the services using Docker Compose:
   ```bash
   docker-compose up
   ```

## gRPC Services

### Wallets

| Method Name  | Request Type                     | Response Type                    | Description                     |
|--------------|----------------------------------|----------------------------------|---------------------------------|
| CreateWallet | [Wallet](/Docs/wallet.md#wallet) | [Wallet](/Docs/wallet.md#wallet) | Create a new wallet for a user. |
| UpdateWallet | [Wallet](/Docs/wallet.md#wallet) | [Wallet](/Docs/wallet.md#wallet) | Update an existing wallet.      |

### Orders

| Method Name | Request Type                                | Response Type                               | Description                        |
|-------------|---------------------------------------------|---------------------------------------------|------------------------------------|
| CreateOrder | [OrderDetails](/Docs/order.md#orderdetails) | [OrderDetails](/Docs/order.md#orderdetails) | Create a new order.                |
| UpdateOrder | [OrderDetails](/Docs/order.md#orderdetails) | [OrderDetails](/Docs/order.md#orderdetails) | Update an existing pending order.  |
| DeleteOrder | [OrderID](/Docs/order.md#orderid)           | [OrderDetails](/Docs/order.md#orderdetails) | Delete an existing pending  order. |
