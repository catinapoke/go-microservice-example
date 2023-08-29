package domain

import (
	"context"

	"github.com/catinapoke/go-microservice-example/internal/repository"
)

func (m *Model) Remove(ctx context.Context, id int, projectId int) (Item, error) {
	var item *repository.GoodsItem
	err := m.txm.RunRepeatableRead(ctx, func(ctxTx context.Context) error {
		var err error
		item, err = m.repo.DeleteItem(ctx, id, projectId)
		return err
	})

	if err != nil {
		return Item{}, err
	}

	return mapRepoItemToModel(*item), nil
}
