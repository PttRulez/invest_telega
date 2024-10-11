package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

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

	investorEndpoint := os.Getenv("TG_CLIENT_INVESTOR_ENDPOINT")
	// Create a new Telega service
	svc, err := telega.New(botToken, investorEndpoint, logger)
	if err != nil {
		log.Fatal(err)
	}

	// Make a TCP Listener
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", listenPort))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		e := ln.Close()
		if e != nil {
			log.Fatal(e)
		}
		fmt.Printf("TCP :%s has been closed", listenPort)
	}()

	// Make a new GRPC native server with (options)
	grpcServer := grpc.NewServer([]grpc.ServerOption{}...)

	// Register (OUR) GRPC server implementation to the GRPC package.
	protogen.RegisterTelegaServer(grpcServer, grpctransport.NewGRPCTelegaServer(svc))
	fmt.Println("GRPC Telega is running on port", listenPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-quit
		svc.Close()
		e := ln.Close()
		if e != nil {
			log.Fatal(e)
		}
		fmt.Printf("TCP :%s has been closed inside graceful", listenPort)
	}()
	grpcServer.Serve(ln)
}
