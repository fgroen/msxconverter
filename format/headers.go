package format

// matchesHeader checks if the data starts with the given header
func matchesHeader(data, header []byte) bool {
	if len(data) < len(header) {
		return false
	}
	for i := 0; i < len(header); i++ {
		if data[i] != header[i] {
			return false
		}
	}
	return true
}


