package controller

import (
	"context"

	"backend/domain/model"
	"backend/usecase"

	pb "backend/pkg/pb"

	"github.com/labstack/gommon/log"
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

// GetTasks handles retrieval of all tasks with optional filtering.
func (h *TaskHandler) GetTasks(ctx context.Context, in *pb.GetTasksRequest) (*pb.TaskList, error) {
	log.Infof("received GetTasks request")
	var categoryID *uint64
	if in != nil {
		categoryID = in.CategoryId
	}
	tasks, err := h.usecase.ListTasks(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	pbTasks := make([]*pb.Task, 0, len(tasks))
	for _, task := range tasks {
		converted, err := toPBTask(task)
		if err != nil {
			log.Errorf("failed to convert task to pb.Task: %v", err)
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
	task, err := h.usecase.UpdateTask(ctx, toUpdateTaskRequest(in))
	if err != nil {
		return nil, err
	}

	return toPBTask(*task)
}

// DeleteTask handles deleting a task.
func (h *TaskHandler) DeleteTask(ctx context.Context, in *pb.TaskId) (*pb.DeleteTaskResponse, error) {
	if err := h.usecase.DeleteTask(ctx, in.Id); err != nil {
		return &pb.DeleteTaskResponse{Success: false}, err
	}

	return &pb.DeleteTaskResponse{Success: true}, nil
}

func toModelTaskFromCreateTaskRequest(in *pb.CreateTaskRequest) model.Task {
	return model.Task{
		Title:      in.Input.Title,
		Note:       in.Input.Note,
		CategoryID: in.Input.CategoryId,
		Completed:  0,
	}
}

func toPBTask(task model.Task) (*pb.Task, error) {
	return &pb.Task{
		Id:         task.ID,
		Title:      task.Title,
		Note:       task.Note,
		Completed:  task.Completed,
		CategoryId: task.CategoryID,
		CreatedAt:  timestamppb.New(task.CreatedAt),
		UpdatedAt:  timestamppb.New(task.UpdatedAt),
	}, nil
}

func toUpdateTaskRequest(in *pb.UpdateTaskRequest) model.UpdateTaskRequest {
	req := model.UpdateTaskRequest{
		ID: in.Input.Id,
	}
	if in.Input.Title != nil {
		req.Title = in.Input.Title
	}
	if in.Input.Note != nil {
		req.Note = in.Input.Note
	}
	if in.Input.Completed != nil {
		req.Completed = in.Input.Completed
	}
	if in.Input.CategoryId != nil {
		req.CategoryID = in.Input.CategoryId
	}
	return req
}
