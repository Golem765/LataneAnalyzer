package db

import (
	"gopkg.in/mgo.v2"
)

type mgoCollection struct {
	*mgo.Collection
}
type mgoDB struct {
	*mgo.Database
}

type MgoConfig struct {
	Url      string
	User     string
	Password string
	Name     string
}

func NewMgo(config MgoConfig) (DB, error) {
	session, err := mgo.Dial(config.Url)
	if err != nil {
		return nil, err
	}

	cred := mgo.Credential{Username: config.User, Password: config.Password, Source: config.Name}
	err = session.Login(&cred)
	if err != nil {
		return nil, err
	}

	db := session.DB(config.Name)
	return &mgoDB{db}, nil
}

func (col *mgoCollection) Find(query, result interface{}) error {
	return col.Collection.Find(query).One(result)
}

func (col *mgoCollection) FindAll(query, result interface{}) error {
	return col.Collection.Find(query).All(result)
}

func (col *mgoCollection) FindById(id, result interface{}) error {
	return col.Collection.FindId(id).One(result)
}

func (col *mgoCollection) Insert(docs interface{}) error {
	return col.Collection.Insert(docs)
}

func (col *mgoCollection) Update(selector interface{}, update interface{}) error {
	return col.Collection.Update(selector, update)
}

func (db *mgoDB) GetRepository(name RepositoryName) Repository {
	return &mgoCollection{db.Database.C(string(name))}
}

func (db mgoDB) Close() {
	db.Database.Session.Close()
}
