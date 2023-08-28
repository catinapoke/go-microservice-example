package domain

import "context"

func (m *Model) UpdatePriority(ctx context.Context, id int, projectId int, newPriority int) ([]ItemPriority, error) {
	items, err := m.repo.Reprioritize(ctx, id, projectId, newPriority)

	if err != nil {
		return nil, err
	}

	result := make([]ItemPriority, 0, len(items))
	for _, elem := range items {
		result = append(result, mapRepoPriorityToModel(elem))
	}

	return result, nil
}
