package main

import (
	"log"
	"net"

	infrastructure "backend/Infrastructure"
	"backend/controller"

	"google.golang.org/grpc"
)

func main() {
	db, err := infrastructure.NewMySQLConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	controller.RegisterService(grpcServer, db)

	log.Println("Server is running on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
