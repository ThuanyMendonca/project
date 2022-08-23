package repository

import "github.com/ThuanyMendonca/project/model"

type TypeRepositorySpy struct {
	ITypeRepository
	TypeResp *model.Type
	TypeErr  error
}

func (t *TypeRepositorySpy) Get(id int64) (*model.Type, error) {
	return t.TypeResp, t.TypeErr
}
