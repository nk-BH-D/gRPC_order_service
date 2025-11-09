# gRPC_order_service

# Описание
Сервис для обработки заказов на базе gRPC обработчика

# Команды
grpcurl -plaintext -d "{\"item\":\"iPhone 52 Pro Max\",\"quantity\":5}" localhost:50051 api.OrderService/CreateOrder

grpcurl -plaintext -d "{\"id\":\"your id\"}" localhost:50051 api.OrderService/GetOrder

grpcurl -plaintext -d "{\"id\":\"your id\", \"item\":\"GameLaptop\", \"quantity\":\"4\"}" localhost:50051 api.OrderService/UpdateOrder

grpcurl -plaintext -d "" localhost:50051 api.OrderService/ListOrders

grpcurl -plaintext -d "{\"id\":\"your id\"}" localhost:50051 api.OrderService/DeleteOrder
