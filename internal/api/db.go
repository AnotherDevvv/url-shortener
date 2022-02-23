package api

//go:generate moq -pkg mock -out mock/repository.go . Repository
type Repository interface {
	Insert(key string, value string) error
	Get(key string) (string, error)
	Close() error
}
