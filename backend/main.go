package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	pb "github.com/naoyakurokawa/go_grpc_graphql_proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Task モデル (GORM用)
type Task struct {
	ID        string `gorm:"primaryKey"`
	Title     string `gorm:"not null"`
	Note      string
	Completed int32     `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TaskServiceServer は gRPC サーバーを実装
type TaskServiceServer struct {
	pb.UnimplementedTaskServiceServer
	db *gorm.DB
}

// NewTaskServiceServer はサーバーを初期化
func NewTaskServiceServer(db *gorm.DB) *TaskServiceServer {
	return &TaskServiceServer{db: db}
}

// GetTasks はすべてのタスクを取得
func (s *TaskServiceServer) GetTasks(ctx context.Context, in *emptypb.Empty) (*pb.TaskList, error) {
	var tasks []Task
	if err := s.db.Find(&tasks).Error; err != nil {
		return nil, err
	}

	var pbTasks []*pb.Task
	for _, task := range tasks {
		pbTasks = append(pbTasks, &pb.Task{
			Id:        task.ID,
			Title:     task.Title,
			Note:      task.Note,
			Completed: task.Completed,
			CreatedAt: timestamppb.New(task.CreatedAt),
			UpdatedAt: timestamppb.New(task.UpdatedAt),
		})
	}
	return &pb.TaskList{Tasks: pbTasks}, nil
}

// CreateTask は新しいタスクを作成
func (s *TaskServiceServer) CreateTask(ctx context.Context, in *pb.CreateTaskRequest) (*pb.Task, error) {
	task := Task{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Title:     in.Input.Title,
		Note:      in.Input.Note,
		Completed: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.db.Create(&task).Error; err != nil {
		return nil, err
	}

	return &pb.Task{
		Id:        task.ID,
		Title:     task.Title,
		Note:      task.Note,
		Completed: task.Completed,
		CreatedAt: timestamppb.New(task.CreatedAt),
		UpdatedAt: timestamppb.New(task.UpdatedAt),
	}, nil
}

// UpdateTask はタスクを更新
func (s *TaskServiceServer) UpdateTask(ctx context.Context, in *pb.UpdateTaskRequest) (*pb.Task, error) {
	var task Task
	if err := s.db.First(&task, "id = ?", in.Input.Id).Error; err != nil {
		return nil, fmt.Errorf("task not found")
	}

	task.Title = in.Input.Title
	task.Note = in.Input.Note
	task.Completed = in.Input.Completed
	task.UpdatedAt = time.Now()

	if err := s.db.Save(&task).Error; err != nil {
		return nil, err
	}

	return &pb.Task{
		Id:        task.ID,
		Title:     task.Title,
		Note:      task.Note,
		Completed: task.Completed,
		CreatedAt: timestamppb.New(task.CreatedAt),
		UpdatedAt: timestamppb.New(task.UpdatedAt),
	}, nil
}

// DeleteTask はタスクを削除
func (s *TaskServiceServer) DeleteTask(ctx context.Context, in *pb.TaskId) (*pb.DeleteTaskResponse, error) {
	if err := s.db.Delete(&Task{}, "id = ?", in.Id).Error; err != nil {
		return &pb.DeleteTaskResponse{Success: false}, err
	}
	return &pb.DeleteTaskResponse{Success: true}, nil
}

func main() {
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root", "password", "db", 3306, "test",
	)

	db, err := gorm.Open("mysql", connectionString) // 修正
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	taskService := NewTaskServiceServer(db)

	pb.RegisterTaskServiceServer(grpcServer, taskService)
	log.Println("Server is running on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
