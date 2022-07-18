package slice

func Contains[T comparable](needle T, search []T) bool {
	for _, v := range search {
		if v == needle {
			return true
		}
	}

	return false
}
