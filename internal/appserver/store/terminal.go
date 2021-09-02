package store

type TerminalStore interface {
	Get() (interface{}, error)
}
