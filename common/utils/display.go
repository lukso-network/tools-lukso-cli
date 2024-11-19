package utils

// MaxLength returns the length of the longest string - used for creating padding from a given client set.
func MaxLength(strs []string) (length int) {
	for _, str := range strs {
		if len(str) > length {
			length = len(str)
		}
	}

	return
}
