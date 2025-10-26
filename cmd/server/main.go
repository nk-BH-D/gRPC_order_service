package main

import (
	"log"
	"net"

	pb "github.com/nk-BH-D/three_one/api/pkg/api/test"
	service "github.com/nk-BH-D/three_one/internal/service"
	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed listening: %v", err)
	}

	log.Println("Server listening port 50051")

	server := grpc.NewServer()
	order := service.NewOrderServiceServer()
	pb.RegisterOrderServiceServer(server, order)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed listening: %v", err)
	}
}
