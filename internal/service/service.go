package service

import (
	"context"
	"database/sql"

	//"sync"

	"github.com/google/uuid"
	pb "github.com/nk-BH-D/gRPC_os/api/pkg/api/test"
	order_db "github.com/nk-BH-D/gRPC_os/internal/order_db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderServiceServer struct {
	pb.UnimplementedOrderServiceServer
	ord_db *order_db.Postgres
}

func NewOrderServiceServer(datab *order_db.Postgres) *OrderServiceServer {
	return &OrderServiceServer{
		ord_db: datab,
	}
}

func (oss *OrderServiceServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, status.Errorf(codes.Canceled, "request canceled %v", ctx.Err())
	}

	if req == nil || req.Item == "" || req.Quantity <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request")
	}

	id := uuid.NewString()
	if err := oss.ord_db.InsertOrder(ctx, id, req.Item, req.Quantity); err != nil {
		return nil, status.Errorf(codes.Internal, "db error: %v", err)
	}

	return &pb.CreateOrderResponse{Id: id}, nil
}

func (oss *OrderServiceServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, status.Errorf(codes.Canceled, "request canceled %v", ctx.Err())
	}

	item, quantity, err := oss.ord_db.GetOrder(ctx, req.Id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "order %q not found", req.Id)
		}
		return nil, status.Errorf(codes.Internal, "db error: %v", err)
	}

	return &pb.GetOrderResponse{
		Order: &pb.Order{
			Id:       req.Id,
			Item:     item,
			Quantity: quantity,
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

	err := oss.ord_db.UpdateOrder(ctx, req.Id, req.Item, req.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "order %q not found", req.Id)
		}
		return nil, status.Errorf(codes.Internal, "db error: %v", err)
	}

	return &pb.UpdateOrderResponse{
		Order: &pb.Order{
			Id:       req.Id,
			Item:     req.Item,
			Quantity: req.Quantity,
		},
	}, nil
}

func (oss *OrderServiceServer) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	if err := ctx.Err(); err != nil {
		return &pb.DeleteOrderResponse{Success: false}, status.Errorf(codes.Canceled, "request canceled %v", ctx.Err())
	}

	if req == nil {
		return &pb.DeleteOrderResponse{Success: false}, status.Errorf(codes.InvalidArgument, "invalid request")
	}

	err := oss.ord_db.DeleteOrder(ctx, req.Id)
	if err != nil {
		return &pb.DeleteOrderResponse{Success: false}, status.Errorf(codes.NotFound, "order %q not found", req.Id)
	}

	return &pb.DeleteOrderResponse{Success: true}, nil
}

func (oss *OrderServiceServer) ListOrders(ctx context.Context, in *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, status.Errorf(codes.Canceled, "request canceled %v", ctx.Err())
	}

	ordersData, err := oss.ord_db.ListOrders(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db error: %v", err)
	}

	orders := make([]*pb.Order, 0, len(ordersData))
	for _, order := range ordersData {
		orders = append(orders, &pb.Order{
			Id:       order["id"].(string),
			Item:     order["item"].(string),
			Quantity: order["quantity"].(int32),
		})
	}

	return &pb.ListOrdersResponse{Orders: orders}, nil
}
