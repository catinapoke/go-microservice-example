package goodsremove

import (
	"context"
	"errors"
	"log"
)

type Handler struct {
	// Model *domain.Model
}

type Response struct {
	Id         int  `json:"id"`
	CampaignId int  `json:"campaignId"`
	Removed    bool `json:"removed"` // true
}

type Request struct {
	Id        int `path:"id"`
	ProjectId int `path:"projectId"`
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
