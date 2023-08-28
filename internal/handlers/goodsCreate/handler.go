package goodscreate

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
	ProjectId int    `query:"projectId"`
	Name      string `json:"name"`
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

	item, err := h.Model.Create(ctx, r.ProjectId, r.Name)
	return Response{item}, err
}
