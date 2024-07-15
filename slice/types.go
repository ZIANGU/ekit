package slice

type equalFunc[T any] func(src, dst T) bool

type matchFunc[T any] func(src T) bool
