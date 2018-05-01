package db

type RepositoryName string

const (
	Users    RepositoryName = "users"
	Simulate RepositoryName = "userssim"
)

type DB interface {
	GetRepository(name RepositoryName) Repository
	Close()
}

type Repository interface {
	Find(query, results interface{}) error
	FindAll(query, result interface{}) error
	FindById(id, result interface{}) error
	Insert(docs interface{}) error
	Update(selector interface{}, update interface{}) error
}
