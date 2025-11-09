PROTO_DIR=api
OUT_DIR=api/pkg/api/test
GOOGLE_APIS=api/googleapis

.PHONY: proto build run

all: proto build run

proto:
	protoc -I ${PROTO_DIR} -I ${GOOGLE_APIS} \
		--go_out ${OUT_DIR} --go_opt paths=source_relative \
		--go-grpc_out ${OUT_DIR} --go-grpc_opt paths=source_relative \
		--grpc-gateway_out ${OUT_DIR} --grpc-gateway_opt paths=source_relative \
		${PROTO_DIR}/order.proto

build:
	go build -o bin/server ./cmd/server/main.go

run:
	go run ./cmd/server/main.go
