package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/nk-BH-D/three_one/api/pkg/api/test"
	config "github.com/nk-BH-D/three_one/internal/config"
	interceptor "github.com/nk-BH-D/three_one/internal/interceptor"
	service "github.com/nk-BH-D/three_one/internal/service"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	conf := config.Loader()

	//gRPC
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.GRPC_Port))
	if err != nil {
		log.Fatalf("Failed listening: %v", err)
	}

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.LogerInteceptor))
	pb.RegisterOrderServiceServer(server, service.NewOrderServiceServer())
	// позволяет запросам просматривать содержимое прото и файлов им сгенерированыx, благодаря этому в запросе не надо писать путь до него
	reflection.Register(server)

	//запуск сервера
	go func() {
		log.Printf("Server listening port :%s", conf.GRPC_Port)
		if err := server.Serve(listener); err != nil {
			log.Fatalf("Failed to listening: %v", err)
		}
	}()

	//http
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()

	conect_grpc_protocol := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%s", conf.GRPC_Port), conect_grpc_protocol); err != nil {
		log.Fatalf("failed to start gateway: %v", err)
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", conf.HTTP_Port),
		Handler: mux,
	}

	go func() {
		log.Printf("HTTP gateway listening port %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to serve HTTP gateway: %v", err)
		}
	}()

	// обработка сигналов Graceful Shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	httpServer.Shutdown(ctx)
	server.GracefulStop()
	log.Println("Server Stoppe")
}
