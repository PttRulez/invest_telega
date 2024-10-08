package main

import (
	"fmt"
	tgGrpc "invest_telega/pkg/grpc"
	"invest_telega/pkg/grpc/server"
	"invest_telega/pkg/logger"
	"invest_telega/telega"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var (
		botToken   = os.Getenv("TOKEN")
		listenAddr = os.Getenv("GRPC_LISTEN_ADDR")
	)

	logger := logger.NewLogger(logger.SetupPrettySlog())

	// Create a new Telega service
	svc, err := telega.New(botToken, logger)
	if err != nil {
		log.Fatal(err)
	}

	// Make a TCP Listener
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	// Make a new GRPC native server with (options)
	grpcServer := grpc.NewServer([]grpc.ServerOption{}...)

	// Register (OUR) GRPC server implementation to the GRPC package.
	tgGrpc.RegisterTelegaServer(grpcServer, server.NewGRPCTelegaServer(svc))
	fmt.Println("GRPC Telega is running on port", listenAddr)

	grpcServer.Serve(ln)
}
