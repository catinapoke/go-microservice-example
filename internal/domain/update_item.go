package domain

import "context"

func (m *Model) Update(ctx context.Context, id int, projectId int, name string, description string) (Item, error) {
	return Item{}, nil
}
