package store

import (
	"context"

	"backend/Infrastructure/store/dto"
	"backend/domain/model"
	"backend/domain/repository"

	"github.com/jinzhu/gorm"
)

type SubTaskRepository struct {
	db *gorm.DB
}

func NewSubTaskRepository(db *gorm.DB) repository.SubTaskRepository {
	return &SubTaskRepository{db: db}
}

func (r *SubTaskRepository) ListByTaskID(ctx context.Context, taskID uint64) ([]model.SubTask, error) {
	var subTaskDTOs []dto.SubTask
	if err := r.db.Where("task_id = ?", taskID).Find(&subTaskDTOs).Error; err != nil {
		return nil, err
	}

	result := make([]model.SubTask, 0, len(subTaskDTOs))
	for _, st := range subTaskDTOs {
		result = append(result, st.ToModel())
	}

	return result, nil
}

func (r *SubTaskRepository) Create(ctx context.Context, in model.SubTask) (*model.SubTask, error) {
	d := dto.SubTaskFromModel(in)
	if err := r.db.Create(&d).Error; err != nil {
		return nil, err
	}
	res := d.ToModel()
	return &res, nil
}

func (r *SubTaskRepository) Update(ctx context.Context, in model.SubTask) (*model.SubTask, error) {
	d := dto.SubTaskFromModel(in)
	if err := r.db.Save(&d).Error; err != nil {
		return nil, err
	}
	res := d.ToModel()
	return &res, nil
}

func (r *SubTaskRepository) FindByID(ctx context.Context, id uint64) (*model.SubTask, error) {
	var d dto.SubTask
	if err := r.db.First(&d, "id = ?", id).Error; err != nil {
		return nil, err
	}
	res := d.ToModel()
	return &res, nil
}
