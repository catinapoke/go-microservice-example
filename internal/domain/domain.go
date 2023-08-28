package domain

import (
	"context"
	"time"

	"github.com/catinapoke/go-microservice-example/internal/repository"
)

type Item struct {
	Id          int       `json:"id"`
	ProjectId   int       `json:"projectId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ItemPriority struct {
	Id       int `json:"id"`
	Priority int `json:"priority"`
}

type TransactionManager interface {
	RunRepeatableRead(ctx context.Context, fn func(ctxTx context.Context) error) error
}

type GoodsRepository interface {
	CreateItem(ctx context.Context, projectId int, name string) (*repository.GoodsItem, error)
	GetItem(ctx context.Context, id int) (*repository.GoodsItem, error)
	UpdateItem(ctx context.Context, id int, projectId int, name string, description string) (*repository.GoodsItem, error)
	DeleteItem(ctx context.Context, id int, projectId int) (*repository.GoodsItem, error)
	ListItems(ctx context.Context, limit int, offset int) ([]repository.GoodsItem, error)
	Reprioritize(ctx context.Context, id int, projectId int, startPriority int) ([]repository.GoodsPriority, error)
}

type Model struct {
	repo GoodsRepository
}

func New(repository GoodsRepository) *Model {
	return &Model{repo: repository}
}

func mapRepoItemToModel(x repository.GoodsItem) Item {
	return Item{
		Id:          x.Id,
		ProjectId:   x.ProjectId,
		Name:        x.Name,
		Description: x.Description,
		Priority:    x.Priority,
		Removed:     x.Removed,
		CreatedAt:   x.CreatedAt,
	}
}

func mapRepoPriorityToModel(x repository.GoodsPriority) ItemPriority {
	return ItemPriority{
		Id:       x.Id,
		Priority: x.Priority,
	}
}
