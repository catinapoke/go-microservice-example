package repository

import "github.com/catinapoke/go-microservice-example/utils/tx"

type Repository struct {
	provider tx.DBProvider
}

func New(provider tx.DBProvider) *Repository {
	return &Repository{provider: provider}
}

