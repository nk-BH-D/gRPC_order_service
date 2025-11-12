# gRPC_order_service

## Описание

Сервис для обработки заказов на базе gRPC обработчика

## Команды

### Примеры с использованием grpcurl

**Создание заказа**:

```bash
grpcurl -plaintext -d '{"item":"iPhone 52 Pro Max","quantity":5}' \
localhost:50051 api.OrderService/CreateOrder
```

**Получение заказа по ID**:

```bash
grpcurl -plaintext -d '{"id":"your_id"}' \
localhost:50051 api.OrderService/GetOrder
```

**Обновление заказа по ID**:

```bash
grpcurl -plaintext -d '{"id":"your_id", "item":"GameLaptop", "quantity":"4"}' \
localhost:50051 api.OrderService/UpdateOrder
```

**Удаление заказа по ID**:

```bash
grpcurl -plaintext -d '{"id":"your_id"}' \
localhost:50051 api.OrderService/DeleteOrder
```

**Получение списка всех заказов**:

```bash
grpcurl -plaintext -d '{}' \
localhost:50051 api.OrderService/ListOrders
```

### Примеры с использованием curl

**Создание заказа**:

```bash
curl -X POST http://localhost:8080/BH/order \
  -H "Content-Type: application/json" \
  -d '{"item": "Laptop", "quantity": 2}'
```

**Получение заказа по ID**:

```bash
curl -X GET http://localhost:8080/BH/order/your_id
```

**Обновление заказа по ID**:

```bash
curl -X PUT http://localhost:8080/BH/order/your_id \
  -H "Content-Type: application/json" \
  -d '{"item": "Tablet", "quantity": 5}'
```

**Удаление заказа по ID**:

```bash
curl -X DELETE http://localhost:8080/BH/order/your_id
```

**Получение списка всех заказов**:

```bash
curl -X GET http://localhost:8080/BH/orders
```
