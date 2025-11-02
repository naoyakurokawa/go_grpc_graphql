package main

import (
	"log"
	"net"

	infrastructure "backend/Infrastructure"
	"backend/Infrastructure/store"
	"backend/controller"
	"backend/usecase"

	"google.golang.org/grpc"
)

func main() {
	db, err := infrastructure.NewMySQLConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	taskRepo := store.NewTaskRepository(db)
	taskUsecase := usecase.NewTaskUseCase(taskRepo)
	taskHandler := controller.NewTaskHandler(taskUsecase)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	controller.RegisterTaskService(grpcServer, taskHandler)

	log.Println("Server is running on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
