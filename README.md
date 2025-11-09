# gRPC_order_service

## Описание

Сервис для обработки заказов на базе gRPC обработчика

## Команды

### Примеры с использованием grpcurl

```bash
grpcurl -plaintext -d '{"item":"iPhone 52 Pro Max","quantity":5}' \
localhost:50051 api.OrderService/CreateOrder
```

```bash
grpcurl -plaintext -d '{"id":"your_id"}' \
localhost:50051 api.OrderService/GetOrder
```

```bash
grpcurl -plaintext -d '{"id":"your_id", "item":"GameLaptop", "quantity":"4"}' \
localhost:50051 api.OrderService/UpdateOrder
```

```bash
grpcurl -plaintext -d '{}' \
localhost:50051 api.OrderService/ListOrders
```

```bash
grpcurl -plaintext -d '{"id":"your_id"}' \
localhost:50051 api.OrderService/DeleteOrder
```

### Примеры с использованием curl

```bash
curl -X POST http://localhost:8080/5.1/order \
  -H "Content-Type: application/json" \
  -d '{"item": "Laptop", "quantity": 2}'
```

```bash
curl -X GET http://localhost:8080/5.1/order/your_id
```

```bash
curl -X PUT http://localhost:8080/5.1/order/your_id \
  -H "Content-Type: application/json" \
  -d '{"item": "Tablet", "quantity": 5}'
```

```bash
curl -X DELETE http://localhost:8080/5.1/order/your_id
```

```bash
curl -X GET http://localhost:8080/5.1/orders
```
