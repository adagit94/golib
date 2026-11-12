package utils

func ChanToSlice[Ch chan V | <- chan V, V any](ch Ch) []V {
	vals := make([]V, 0, len(ch))
	
	for range ch {
		vals = append(vals, <- ch)
	}

	return vals
}