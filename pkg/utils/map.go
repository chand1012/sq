package utils

func BreakOutMap(m map[string]string) (keys []string, values []string) {
	for k, v := range m {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}
