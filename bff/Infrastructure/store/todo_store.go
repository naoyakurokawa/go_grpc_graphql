package store

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/naoyakurokawa/go_grpc_graphql/domain/model"
	"github.com/naoyakurokawa/go_grpc_graphql/domain/repository"
	pb "github.com/naoyakurokawa/go_grpc_graphql/pkg/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ repository.TodoRepository = (*TodoStore)(nil)

// TodoStore implements the TodoRepository interface using the gRPC client.
type TodoStore struct {
	client pb.TaskServiceClient
}

// NewTodoStore creates a new TodoStore.
func NewTodoStore(client pb.TaskServiceClient) repository.TodoRepository {
	return &TodoStore{client: client}
}

func (s *TodoStore) CreateTask(ctx context.Context, input model.NewTask) (*model.Task, error) {
	req := &pb.CreateTaskRequest{
		Input: &pb.NewTask{
			Title:      input.Title,
			Note:       input.Note,
			CategoryId: input.CategoryID,
		},
	}

	if input.DueDate != nil {
		ts, err := parseDateString(input.DueDate)
		if err != nil {
			return nil, err
		}
		req.Input.DueDate = ts
	}

	res, err := s.client.CreateTask(ctx, req)
	if err != nil {
		return nil, err
	}

	return toDomainTask(res), nil
}

func (s *TodoStore) UpdateTask(ctx context.Context, input model.UpdateTask) (*model.Task, error) {
	req := &pb.UpdateTaskRequest{
		Input: &pb.UpdateTask{
			Id: input.ID,
		},
	}

	if input.Title != nil {
		req.Input.Title = input.Title
	}
	if input.Note != nil {
		req.Input.Note = input.Note
	}
	if input.Completed != nil {
		req.Input.Completed = input.Completed
	}
	if input.CategoryID != nil {
		value := *input.CategoryID
		req.Input.CategoryId = &value
	}
	if input.DueDate != nil {
		ts, err := parseDateString(input.DueDate)
		if err != nil {
			return nil, err
		}
		req.Input.DueDate = ts
	}

	res, err := s.client.UpdateTask(ctx, req)
	if err != nil {
		return nil, err
	}

	return toDomainTask(res), nil
}

func (s *TodoStore) DeleteTask(ctx context.Context, id uint64) (bool, error) {
	req := &pb.TaskId{Id: id}

	res, err := s.client.DeleteTask(ctx, req)
	if err != nil {
		return false, err
	}

	return res.Success, nil
}

func (s *TodoStore) ListTasks(ctx context.Context, filter repository.TaskFilter) ([]*model.Task, error) {
	req := &pb.GetTasksRequest{}
	if filter.CategoryID != nil {
		req.CategoryId = filter.CategoryID
	}
	if filter.DueDateStart != nil {
		ts, err := parseDateString(filter.DueDateStart)
		if err != nil {
			return nil, err
		}
		req.DueDateStart = ts
	}
	if filter.DueDateEnd != nil {
		ts, err := parseDateString(filter.DueDateEnd)
		if err != nil {
			return nil, err
		}
		req.DueDateEnd = ts
	}

	res, err := s.client.GetTasks(ctx, req)
	if err != nil {
		return nil, err
	}

	tasks := make([]*model.Task, 0, len(res.Tasks))
	for _, task := range res.Tasks {
		tasks = append(tasks, toDomainTask(task))
	}

	return tasks, nil
}

func toDomainTask(task *pb.Task) *model.Task {
	if task == nil {
		return nil
	}

	return &model.Task{
		ID:         task.GetId(),
		Title:      task.GetTitle(),
		Note:       task.GetNote(),
		Completed:  task.GetCompleted(),
		CategoryID: toUint64Ptr(task.GetCategoryId()),
		DueDate:    formatDate(task.GetDueDate()),
		CreatedAt:  formatTimestamp(task.GetCreatedAt()),
		UpdatedAt:  formatTimestamp(task.GetUpdatedAt()),
	}
}

func toUint64Ptr(v uint64) *uint64 {
	if v == 0 {
		return nil
	}
	val := v
	return &val
}

func formatTimestamp(ts *timestamppb.Timestamp) string {
	if ts == nil {
		return ""
	}

	return ts.AsTime().In(time.Local).Format("2006-01-02 15:04:05")
}

func formatDate(ts *timestamppb.Timestamp) *string {
	if ts == nil {
		return nil
	}

	formatted := ts.AsTime().In(time.Local).Format(dateLayout)
	return &formatted
}

const dateLayout = "2006-01-02"

func parseDateString(value *string) (*timestamppb.Timestamp, error) {
	if value == nil {
		return nil, nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil, nil
	}

	parsed, err := time.ParseInLocation(dateLayout, trimmed, time.Local)
	if err != nil {
		return nil, fmt.Errorf("invalid due_date format (expected YYYY-MM-DD): %w", err)
	}

	return timestamppb.New(parsed), nil
}
