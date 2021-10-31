package repository

import _interface "github.com/nightsilvertech/bar/repository/interface"

type Repository struct {
	Data  _interface.DRW
	Cache _interface.CRW
}

func NewRepository() *Repository{
	return &Repository{

	}
}
