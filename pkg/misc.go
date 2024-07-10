package pkg

func Contains[A comparable](list []A, item A) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}
