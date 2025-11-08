package dto

import (
	"backend/domain/model"
	"time"
)

type SubTask struct {
	ID          uint64     `gorm:"column:id;primaryKey;autoIncrement;type:bigint unsigned"`
	TaskID      uint64     `gorm:"column:task_id;type:bigint unsigned"`
	Title       string     `gorm:"column:title;type:varchar(255)"`
	Note        string     `gorm:"column:note;type:text"`
	Completed   int        `gorm:"column:completed;type:tinyint"`
	CompletedAt *time.Time `gorm:"column:completed_at;type:datetime"`
	DueDate     *time.Time `gorm:"column:due_date;type:date"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (SubTask) TableName() string {
	return "sub_tasks"
}

func (s SubTask) ToModel() model.SubTask {
	return model.SubTask{
		ID:          s.ID,
		TaskID:      s.TaskID,
		Title:       s.Title,
		Note:        s.Note,
		Completed:   int32(s.Completed),
		CompletedAt: s.CompletedAt,
		DueDate:     s.DueDate,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

func SubTaskFromModel(m model.SubTask) SubTask {
	return SubTask{
		ID:          m.ID,
		TaskID:      m.TaskID,
		Title:       m.Title,
		Note:        m.Note,
		Completed:   int(m.Completed),
		CompletedAt: m.CompletedAt,
		DueDate:     m.DueDate,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
