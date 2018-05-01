package userservice

import (
	"univ/course/model"
	"univ/course/db"
	"gopkg.in/mgo.v2/bson"
)

type SimpleService struct {
	repository db.Repository
}

func NewService(database db.DB, name db.RepositoryName) Service {
	simpleService := SimpleService{
		repository: database.GetRepository(name),
	}
	return &simpleService
}

func (userService *SimpleService) Insert(user *model.User) error {
	return userService.repository.Insert(user)
}

func (userService *SimpleService) InsertAll(users []model.User) error {
	var err error
	for _, u := range users {
		err = userService.repository.Insert(u)
		if err != nil {
			return err
		}
	}
	return nil
}

func (userService *SimpleService) LoadChecked() (Array, error) {
	var users []model.User
	err := userService.repository.FindAll(model.User{Checked: model.CHECKED}, &users)

	return users, err
}

func (userService *SimpleService) LoadAll() (Array, error) {
	var users []model.User
	err := userService.repository.Find(model.User{}, &users)

	return users, err
}

func (userService *SimpleService) GetOne(id model.ObjectId) (model.User, error) {
	var user model.User
	err := userService.repository.FindById(bson.ObjectId(id), &user)
	return user, err
}
