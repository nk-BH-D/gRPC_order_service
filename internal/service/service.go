package service

import (
	"context"
	"sync"

	"github.com/google/uuid"
	pb "github.com/nk-BH-D/three_one/api/pkg/api/test"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderServiceServer struct {
	pb.UnimplementedOrderServiceServer
	rmu    sync.RWMutex
	orders map[string]*pb.Order
}

func NewOrderServiceServer() *OrderServiceServer {
	return &OrderServiceServer{
		orders: make(map[string]*pb.Order),
	}
}

func (oss *OrderServiceServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, status.Errorf(codes.Canceled, "request canceled %v", ctx.Err())
	}

	if req == nil || req.Item == "" || req.Quantity <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request")
	}

	oss.rmu.Lock()
	defer oss.rmu.Unlock()

	id := uuid.NewString()
	oss.orders[id] = &pb.Order{
		Id:       id,
		Item:     req.Item,
		Quantity: req.Quantity,
	}

	return &pb.CreateOrderResponse{Id: id}, nil
}

func (oss *OrderServiceServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, status.Errorf(codes.Canceled, "request canceled %v", ctx.Err())
	}

	oss.rmu.RLock()
	data, ok := oss.orders[req.Id]
	oss.rmu.RUnlock()

	if !ok {
		return nil, status.Errorf(codes.NotFound, "order %q not found", req.Id)
	}

	return &pb.GetOrderResponse{
		Order: &pb.Order{
			Id:       data.Id,
			Item:     data.Item,
			Quantity: data.Quantity,
		},
	}, nil
}

func (oss *OrderServiceServer) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, status.Errorf(codes.Canceled, "request canceled %v", ctx.Err())
	}

	if req == nil || req.Item == "" || req.Quantity <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request")
	}

	oss.rmu.Lock()
	defer oss.rmu.Unlock()

	data, ok := oss.orders[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "order %q not found", req.Id)
	}

	data.Item = req.Item
	data.Quantity = req.Quantity

	return &pb.UpdateOrderResponse{
		Order: data,
	}, nil
}

func (oss *OrderServiceServer) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	if err := ctx.Err(); err != nil {
		return &pb.DeleteOrderResponse{Success: false}, status.Errorf(codes.Canceled, "request canceled %v", ctx.Err())
	}

	if req == nil {
		return &pb.DeleteOrderResponse{Success: false}, status.Errorf(codes.InvalidArgument, "invalid request")
	}

	oss.rmu.Lock()
	defer oss.rmu.Unlock()

	_, ok := oss.orders[req.Id]
	if !ok {
		return &pb.DeleteOrderResponse{Success: false}, status.Errorf(codes.NotFound, "order %q not found", req.Id)
	}

	delete(oss.orders, req.Id)

	return &pb.DeleteOrderResponse{Success: true}, nil
}

func (oss *OrderServiceServer) ListOrders(ctx context.Context, in *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, status.Errorf(codes.Canceled, "request canceled %v", ctx.Err())
	}

	oss.rmu.Lock()
	defer oss.rmu.Unlock()

	orders := make([]*pb.Order, 0, len(oss.orders))

	for _, order := range oss.orders {
		orders = append(orders, order)
	}

	return &pb.ListOrdersResponse{Orders: orders}, nil
}
