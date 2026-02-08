package utils

func Find[T any](arr []T, pred func(T) bool) (*T, int) {
	for i, ele := range arr {
		if pred(ele) {
			return &ele, i
		}
	}
	return nil, -1
}