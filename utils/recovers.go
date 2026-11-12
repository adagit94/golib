package utils

// Safely will attempt to execute f() and will call onPanic(x) with recovered value in case there would be one.
func Safely(f func(), onPanic func(x any)) {
	defer func() {
		if r := recover(); r != nil {
			onPanic(r)
		}
	}()

	f()
}