package store

var client Factory

type Factory interface {
	Images() ImageStore
	Terminals() TerminalStore
	Close() error
}

func Client() Factory {
	return client
}

func SetClient(factory Factory) {
	client = factory
}
