package dedupe

// DedupeStrings remove duplicate elements from a slice
func DedupeStrings(arr []string) []string {
	m, uniq := make(map[string]struct{}), make([]string, 0, len(arr))
	for _, v := range arr {
		if _, ok := m[v]; !ok {
			m[v], uniq = struct{}{}, append(uniq, v)
		}
	}
	return uniq
}
