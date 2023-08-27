package goodslist

import (
	"context"
	"errors"
	"log"

	"github.com/catinapoke/go-microservice-example/internal/domain"
)

type Handler struct {
	// Model *domain.Model
}

type Response struct {
	Meta  MetaResponse  `json:"meta"`
	Goods []domain.Item `json:"goods"`
}

type Request struct {
	Limit  int `path:"limit"`
	Offset int `path:"offset"`
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
	log.Printf("%+v", r)

	// err := h.Model.AddToCart(ctx, req.User, req.SKU, req.Count)
	return Response{Goods: make([]domain.Item, 0)}, nil
}
