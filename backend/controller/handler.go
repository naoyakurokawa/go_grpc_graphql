package controller

import (
	"backend/Infrastructure/store"
	"backend/usecase"

	pb "backend/pkg/pb"

	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"
)

// RegisterTaskService wires the Task service into the provided gRPC server.
func RegisterService(grpcServer *grpc.Server, db *gorm.DB) {
	taskRepo := store.NewTaskRepository(db)
	taskUsecase := usecase.NewTaskUseCase(taskRepo)
	taskController := NewTaskController(taskUsecase)
	pb.RegisterTaskServiceServer(grpcServer, taskController)

	categoryRepo := store.NewCategoryRepository(db)
	categoryUsecase := usecase.NewCategoryUseCase(categoryRepo)
	categoryController := NewCategoryController(categoryUsecase)
	pb.RegisterCategoryServiceServer(grpcServer, categoryController)
}
