package domain

import "context"

func (m *Model) Remove(ctx context.Context, id int, projectId int) (Item, error) {
	item, err := m.repo.DeleteItem(ctx, id, projectId)

	if err != nil {
		return Item{}, err
	}

	return mapRepoItemToModel(*item), nil
}
