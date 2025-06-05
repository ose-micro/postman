package common


type IDomain[T any, P any] interface {
	New(param P) (*T, error)
	Existing(param P) (*T, error)
}
