package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/pttrulez/invest_telega/internal/grpctransport"
	"github.com/pttrulez/invest_telega/internal/telega"
	"github.com/pttrulez/invest_telega/pkg/logger"
	"github.com/pttrulez/invest_telega/pkg/protogen"

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
		listenPort = os.Getenv("GRPC_LISTEN_PORT")
	)

	logger := logger.NewLogger(logger.SetupPrettySlog())

	// Create a new Telega service
	svc, err := telega.New(botToken, logger)
	if err != nil {
		log.Fatal(err)
	}

	// Make a TCP Listener
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", listenPort))
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	// Make a new GRPC native server with (options)
	grpcServer := grpc.NewServer([]grpc.ServerOption{}...)

	// Register (OUR) GRPC server implementation to the GRPC package.
	protogen.RegisterTelegaServer(grpcServer, grpctransport.NewGRPCTelegaServer(svc))
	fmt.Println("GRPC Telega is running on port", listenPort)

	grpcServer.Serve(ln)
}
