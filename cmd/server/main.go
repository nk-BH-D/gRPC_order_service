package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/nk-BH-D/three_one/api/pkg/api/test"
	config "github.com/nk-BH-D/three_one/internal/config"
	interceptor "github.com/nk-BH-D/three_one/internal/interceptor"
	service "github.com/nk-BH-D/three_one/internal/service"

	"google.golang.org/grpc"
)

func main() {
	conf := config.Loader()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.Port))
	if err != nil {
		log.Fatalf("Failed listening: %v", err)
	}

	log.Printf("Server listening port %s", conf.Port)

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.LogerInteceptor))
	order := service.NewOrderServiceServer()
	pb.RegisterOrderServiceServer(server, order)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed listening: %v", err)
	}
}
