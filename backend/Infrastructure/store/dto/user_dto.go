package dto

import (
	"backend/domain/model"
	"time"
)

type User struct {
	ID           uint64    `gorm:"column:id;primaryKey;autoIncrement;type:bigint unsigned"`
	Email        string    `gorm:"column:email;type:varchar(255);unique"`
	PasswordHash string    `gorm:"column:password_hash;type:varchar(255)"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}

func (u User) ToModel() model.User {
	return model.User{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func UserFromModel(m model.User) User {
	return User{
		ID:           m.ID,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
