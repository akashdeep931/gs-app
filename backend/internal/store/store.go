package store

type Store interface {
	GetAll() []int
	Add(item int) error
	Remove(item int) error
	Set(items []int) error
}
