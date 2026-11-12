package types

type Result[R any] struct {
	Result R
	Error error
}