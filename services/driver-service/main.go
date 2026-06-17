package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpcserver "google.golang.org/grpc"
)

var GrpcAddr = ":9092"

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	lis, err := net.Listen("tcp", GrpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// RabbitMQ connection
	conn, err := amqp.Dial(env.GetString("RABBITMQ_DEFAULT_URI", "amqp://guest:guest@rabbitmq:5672/"))
	if err != nil {
		log.Fatal("failed to connect to rabbitmq")
	}
	defer conn.Close()

	// start the grpc
	service := NewService()

	grpcServer := grpcserver.NewServer()
	// todo initialize grpc handler implementation
	NewGrpcHandler(grpcServer, service)

	log.Printf("starting gRPC server Driver Service on port %s", lis.Addr().String())

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
			cancel()
		}
	}()

	// todo wait for the shutdown signal
	<-ctx.Done()
	log.Printf("Shutting down the server...")
	grpcServer.GracefulStop()

}
