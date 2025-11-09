package store

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/naoyakurokawa/go_grpc_graphql/domain/model"
	"github.com/naoyakurokawa/go_grpc_graphql/domain/repository"
	"github.com/naoyakurokawa/go_grpc_graphql/middleware/session"
	pb "github.com/naoyakurokawa/go_grpc_graphql/pkg/pb"
	"google.golang.org/grpc/metadata"
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
	ctxWithUser, err := withUserMetadata(ctx)
	if err != nil {
		return nil, err
	}
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

	res, err := s.client.CreateTask(ctxWithUser, req)
	if err != nil {
		return nil, err
	}

	return toDomainTask(res), nil
}

func (s *TodoStore) UpdateTask(ctx context.Context, input model.UpdateTask) (*model.Task, error) {
	ctxWithUser, err := withUserMetadata(ctx)
	if err != nil {
		return nil, err
	}
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

	res, err := s.client.UpdateTask(ctxWithUser, req)
	if err != nil {
		return nil, err
	}

	return toDomainTask(res), nil
}

func (s *TodoStore) DeleteTask(ctx context.Context, id uint64) (bool, error) {
	ctxWithUser, err := withUserMetadata(ctx)
	if err != nil {
		return false, err
	}
	req := &pb.TaskId{Id: id}

	res, err := s.client.DeleteTask(ctxWithUser, req)
	if err != nil {
		return false, err
	}

	return res.Success, nil
}

func (s *TodoStore) ListTasks(ctx context.Context, filter repository.TaskFilter) ([]*model.Task, error) {
	ctxWithUser, err := withUserMetadata(ctx)
	if err != nil {
		return nil, err
	}
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
	if filter.IncompleteOnly != nil {
		req.IncompleteOnly = filter.IncompleteOnly
	}

	res, err := s.client.GetTasks(ctxWithUser, req)
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

	subTasks := make([]*model.SubTask, 0, len(task.GetSubTasks()))
	for _, st := range task.GetSubTasks() {
		subTasks = append(subTasks, toDomainSubTask(st))
	}

	return &model.Task{
		ID:          task.GetId(),
		Title:       task.GetTitle(),
		Note:        task.GetNote(),
		Completed:   task.GetCompleted(),
		CategoryID:  toUint64Ptr(task.GetCategoryId()),
		DueDate:     formatDate(task.GetDueDate()),
		CompletedAt: formatTimestampPtr(task.GetCompletedAt()),
		CreatedAt:   formatTimestamp(task.GetCreatedAt()),
		UpdatedAt:   formatTimestamp(task.GetUpdatedAt()),
		SubTasks:    subTasks,
	}
}

func (s *TodoStore) CreateSubTask(ctx context.Context, input model.NewSubTask) (*model.SubTask, error) {
	ctxWithUser, err := withUserMetadata(ctx)
	if err != nil {
		return nil, err
	}

	note := ""
	if input.Note != nil {
		note = *input.Note
	}
	req := &pb.CreateSubTaskRequest{
		Input: &pb.NewSubTask{
			TaskId: input.TaskID,
			Title:  input.Title,
			Note:   note,
		},
	}

	if input.DueDate != nil {
		ts, err := parseDateString(input.DueDate)
		if err != nil {
			return nil, err
		}
		req.Input.DueDate = ts
	}

	res, err := s.client.CreateSubTask(ctxWithUser, req)
	if err != nil {
		return nil, err
	}

	return toDomainSubTask(res), nil
}

func (s *TodoStore) ToggleSubTask(ctx context.Context, id uint64, completed bool) (*model.SubTask, error) {
	ctxWithUser, err := withUserMetadata(ctx)
	if err != nil {
		return nil, err
	}
	req := &pb.ToggleSubTaskRequest{
		Id:        id,
		Completed: completed,
	}

	res, err := s.client.ToggleSubTask(ctxWithUser, req)
	if err != nil {
		return nil, err
	}

	return toDomainSubTask(res), nil
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

func formatTimestampPtr(ts *timestamppb.Timestamp) *string {
	if ts == nil {
		return nil
	}

	formatted := formatTimestamp(ts)
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

func toDomainSubTask(sub *pb.SubTask) *model.SubTask {
	if sub == nil {
		return nil
	}

	return &model.SubTask{
		ID:          sub.GetId(),
		TaskID:      sub.GetTaskId(),
		Title:       sub.GetTitle(),
		Note:        sub.GetNote(),
		Completed:   sub.GetCompleted(),
		CompletedAt: formatTimestampPtr(sub.GetCompletedAt()),
		DueDate:     formatDate(sub.GetDueDate()),
		CreatedAt:   formatTimestamp(sub.GetCreatedAt()),
		UpdatedAt:   formatTimestamp(sub.GetUpdatedAt()),
	}
}

func withUserMetadata(ctx context.Context) (context.Context, error) {
	userID, ok := session.UserIDFromContext(ctx)
	if !ok {
		return nil, errors.New("unauthenticated")
	}
	md := metadata.Pairs("x-user-id", fmt.Sprintf("%d", userID))
	return metadata.NewOutgoingContext(ctx, md), nil
}
