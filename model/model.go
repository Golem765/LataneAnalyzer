package model

import (
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"encoding/json"
)

type ObjectId bson.ObjectId
type UserMentions map[ObjectId]int

type User struct {
	ID        ObjectId     `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string       `bson:"name,omitempty" json:"name,omitempty"`
	Link      string       `bson:"fb_link,omitempty" json:"link,omitempty"`
	Type      string       `bson:"type,omitempty" json:"type,omitempty"`
	Friends   []ObjectId   `bson:"friends,omitempty" json:"friends,omitempty"`
	Followers []ObjectId   `bson:"followers,omitempty" json:"followers,omitempty"`
	Following []ObjectId   `bson:"following,omitempty" json:"following,omitempty"`
	Posts     []Post       `bson:"posts,omitempty" json:"posts,omitempty"`
	Checked   Checked      `bson:"checked,omitempty" json:"checked,omitempty"`
	Mentions  UserMentions `bson:"mentions,omitempty" json:"mentions,omitempty"`
	Pressure  float64      `bson:"pressure,omitempty" json:"pressure"`
	Influence float64      `bson:"influence,omitempty" json:"influence"`
}

type UserType string

const (
	USER UserType = "User"
	PAGE UserType = "Page"
)

type Checked int

const (
	NONE       Checked = iota
	PROCESSING
	CHECKED
)

type Post struct {
	AuthorId       ObjectId  `bson:"authorId,omitempty" json:"authorId,omitempty"`
	LikesAmount    int       `bson:"likesAmount,omitempty" json:"likesAmount,omitempty"`
	CommentsAmount int       `bson:"commentsAmount,omitempty" json:"commentsAmount,omitempty"`
	Text           string    `bson:"text,omitempty" json:"text,omitempty"`
	Tags           []string  `bson:"tags,omitempty" json:"tags,omitempty"`
	Urls           []string  `bson:"urls,omitempty" json:"urls,omitempty"`
	Date           string    `bson:"date,omitempty" json:"date,omitempty"`
	Type           PostType  `bson:"type,omitempty" json:"type,omitempty"`
	Mentions       []Mention `bson:"mentions,omitempty" json:"mentions,omitempty"`
}

type PostType string

const (
	POST   PostType = "post"
	REPOST PostType = "repost"
)

type Mention struct {
	Type   MentionType `bson:"type,omitempty" json:"type,omitempty"`
	UserId ObjectId    `bson:"user_id,omitempty" json:"user_id,omitempty"`
	Link   string      `bson:"link,omitempty" json:"link,omitempty"`
	Name   string      `bson:"name,omitempty" json:"name,omitempty"`
}

type MentionType string

const (
	NETWORK   MentionType = "fb"
	FB        MentionType = "fb"
	INSTAGRAM MentionType = "instagram"
)

func (mentions UserMentions) MarshalJSON() ([]byte, error) {
	target := make(map[string]int)
	for key, value := range mentions {
		stringKey := bson.ObjectId(key).Hex()
		target[stringKey] = value
	}
	return json.Marshal(target)
}

func (id ObjectId) MarshalJSON() ([]byte, error) {
	return bson.ObjectId(id).MarshalJSON()
}

func (id ObjectId) String() string {
	return bson.ObjectId(id).String()
}

func (user User) String() string {
	return fmt.Sprintf("{id: %v, name: %v, link: %v, friends: %v}", user.ID, user.Name, user.Link, len(user.Friends))
}
