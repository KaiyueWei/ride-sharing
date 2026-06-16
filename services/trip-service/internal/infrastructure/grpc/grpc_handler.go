package grpc

import (
	"context"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


type gRPCHandler struct {
	pb.UnimplementedTripServiceServer
	service domain.TripService
}

func NewGRPCHandler(server *grpc.Server, service domain.TripService) *gRPCHandler {
	handler := &gRPCHandler{
		service: service,
	}
	pb.RegisterTripServiceServer(server, handler)
	return handler
}

func (h *gRPCHandler) PreviewTrip(ctx context.Context, req *pb.PreviewTripRequest) (*pb.PreviewTripResponse, error) {
	pickup := req.GetStartLocation()
	destination := req.GetEndLocation()
	userID := req.GetUserID()
	pickupCoords := &types.Coordinate{
		Latitude:  pickup.Latitude,
		Longitude: pickup.Longitude,
	}

	destCoords := &types.Coordinate{
		Latitude:  destination.Latitude,
		Longitude: destination.Longitude,
	}
	route, err := h.service.GetRoute(ctx, pickupCoords, destCoords)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to get route: %v", err)
	}
	estimatedFares := h.service.EstimatedPackagesPriceWithRoute(route)
	fares, err := h.service.GenerateTripFares(ctx, estimatedFares, userID, route)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to generate fares: %v", err)
	}
	return &pb.PreviewTripResponse{
		Route: route.ToProto(),
		RideFares: domain.ToRideFaresProto(fares),
	}, nil
}

func (h *gRPCHandler) CreateTrip(ctx context.Context, req *pb.CreateTripRequest) (*pb.CreateTripResponse, error) {
	fareID := req.GetRideFareID()
	userID := req.GetUserID()
	// 1. Fetch and validate the fare
	rideFare, err := h.service.GetAndValidateFare(ctx, fareID, userID)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to get and validate fares: %v", err)
	}
	// 2. Call create trip
	trip, err := h.service.CreateTrip(ctx, rideFare)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to create a trip: %v", err)
	}
	// 3. We also need to initlize a driver
	// 4. add a comment at the end of the function to publish an event
	return &pb.CreateTripResponse{
		TripID: trip.ID.Hex(),
	}, nil
}