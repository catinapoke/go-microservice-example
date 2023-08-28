package localgoods

import (
	"context"
	"errors"
	"time"

	"github.com/catinapoke/go-microservice-example/internal/repository"
)

type Repository struct {
	items    []repository.GoodsItem
	projects []repository.ProjectItem
}

const (
	AllocateLimit = 64 // In case too bit limit in ListItems
)

var (
	ErrNotFound = errors.New("can't find item")
)

func New() *Repository {
	return &Repository{
		items: make([]repository.GoodsItem, 0),
		projects: []repository.ProjectItem{
			{
				Id:        1,
				Name:      "Первая запись",
				CreatedAt: time.Now(),
			},
		},
	}
}

func (r *Repository) CreateItem(ctx context.Context, projectId int, name string) (*repository.GoodsItem, error) {
	// Check project id
	pIndex := r.getProjectIndex(func(x repository.ProjectItem) bool { return x.Id == projectId })
	if pIndex == -1 {
		return nil, ErrNotFound
	}

	// Get last id and priority
	id := 0
	priority := 0

	for _, item := range r.items {
		if item.Id > id {
			id = item.Id
		}

		if item.Priority > priority {
			priority = item.Priority
		}
	}

	id = id + 1
	priority = priority + 1

	item := repository.GoodsItem{
		Id:          id,
		ProjectId:   projectId,
		Name:        name,
		Description: "",
		Priority:    priority,
		Removed:     false,
		CreatedAt:   time.Now(),
	}

	r.items = append(r.items, item)

	return &item, nil
}

func (r *Repository) GetItem(ctx context.Context, id int, projectId int) (*repository.GoodsItem, error) {
	index := r.getItemIndex(func(x repository.GoodsItem) bool { return x.Id == id })

	if index == -1 {
		return nil, ErrNotFound
	}

	return &r.items[index], nil
}

func (r *Repository) UpdateItem(ctx context.Context, id int, projectId int, name string, description string) (*repository.GoodsItem, error) {
	index := r.getItemIndex(func(x repository.GoodsItem) bool { return x.Id == id && x.ProjectId == projectId })

	if index == -1 {
		return nil, ErrNotFound
	}

	item := &r.items[index]

	item.Name = name
	item.Description = description

	return item, nil
}

func (r *Repository) DeleteItem(ctx context.Context, id int, projectId int) (*repository.GoodsItem, error) {
	index := r.getItemIndex(func(x repository.GoodsItem) bool { return x.Id == id && x.ProjectId == projectId })

	if index == -1 {
		return nil, ErrNotFound
	}

	r.items[index].Removed = true

	return &r.items[index], nil
}

func (r *Repository) ListItems(ctx context.Context, limit int, offset int) ([]repository.GoodsItem, error) {
	result := make([]repository.GoodsItem, 0, min(limit, AllocateLimit))

	for _, elem := range r.items {
		if elem.Id >= offset && elem.Id < limit+offset { // Incorrect logic, but it's debug one
			result = append(result, elem)
		}
	}

	return result, nil
}

func (r *Repository) Reprioritize(ctx context.Context, id int, projectId int, startPriority int) ([]repository.GoodsPriority, error) {
	result := make([]repository.GoodsPriority, 0)

	for _, item := range r.items {
		if item.ProjectId != projectId {
			continue
		}

		priority := item.Id - id + startPriority
		if item.Priority == priority {
			continue
		}

		item.Priority = priority

		result = append(result, repository.GoodsPriority{
			Id:        item.Id,
			ProjectId: item.ProjectId,
			Priority:  item.Priority,
		})
	}

	return result, nil
}

func (r *Repository) getItemIndex(filter func(x repository.GoodsItem) bool) int {
	for i, elem := range r.items {
		if filter(elem) {
			return i
		}
	}

	return -1
}

func (r *Repository) getProjectIndex(filter func(x repository.ProjectItem) bool) int {
	for i, elem := range r.projects {
		if filter(elem) {
			return i
		}
	}

	return -1
}

func min(a int, b int) int {
	if a < b {
		return a
	}

	return b
}
