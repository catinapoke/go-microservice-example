package domain

import "context"

func (m *Model) List(ctx context.Context, limit int, offset int) ([]Item, error) {
	items, err := m.repo.ListItems(ctx, limit, offset)

	if err != nil {
		return nil, err
	}

	result := make([]Item, 0, len(items))
	for _, elem := range items {
		result = append(result, mapRepoItemToModel(elem))
	}

	return result, nil
}
