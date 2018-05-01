package userservice

import (
	"univ/course/model"
	"fmt"
	"math/rand"
	"gopkg.in/mgo.v2/bson"
	"math"
)

type FakeService struct {
	names []string
	users Array
}

func NewFake(initialCapacity int) Service {
	fakeService := FakeService{
		names: []string{"John", "Bill", "Sue", "Megan", "Alfred", "Rick", "Jake"},
	}
	users := make([]model.User, initialCapacity)
	for i := 0; i < initialCapacity; i++ {
		var user = model.User{
			Name:      fmt.Sprintf("%v%d", fakeService.names[rand.Intn(len(fakeService.names))], i),
			Type:      "User",
			Checked:   model.CHECKED,
			ID:        model.ObjectId(bson.NewObjectId()),
			Followers: []model.ObjectId{},
			Following: []model.ObjectId{},
		}
		users[i] = user
	}
	for i := 0; i < initialCapacity; i++ {
		user := users[i]
		friendsCount := normalInt(45, 0, initialCapacity-2)
		friends := make([]model.ObjectId, 0)
		for j := 0; j < friendsCount; j++ {
			friendIdx := rand.Intn(initialCapacity)
			for ; friendIdx == i; {
				friendIdx = rand.Intn(initialCapacity)
			}
			friend := users[friendIdx]
			friend.Friends = append(friend.Friends, user.ID)
			friends = append(friends, friend.ID)
			users[friendIdx] = friend
		}
		user.Friends = friends

		center := math.Abs(float64(i)-float64(initialCapacity)/2.0) < float64(initialCapacity)/5.0
		var postsCount int
		if center {
			postsCount = initialCapacity/10 + rand.Intn(100)
		} else {
			postsCount = 3 + rand.Intn(5)
		}
		posts := make([]model.Post, 0)
		for j := 0; j < postsCount; j++ {
			mentionsCount := int(math.Min(float64(initialCapacity-1), float64(normalInt(1, 0, 5))))
			mentions := make([]model.Mention, mentionsCount)
			for k := 0; k < len(mentions); k++ {
				mentionedIdx := rand.Intn(initialCapacity)
				for ; mentionedIdx == i; {
					mentionedIdx = rand.Intn(initialCapacity)
				}
				mentionedUser := users[mentionedIdx]
				mentions[k] = model.Mention{
					Type:   model.NETWORK,
					UserId: mentionedUser.ID,
				}
			}
			post := model.Post{
				AuthorId:    user.ID,
				LikesAmount: rand.Intn(initialCapacity),
				Type:        model.POST,
				Mentions:    mentions,
			}
			posts = append(posts, post)
		}
		user.Posts = posts

		users[i] = user
	}
	fakeService.users = users
	return &fakeService
}

func normalInt(mean, min, max int) int {
	return int(math.Min(math.Max(float64(min), rand.NormFloat64()+float64(mean)), float64(max)))
}

func (userService *FakeService) Insert(user *model.User) error {
	userService.users.Add(*user)
	return nil
}

func (userService *FakeService) InsertAll(users []model.User) error {
	for _, user := range users {
		userService.users.Add(user)
	}
	return nil
}

func (userService *FakeService) LoadChecked() (Array, error) {
	return userService.users, nil
}

func (userService *FakeService) LoadAll() (Array, error) {
	return userService.users, nil
}

func (userService *FakeService) GetOne(id model.ObjectId) (model.User, error) {
	for _, u := range userService.users {
		if u.ID == id {
			return u, nil
		}
	}
	return model.User{}, fmt.Errorf("no such user in mock repo")
}
