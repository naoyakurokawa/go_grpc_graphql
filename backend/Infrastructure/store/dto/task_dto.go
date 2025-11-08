package dto

import (
	"backend/domain/model"
	"time"
)

// Task represents the persistence model for the tasks table.
type Task struct {
	ID         uint64    `gorm:"column:id;primaryKey;autoIncrement;type:bigint unsigned"`
	Title      string    `gorm:"column:title;type:varchar(255)"`
	Note       string    `gorm:"column:note;type:text"`
	Completed  int       `gorm:"column:completed;type:tinyint"`
	CategoryID uint64    `gorm:"column:category_id;type:bigint unsigned"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"` // 自動で現在時刻が設定される
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime"` // 更新時に自動更新される
}

// TableName allows GORM to map the DTO to the tasks table.
func (Task) TableName() string {
	return "tasks"
}

// ToModel converts the DTO into the domain Task entity.
func (t Task) ToModel() model.Task {
	return model.Task{
		ID:         t.ID,
		Title:      t.Title,
		Note:       t.Note,
		Completed:  int32(t.Completed),
		CategoryID: t.CategoryID,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}
}

// FromModel converts the domain Task entity into the DTO form.
func FromModel(task model.Task) Task {
	return Task{
		ID:         task.ID,
		Title:      task.Title,
		Note:       task.Note,
		Completed:  int(task.Completed),
		CategoryID: task.CategoryID,
		CreatedAt:  task.CreatedAt,
		UpdatedAt:  task.UpdatedAt,
	}
}
