package repository

import (
	"context"

	"github.com/kiarrobino/diveLog/internal/model"
)

type DiveRepository interface {
	Create(ctx context.Context, dive *model.Dive) error
	GetByID(ctx context.Context, id string) (*model.Dive, error)
	List(ctx context.Context, limit, offset int) ([]*model.Dive, error)
	Delete(ctx context.Context, id string) error
	NextDiveNumber(ctx context.Context) (int, error)
}
