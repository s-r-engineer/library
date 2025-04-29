package libraryNetwork

type GenericConnection interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
}
