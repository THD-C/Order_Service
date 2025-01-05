# Order Service

## Requirements
- Go 1.23+

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/order_service.git
   cd order_service
   ```

2. Install Go dependencies:
   ```bash
   go mod tidy
   ```

## Configuration

Set environment variables in the `.env` file or directly in your system environment:

- `APPLICATION_ADDR` - Application address (default: `0.0.0.0`)
- `APPLICATION_PORT` - Application port (default: `50051`)
- `PROMETHEUS_PORT` - Prometheus port (default: `8111`)
- `DB_MANAGER_ADDRESS` - Database manager address
- `COINGECKO_SERVICE_ADDRESS` - Coingecko service address
- `DB_MANAGER_TIMEOUT` - Database manager timeout (in seconds, default: `30`)
- `COINGECKO_SERVICE_TIMEOUT` - Coingecko service timeout (in seconds, default: `30`)
- `COINGECKO_POLLING_FREQUENCY` - Coingecko polling frequency (in seconds, default: `60`)
- `PENDING_ORDER_CHECK_FREQUENCY` - Pending orders check frequency (in seconds, default: `60`)

## Running

1. Start the Go service:
   ```bash
   go build cmd/order_service/main.go
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
