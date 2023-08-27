package domain

import "context"

func (m *Model) List(ctx context.Context, int, offset int) ([]Item, error) {
	return make([]Item, 0), nil
}
