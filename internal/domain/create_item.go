package domain

import "context"

func (m *Model) Create(ctx context.Context, projectId int, name string) (Item, error) {
	return Item{}, nil
}
