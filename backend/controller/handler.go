package controller

import (
	"backend/Infrastructure/store"
	"backend/usecase"

	"github.com/jinzhu/gorm"
	pb "github.com/naoyakurokawa/go_grpc_graphql_proto/pb"
	"google.golang.org/grpc"
)

// RegisterTaskService wires the Task service into the provided gRPC server.
func RegisterService(grpcServer *grpc.Server, db *gorm.DB) {
	taskRepo := store.NewTaskRepository(db)
	taskUsecase := usecase.NewTaskUseCase(taskRepo)
	taskController := NewTaskHandler(taskUsecase)
	pb.RegisterTaskServiceServer(grpcServer, taskController)
}
