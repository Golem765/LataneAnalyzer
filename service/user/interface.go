package userservice

import "univ/course/model"

type Set interface {
	Add(user model.User)
}

type Service interface {
	Insert(*model.User) error
	InsertAll([]model.User) error
	LoadChecked() (Array, error)
	LoadAll() (Array, error)
	GetOne(model.ObjectId) (model.User, error)
}
