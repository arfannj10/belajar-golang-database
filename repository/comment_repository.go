package repository

import (
	"belajar-golang-database/entity"
	"context"
)

type CommentRepository interface {
	Insert(ctx context.Context, comment entity.Comment)(entity.Comment, error)
	FindByid(ctx context.Context, id int32)(entity.Comment, error)
	FindAll(ctx context.Context) ([]entity.Comment, error)
}