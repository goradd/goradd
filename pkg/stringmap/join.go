package stringmap

// JoinStrings will join a map[string]string with kvSep in between each key and value, and itemSep in between each group
// of those
func JoinStrings (m map[string]string, kvSep string, itemSep string) (ret string) {
	if m == nil {
		return ""
	}
	keys := SortedKeys(m)
	for _,k := range keys {
		v := m[k]
		if ret != "" {
			ret += itemSep
		}
		ret += k + kvSep + v
	}
	return
}
