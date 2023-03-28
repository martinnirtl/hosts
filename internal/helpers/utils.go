package helpers

func SliceContains(stringSlice []string, search string) bool {
	for _, v := range stringSlice {
		if v == search {
			return true
		}
	}

	return false
}
