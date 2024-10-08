proto:
	protoc --go_out=. --go-grpc_out=. proto/telegram.proto

.PHONY: proto
