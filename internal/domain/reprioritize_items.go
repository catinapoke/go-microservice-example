package domain

import (
	"context"

	"github.com/catinapoke/go-microservice-example/internal/repository"
)

func (m *Model) UpdatePriority(ctx context.Context, id int, projectId int, newPriority int) ([]ItemPriority, error) {
	var items []repository.GoodsItem
	err := m.txm.RunRepeatableRead(ctx, func(ctxTx context.Context) error {
		var err error
		items, err = m.repo.Reprioritize(ctx, id, projectId, newPriority)
		return err
	})

	if err != nil {
		return nil, err
	}

	result := make([]ItemPriority, 0, len(items))
	for _, elem := range items {
		result = append(result, mapRepoPriorityToModel(elem))
	}

	return result, nil
}
