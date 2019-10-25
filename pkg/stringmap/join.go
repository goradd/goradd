package stringmap

// JoinStrings will join a map[string]string with kvSep in between each key and value, and itemSep in between each group
// of those
func JoinStrings (m map[string]string, kvSep string, itemSep string) (ret string) {
	if m == nil {
		return ""
	}
	for k,v := range m {
		if ret != "" {
			ret += itemSep
		}
		ret += k + kvSep + v
	}
	return
}
