package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"backend/domain/model"
	mockrepository "backend/domain/repository/mock"

	"github.com/golang/mock/gomock"
)

func TestCategoryUseCase_ListCategories(t *testing.T) {
	t.Parallel()

	errRepository := errors.New("db unavailable")

	tests := []struct {
		name       string
		repoResult []model.Category
		repoErr    error
		wantErr    error
	}{
		{
			name: "success",
			repoResult: []model.Category{
				{ID: 1, Name: "Work"},
				{ID: 2, Name: "Personal"},
			},
		},
		{
			name:    "repository error",
			repoErr: errRepository,
			wantErr: errRepository,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockRepo := mockrepository.NewMockCategoryRepository(ctrl)
			mockRepo.EXPECT().ListCategories(ctx).Return(tt.repoResult, tt.repoErr)

			uc := NewCategoryUseCase(mockRepo)

			got, err := uc.ListCategories(ctx)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("ListCategories error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("ListCategories returned error: %v", err)
			}

			if !reflect.DeepEqual(got, tt.repoResult) {
				t.Fatalf("ListCategories = %#v, want %#v", got, tt.repoResult)
			}
		})
	}
}
