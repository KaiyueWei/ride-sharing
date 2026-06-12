package main

import (
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"
	pb "ride-sharing/shared/proto/trip"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	var reqBody previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Failed to parse json data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	// Validate
	if reqBody.UserID == "" {
		http.Error(w, "user ID is required", http.StatusBadRequest)
		return
	}

	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer tripService.Close()

	tripPreview, err := tripService.Client.PreviewTrip(r.Context(), &pb.PreviewTripRequest{
		UserID: reqBody.UserID,
		StartLocation: &pb.Coordinate{
			Latitude:  reqBody.Pickup.Latitude,
			Longitude: reqBody.Pickup.Longitude,
		},
		EndLocation: &pb.Coordinate{
			Latitude:  reqBody.Destination.Latitude,
			Longitude: reqBody.Destination.Longitude,
		},
	})
	if err != nil {
		log.Printf("failed to preview trip: %v", err)
		http.Error(w, "failed to preview trip", http.StatusInternalServerError)
		return
	}

	response := contracts.APIResponse{Data: tripPreview}

	writeJSON(w, http.StatusCreated, response)

}
