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
	usecase        usecase.TaskUseCase
	subTaskUsecase usecase.SubTaskUseCase
}

// NewTaskController constructs a TaskController.
func NewTaskController(uc usecase.TaskUseCase, sub usecase.SubTaskUseCase) *TaskController {
	return &TaskController{usecase: uc, subTaskUsecase: sub}
}

// GetTasks handles retrieval of all tasks with optional filtering.
func (h *TaskController) GetTasks(ctx context.Context, in *pb.GetTasksRequest) (*pb.TaskList, error) {
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
	for i := range tasks {
		subTasks, err := h.subTaskUsecase.ListByTaskID(ctx, tasks[i].ID)
		if err != nil {
			return nil, err
		}
		tasks[i].SubTasks = subTasks
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

// CreateSubTask handles creation of a sub task.
func (h *TaskController) CreateSubTask(ctx context.Context, in *pb.CreateSubTaskRequest) (*pb.SubTask, error) {
	subTask := toModelSubTaskFromCreateRequest(in)
	res, err := h.subTaskUsecase.Create(ctx, subTask)
	if err != nil {
		return nil, err
	}
	return toPBSubTask(*res), nil
}

// ToggleSubTask handles toggling completion of a sub task.
func (h *TaskController) ToggleSubTask(ctx context.Context, in *pb.ToggleSubTaskRequest) (*pb.SubTask, error) {
	res, err := h.subTaskUsecase.ToggleCompletion(ctx, in.Id, in.Completed)
	if err != nil {
		return nil, err
	}
	return toPBSubTask(*res), nil
}

// ListSubTasks returns subtasks for a task.
func (h *TaskController) ListSubTasks(ctx context.Context, in *pb.TaskId) (*pb.SubTaskList, error) {
	subTasks, err := h.subTaskUsecase.ListByTaskID(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	pbSubTasks := make([]*pb.SubTask, 0, len(subTasks))
	for _, st := range subTasks {
		pbSubTasks = append(pbSubTasks, toPBSubTask(st))
	}

	return &pb.SubTaskList{SubTasks: pbSubTasks}, nil
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
	pbSubTasks := make([]*pb.SubTask, 0, len(task.SubTasks))
	for _, st := range task.SubTasks {
		pbSubTasks = append(pbSubTasks, toPBSubTask(st))
	}
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
		SubTasks:    pbSubTasks,
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

func toModelSubTaskFromCreateRequest(in *pb.CreateSubTaskRequest) model.SubTask {
	return model.SubTask{
		TaskID:    in.Input.TaskId,
		Title:     in.Input.Title,
		Note:      in.Input.Note,
		DueDate:   timestampToTime(in.Input.DueDate),
		Completed: 0,
	}
}

func toPBSubTask(sub model.SubTask) *pb.SubTask {
	return &pb.SubTask{
		Id:          sub.ID,
		TaskId:      sub.TaskID,
		Title:       sub.Title,
		Note:        sub.Note,
		Completed:   sub.Completed,
		CompletedAt: timeToTimestamp(sub.CompletedAt),
		DueDate:     timeToTimestamp(sub.DueDate),
		CreatedAt:   timestamppb.New(sub.CreatedAt),
		UpdatedAt:   timestamppb.New(sub.UpdatedAt),
	}
}
