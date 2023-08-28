package domain

import "context"

func (m *Model) Update(ctx context.Context, id int, projectId int, name string, description string) (Item, error) {
	item, err := m.repo.UpdateItem(ctx, id, projectId, name, description)

	if err != nil {
		return Item{}, err
	}

	return mapRepoItemToModel(*item), nil
}
