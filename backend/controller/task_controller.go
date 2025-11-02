package controller

import (
	"context"
	"fmt"
	"strconv"

	"backend/domain/model"
	"backend/usecase"

	"github.com/labstack/gommon/log"
	pb "github.com/naoyakurokawa/go_grpc_graphql_proto/pb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TaskHandler bridges gRPC requests with task use cases.
type TaskHandler struct {
	pb.UnimplementedTaskServiceServer
	usecase usecase.TaskUseCase
}

// NewTaskHandler constructs a TaskHandler.
func NewTaskHandler(uc usecase.TaskUseCase) *TaskHandler {
	return &TaskHandler{usecase: uc}
}

// GetTasks handles retrieval of all tasks.
func (h *TaskHandler) GetTasks(ctx context.Context, _ *emptypb.Empty) (*pb.TaskList, error) {
	log.Infof("received GetTasks request")
	tasks, err := h.usecase.ListTasks(ctx)
	if err != nil {
		return nil, err
	}

	pbTasks := make([]*pb.Task, 0, len(tasks))
	for _, task := range tasks {
		converted, err := toPBTask(task)
		if err != nil {
			return nil, err
		}
		pbTasks = append(pbTasks, converted)
	}

	return &pb.TaskList{Tasks: pbTasks}, nil
}

// CreateTask handles creation of a task.
func (h *TaskHandler) CreateTask(ctx context.Context, in *pb.CreateTaskRequest) (*pb.Task, error) {
	task := toModelTaskFromCreateTaskRequest(in)
	res, err := h.usecase.CreateTask(ctx, task)
	if err != nil {
		return nil, err
	}

	return toPBTask(*res)
}

// UpdateTask handles updates to a task.
func (h *TaskHandler) UpdateTask(ctx context.Context, in *pb.UpdateTaskRequest) (*pb.Task, error) {
	input := in.GetInput()
	if input == nil {
		return nil, fmt.Errorf("input is required")
	}
	id := strconv.FormatUint(input.GetId(), 10)
	task, err := h.usecase.UpdateTask(ctx, id, input.GetTitle(), input.GetNote(), input.GetCompleted())
	if err != nil {
		return nil, err
	}

	return toPBTask(*task)
}

// DeleteTask handles deleting a task.
func (h *TaskHandler) DeleteTask(ctx context.Context, in *pb.TaskId) (*pb.DeleteTaskResponse, error) {
	id := strconv.FormatUint(in.GetId(), 10)
	if err := h.usecase.DeleteTask(ctx, id); err != nil {
		return &pb.DeleteTaskResponse{Success: false}, err
	}

	return &pb.DeleteTaskResponse{Success: true}, nil
}

func toModelTaskFromCreateTaskRequest(in *pb.CreateTaskRequest) model.Task {
	return model.Task{
		Title:     in.Input.Title,
		Note:      in.Input.Note,
		Completed: 0,
	}
}

func toPBTask(task model.Task) (*pb.Task, error) {
	return &pb.Task{
		Id:        task.ID,
		Title:     task.Title,
		Note:      task.Note,
		Completed: task.Completed,
		CreatedAt: timestamppb.New(task.CreatedAt),
		UpdatedAt: timestamppb.New(task.UpdatedAt),
	}, nil
}
