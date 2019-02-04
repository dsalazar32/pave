package storage

type Storage interface {
	Read() (string, error)
}
