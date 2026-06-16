package main

import (
	"context"
	"log"
	"net/http"
	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/proto/driver"
	"time"

	"github.com/gorilla/websocket"
)

const unregisterTimeout = 5 * time.Second

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


func handleRidersWebSocket(w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v.", err)
		return 
	}
	defer conn.Close()

	userID := r.URL.Query().Get("userID")
	if userID == ""{
		log.Printf("No userID provided")
		return 		
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket upgrade failed: %v.", err)
			break
		}
		log.Printf("Received message: %s", message)

	}



}

func handleDriversWebSocket(w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v.", err)
		return 
	}
	defer conn.Close()

	userID := r.URL.Query().Get("userID")
	if userID == ""{
		log.Printf("No userID provided")
		return 		
	}
	packageSlug := r.URL.Query().Get("packageSlug")
	if packageSlug == ""{
		log.Printf("No packageSlug provided")
		return 		
	}

	ctx := r.Context()

	driverService, err := grpc_clients.NewDriverServiceClient()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		// The request context is already cancelled once the connection closes,
		// so use a fresh, bounded context for the cleanup RPC.
		cleanupCtx, cancel := context.WithTimeout(context.Background(), unregisterTimeout)
		defer cancel()

		if _, err := driverService.Client.UnregisterDriver(cleanupCtx, &driver.RegisterDriverRequest{
			DriverID:    userID,
			PackageSlug: packageSlug,
		}); err != nil {
			log.Printf("Failed to unregister driver %s: %v", userID, err)
		}
		driverService.Close()
		log.Println("Driver unregistered: ", userID)
	}()

	driverData, err := driverService.Client.RegisterDriver(ctx, &driver.RegisterDriverRequest{
		DriverID:    userID,
		PackageSlug: packageSlug,
	})
	if err != nil {
		log.Printf("Failed to register driver %s: %v", userID, err)
		return
	}


	msg := contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: driverData.Driver,
	}

	if err := conn.WriteJSON(msg); err != nil{
		log.Printf("Error sending message: %v", err)
		return 
	}


	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket upgrade failed: %v.", err)
			break
		}
		log.Printf("Received message: %s", message)

	}

}

