package interceptor

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LogerInteceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// info показывает какой метод интерфейса из .proto ипользуеться при запросе
	// handler вызов этого метода
	if err := ctx.Err(); err != nil {
		return nil, status.Errorf(codes.Canceled, "request canceled %v", ctx.Err())
	}

	log.Printf("operator: %s; req: %+v\n", info.FullMethod, req)

	resp, err := handler(ctx, req)
	if err != nil {
		log.Printf("error: %v", err)
	} else {
		log.Printf("resp: %v", resp)
	}

	return resp, err
}
