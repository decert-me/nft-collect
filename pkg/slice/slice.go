package slice

// DiffArray 求两个切片的差集
func DiffSlice[T string | int](a []T, b []T) []T {
	var diffArray []T
	temp := map[T]struct{}{}

	for _, val := range b {
		if _, ok := temp[val]; !ok {
			temp[val] = struct{}{}
		}
	}

	for _, val := range a {
		if _, ok := temp[val]; !ok {
			diffArray = append(diffArray, val)
		}
	}

	return diffArray
}
