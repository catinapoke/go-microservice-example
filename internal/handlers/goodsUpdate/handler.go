package goodsupdate

import (
	"context"
	"errors"
	"log"

	"github.com/catinapoke/go-microservice-example/internal/domain"
)

type Handler struct {
	Model *domain.Model
}

type Response struct {
	domain.Item
}

type Request struct {
	Id          int    `query:"id"`
	ProjectId   int    `query:"projectId"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var (
	ErrWrongInput = errors.New("wrong input")
)

func (r Request) Validate() error {
	if r.Name == "" {
		return ErrWrongInput
	}

	return nil
}

func (h *Handler) Handle(ctx context.Context, r Request) (Response, error) {
	log.Printf("%+v", r)

	item, err := h.Model.Update(ctx, r.Id, r.ProjectId, r.Name, r.Description)
	return Response{item}, err
}
