package goodslist

import (
	"context"
	"errors"

	"github.com/catinapoke/go-microservice-example/internal/domain"
)

type Handler struct {
	Model *domain.Model
}

type Response struct {
	Meta  MetaResponse  `json:"meta"`
	Goods []domain.Item `json:"goods"`
}

type Request struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}

type MetaResponse struct {
	Total   int `json:"total"`   // сколько всего записей
	Removed int `json:"removed"` // сколько записей со статусом Removed = true
	Limit   int `json:"limit"`   // какое ограничение стоит на вывод объетов (например 20шт)
	Offset  int `json:"offset"`  // от какой позиции выводить данные в списке
}

var (
	ErrWrongInput = errors.New("wrong input")
)

func (r Request) Validate() error {
	if r.Limit < 0 {
		return ErrWrongInput
	}

	if r.Offset < 0 {
		return ErrWrongInput
	}

	return nil
}

func (h *Handler) Handle(ctx context.Context, r Request) (Response, error) {
	// Update default values
	if r.Limit == 0 {
		r.Limit = 10
	}

	if r.Offset == 0 {
		r.Offset = 1
	}

	items, err := h.Model.List(ctx, r.Limit, r.Offset)

	meta := MetaResponse{
		Limit:   r.Limit,
		Offset:  r.Offset,
		Removed: 0,
		Total:   len(items),
	}

	for _, item := range items {
		if item.Removed {
			meta.Removed++
		}
	}

	return Response{Meta: meta, Goods: items}, err
}
