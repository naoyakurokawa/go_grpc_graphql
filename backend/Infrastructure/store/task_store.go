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

// FindAll retrieves every task.
func (r *TaskRepository) FindAll(ctx context.Context) ([]model.Task, error) {
	var tasks []model.Task
	if err := r.db.Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// FindByID retrieves a task by its identifier.
func (r *TaskRepository) FindByID(ctx context.Context, id string) (*model.Task, error) {
	var task model.Task
	if err := r.db.First(&task, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// Create persists a new task entity.
func (r *TaskRepository) Create(ctx context.Context, task model.Task) (*model.Task, error) {
	d := dto.FromModel(task)

	if err := r.db.Create(&d).Error; err != nil {
		return nil, err
	}

	res := d.ToModel()
	return &res, nil
}

// Update persists updates to an existing task entity.
func (r *TaskRepository) Update(ctx context.Context, task *model.Task) error {
	return r.db.Save(task).Error
}

// Delete removes a task by id.
func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	return r.db.Delete(&model.Task{}, "id = ?", id).Error
}
