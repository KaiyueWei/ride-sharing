package main

import (
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	var reqBody previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Failed to parse json data")
		return
	}
	defer r.Body.Close()
	// Validate
	if reqBody.UserID == "" {
		writeJSONError(w, http.StatusBadRequest, "user ID is required")
		return
	}

	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer tripService.Close()

	tripPreview, err := tripService.Client.PreviewTrip(r.Context(), reqBody.toProto())
	if err != nil {
		log.Printf("Failed to preview the trip: %v", err)
		writeJSONError(w, http.StatusBadRequest, "Failed to preview the trip")
		return
	}

	response := contracts.APIResponse{Data:tripPreview}

	writeJSON(w, http.StatusCreated, response)

}


func handleTripStart(w http.ResponseWriter, r *http.Request) {
	var reqBody startTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Failed to parse json data")
		return
	}
	defer r.Body.Close()
	// Validate
	if reqBody.UserID == "" {
		writeJSONError(w, http.StatusBadRequest, "user ID is required")
		return
	}

	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer tripService.Close()

	trip, err := tripService.Client.CreateTrip(r.Context(), reqBody.toProto())
	if err != nil {
		log.Printf("Failed to create the trip: %v", err)
		writeJSONError(w, http.StatusBadRequest, "Failed to create the trip")
		return
	}

	response := contracts.APIResponse{Data:trip}

	writeJSON(w, http.StatusCreated, response)

}

