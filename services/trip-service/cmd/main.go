package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"ride-sharing/services/trip-service/internal/infrastructure/grpc"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"syscall"

	grpcserver "google.golang.org/grpc"
)

var GrpcAddr = ":9093"

func main() {
	inmemRepo := repository.NewInmemRepository()
	svc := service.NewService(inmemRepo)
	// mux := http.NewServeMux()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(){
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	lis, err := net.Listen("tcp", GrpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpcserver.NewServer()
	// todo initialize grpc handler implementation
	grpc.NewGRPCHandler(grpcServer, svc)

	log.Printf("starting gRPC server Trip Service on port %s", lis.Addr().String())

	go func(){
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