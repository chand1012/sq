package utils

func StringSliceToAnySlice(s []string) []any {
	var i []any
	for _, v := range s {
		i = append(i, v)
	}
	return i
}
