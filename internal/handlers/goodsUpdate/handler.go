package goodsupdate

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
	domain.Item
}

type Request struct {
	Id          int `path:"id"`
	ProjectId   int `path:"projectId"`
	Name        int `query:"name"`
	Description int `query:"description"`
}

var (
	ErrWrongInput = errors.New("wrong input")
)

func (r Request) Validate() error {
	return nil
}

func (h *Handler) Handle(ctx context.Context, r Request) (Response, error) {
	log.Printf("%+v", r)

	// err := h.Model.AddToCart(ctx, req.User, req.SKU, req.Count)
	return Response{}, nil
}
