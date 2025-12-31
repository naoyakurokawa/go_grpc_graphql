package condition

import (
	"backend/domain/model"

	"github.com/jinzhu/gorm"
)

// ApplyTaskConditions applies filters for task queries.
func ApplyTaskConditions(query *gorm.DB, in model.GetTasksRequest) *gorm.DB {
	query = query.Where("deleted_at IS NULL")
	if in.CategoryID != nil {
		query = query.Where("category_id = ?", *in.CategoryID)
	}
	if in.DueDateFrom != nil {
		query = query.Where("due_date >= ?", in.DueDateFrom.Format("2006-01-02"))
	}
	if in.DueDateTo != nil {
		query = query.Where("due_date <= ?", in.DueDateTo.Format("2006-01-02"))
	}
	if in.IncompleteOnly != nil && *in.IncompleteOnly {
		query = query.Where("completed = ?", 0)
	}
	if in.UserID != nil {
		query = query.Where("user_id = ?", *in.UserID)
	}

	return query
}
