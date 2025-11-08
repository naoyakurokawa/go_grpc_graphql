package dto

import (
	"backend/domain/model"
	"time"
)

// Category represents the persistence model for the categories table.
type Category struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement;type:bigint unsigned"`
	Name      string    `gorm:"column:name;type:varchar(255)"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName overrides the default table name.
func (Category) TableName() string {
	return "categories"
}

// ToModel converts DTO to domain model.
func (c Category) ToModel() model.Category {
	return model.Category{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
