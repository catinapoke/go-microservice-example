package domain

import "context"

func (m *Model) UpdatePriority(ctx context.Context, id int, projectId int, newPriority int) ([]ItemPriority, error) {
	return make([]ItemPriority, 0), nil
}
