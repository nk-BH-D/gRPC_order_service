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
	"time"

	pb "github.com/nk-BH-D/gRPC_os/api/pkg/api/test"
	config "github.com/nk-BH-D/gRPC_os/internal/config"
	interceptor "github.com/nk-BH-D/gRPC_os/internal/interceptor"
	order_db "github.com/nk-BH-D/gRPC_os/internal/order_db"
	service "github.com/nk-BH-D/gRPC_os/internal/service"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	conf := config.Loader()

	// Подключение к Postgres
	pg_db, err := order_db.NewPostgres(conf.DB_URL)
	if err != nil {
		log.Fatalf("failed to connect postgres_db: %v", err)
	}
	defer pg_db.Close()

	// gRPC сервер
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.GRPC_Port))
	if err != nil {
		log.Fatalf("Failed listening: %v", err)
	}
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.LogerInteceptor))
	pb.RegisterOrderServiceServer(server, service.NewOrderServiceServer(pg_db))
	reflection.Register(server)

	go func() {
		log.Printf("Server listening port :%s", conf.GRPC_Port)
		if err := server.Serve(listener); err != nil {
			log.Fatalf("Failed to listening: %v", err)
		}
	}()

	// HTTP Gateway
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()
	conect_grpc_protocol := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%s", conf.GRPC_Port), conect_grpc_protocol); err != nil {
		log.Fatalf("failed to start gateway: %v", err)
	}

	// Добавляем обработчик для health-check
	muxWithHealth := http.NewServeMux()
	muxWithHealth.Handle("/", mux) // Прокидываем существующий Gateway
	muxWithHealth.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()
		if err := pg_db.DB.PingContext(ctx); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "db not ready")
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", conf.HTTP_Port),
		Handler: muxWithHealth,
	}

	go func() {
		log.Printf("HTTP gateway listening port %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to serve HTTP gateway: %v", err)
		}
	}()

	// Graceful Shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	httpServer.Shutdown(ctx)
	server.GracefulStop()
	log.Println("Server Stopped")
}
