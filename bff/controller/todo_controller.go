package controller

import (
	"context"
	"log"

	"github.com/naoyakurokawa/go_grpc_graphql/domain/model"
	"github.com/naoyakurokawa/go_grpc_graphql/domain/repository"
	"github.com/naoyakurokawa/go_grpc_graphql/usecase"
)

// TodoController orchestrates requests for todo operations.
type TodoController struct {
	usecase usecase.TodoUsecase
}

// NewTodoController constructs a TodoController instance.
func NewTodoController(uc usecase.TodoUsecase) *TodoController {
	return &TodoController{usecase: uc}
}

func (c *TodoController) CreateTask(ctx context.Context, input model.NewTask) (*model.Task, error) {
	task, err := c.usecase.CreateTask(ctx, input)
	if err != nil {
		log.Printf("failed to create task: %v", err)
		return nil, err
	}

	return task, nil
}

func (c *TodoController) UpdateTask(ctx context.Context, input model.UpdateTask) (*model.Task, error) {
	task, err := c.usecase.UpdateTask(ctx, input)
	if err != nil {
		log.Printf("failed to update task: %v", err)
		return nil, err
	}

	return task, nil
}

func (c *TodoController) DeleteTask(ctx context.Context, id uint64) (bool, error) {
	ok, err := c.usecase.DeleteTask(ctx, id)
	if err != nil {
		log.Printf("failed to delete task: %v", err)
		return false, err
	}

	return ok, nil
}

func (c *TodoController) ListTasks(ctx context.Context, filter repository.TaskFilter) ([]*model.Task, error) {
	tasks, err := c.usecase.ListTasks(ctx, filter)
	if err != nil {
		log.Printf("failed to fetch tasks: %v", err)
		return nil, err
	}

	return tasks, nil
}

func (c *TodoController) CreateSubTask(ctx context.Context, input model.NewSubTask) (*model.SubTask, error) {
	subTask, err := c.usecase.CreateSubTask(ctx, input)
	if err != nil {
		log.Printf("failed to create sub task: %v", err)
		return nil, err
	}
	return subTask, nil
}

func (c *TodoController) ToggleSubTask(ctx context.Context, id uint64, completed bool) (*model.SubTask, error) {
	subTask, err := c.usecase.ToggleSubTask(ctx, id, completed)
	if err != nil {
		log.Printf("failed to toggle sub task: %v", err)
		return nil, err
	}
	return subTask, nil
}
