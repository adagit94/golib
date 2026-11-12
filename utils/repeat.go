package utils

// Retried calls f() and in case of non-nil error it attempts to call f() again, but not more than maxRetries.
func Retried[V any](f func() (V, error), maxRetries uint8) (V, error) {
	result, err := f()

	if err == nil {
		return result, nil
	}

	for range maxRetries {
		result, err = f()

		if err == nil {
			return result, nil
		}
	}

	return result, err
}
