package repository

import (
	"github.com/nightsilvertech/bar/repository/data"
	_interface "github.com/nightsilvertech/bar/repository/interface"
)

type Repository struct {
	Data   _interface.DRW
	Cache  _interface.CRW
}

func NewRepository() (repo *Repository) {
	dataReadWriter, _ := data.NewDataReadWriter("root", "root", "localhost", "3306", "foobar",)
	return &Repository{
		Data:   dataReadWriter,
	}
}
