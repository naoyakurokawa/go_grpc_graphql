package store

import (
	"context"

	"backend/Infrastructure/store/dto"
	"backend/domain/model"
	"backend/domain/repository"

	"github.com/jinzhu/gorm"
)

// TaskRepository implements domain.TaskRepository using GORM.
type TaskRepository struct {
	db *gorm.DB
}

// NewTaskRepository creates a new TaskRepository.
func NewTaskRepository(db *gorm.DB) repository.TaskRepository {
	return &TaskRepository{db: db}
}

// FindAll retrieves every task, filtered by provided criteria.
func (r *TaskRepository) FindAll(ctx context.Context, filter repository.TaskFilter) ([]model.Task, error) {
	query := r.db
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.DueDateFrom != nil {
		query = query.Where("due_date >= ?", filter.DueDateFrom.Format("2006-01-02"))
	}
	if filter.DueDateTo != nil {
		query = query.Where("due_date <= ?", filter.DueDateTo.Format("2006-01-02"))
	}
	if filter.IncompleteOnly != nil && *filter.IncompleteOnly {
		query = query.Where("completed = ?", 0)
	}
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}

	var taskDTOs []dto.Task
	if err := query.Find(&taskDTOs).Error; err != nil {
		return nil, err
	}

	tasks := make([]model.Task, 0, len(taskDTOs))
	for _, t := range taskDTOs {
		tasks = append(tasks, t.ToModel())
	}

	return tasks, nil
}

// FindByID retrieves a task by its identifier.
func (r *TaskRepository) FindByID(ctx context.Context, id uint64) (*model.Task, error) {
	var task model.Task
	if err := r.db.First(&task, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

// Create persists a new task entity.
func (r *TaskRepository) Create(ctx context.Context, in model.Task) (*model.Task, error) {
	d := dto.FromModel(in)

	if err := r.db.Create(&d).Error; err != nil {
		return nil, err
	}
	res := d.ToModel()

	return &res, nil
}

// Update persists updates to an existing task entity.
func (r *TaskRepository) Update(ctx context.Context, in model.Task) (*model.Task, error) {
	d := dto.FromModel(in)

	err := r.db.Save(&d).Error
	if err != nil {
		return nil, err
	}

	res := d.ToModel()
	return &res, nil
}

// Delete removes a task by id.
func (r *TaskRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.Delete(&model.Task{}, "id = ?", id).Error
}
