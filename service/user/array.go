package userservice

import "univ/course/model"

type Array []model.User

func (array *Array) Add(user model.User) {
	contains := false
	var realArray = []model.User(*array)
	for i := 0; i < len(realArray); i++ {
		if realArray[i].ID == user.ID {
			contains = true
			break
		}
	}
	if !contains {
		*array = append(*array, user)
	}
}
