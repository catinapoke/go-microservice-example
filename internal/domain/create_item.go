package domain

import "context"

func (m *Model) Create(ctx context.Context, projectId int, name string) (Item, error) {
	item, err := m.repo.CreateItem(ctx, projectId, name)

	if err != nil {
		return Item{}, err
	}

	return mapRepoItemToModel(*item), nil
}
