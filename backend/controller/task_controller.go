package controller

import (
	"context"
	"time"

	"backend/domain/model"
	"backend/domain/repository"
	"backend/usecase"

	pb "backend/pkg/pb"

	"github.com/labstack/gommon/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TaskController bridges gRPC requests with task use cases.
type TaskController struct {
	pb.UnimplementedTaskServiceServer
	usecase usecase.TaskUseCase
}

// NewTaskController constructs a TaskController.
func NewTaskController(uc usecase.TaskUseCase) *TaskController {
	return &TaskController{usecase: uc}
}

// GetTasks handles retrieval of all tasks with optional filtering.
func (h *TaskController) GetTasks(ctx context.Context, in *pb.GetTasksRequest) (*pb.TaskList, error) {
	// log.Infof("received GetTasks request")
	filter := repository.TaskFilter{}
	if in != nil {
		filter.CategoryID = in.CategoryId
		filter.DueDateFrom = timestampToTime(in.DueDateStart)
		filter.DueDateTo = timestampToTime(in.DueDateEnd)
		filter.IncompleteOnly = in.IncompleteOnly
	}
	tasks, err := h.usecase.ListTasks(ctx, filter)
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
func (h *TaskController) CreateTask(ctx context.Context, in *pb.CreateTaskRequest) (*pb.Task, error) {
	task := toModelTaskFromCreateTaskRequest(in)
	res, err := h.usecase.CreateTask(ctx, task)
	if err != nil {
		return nil, err
	}

	return toPBTask(*res)
}

// UpdateTask handles updates to a task.
func (h *TaskController) UpdateTask(ctx context.Context, in *pb.UpdateTaskRequest) (*pb.Task, error) {
	task, err := h.usecase.UpdateTask(ctx, toUpdateTaskRequest(in))
	if err != nil {
		return nil, err
	}

	return toPBTask(*task)
}

// DeleteTask handles deleting a task.
func (h *TaskController) DeleteTask(ctx context.Context, in *pb.TaskId) (*pb.DeleteTaskResponse, error) {
	if err := h.usecase.DeleteTask(ctx, in.Id); err != nil {
		return &pb.DeleteTaskResponse{Success: false}, err
	}

	return &pb.DeleteTaskResponse{Success: true}, nil
}

func toModelTaskFromCreateTaskRequest(in *pb.CreateTaskRequest) model.Task {
	return model.Task{
		Title:      in.Input.Title,
		Note:       in.Input.Note,
		DueDate:    timestampToTime(in.Input.DueDate),
		CategoryID: in.Input.CategoryId,
		Completed:  0,
	}
}

func toPBTask(task model.Task) (*pb.Task, error) {
	return &pb.Task{
		Id:          task.ID,
		Title:       task.Title,
		Note:        task.Note,
		Completed:   task.Completed,
		CompletedAt: timeToTimestamp(task.CompletedAt),
		CategoryId:  task.CategoryID,
		DueDate:     timeToTimestamp(task.DueDate),
		CreatedAt:   timestamppb.New(task.CreatedAt),
		UpdatedAt:   timestamppb.New(task.UpdatedAt),
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
	if in.Input.DueDate != nil {
		req.DueDate = timestampToTime(in.Input.DueDate)
	}
	if in.Input.CompletedAt != nil {
		req.CompletedAt = timestampToTime(in.Input.CompletedAt)
	}
	return req
}
func timestampToTime(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}

	t := ts.AsTime()
	return &t
}

func timeToTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}

	return timestamppb.New(*t)
}
