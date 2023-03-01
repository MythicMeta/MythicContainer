package utils

func StringSliceContains(source []string, str string) bool {
	for _, v := range source {
		if str == v {
			return true
		}
	}
	return false
}

func RemoveStringFromSliceNoOrder(source []string, str string) []string {
	for index, value := range source {
		if str == value {
			source[index] = source[len(source)-1]
			source[len(source)-1] = ""
			source = source[:len(source)-1]
			return source
		}
	}
	// we didn't find the element to remove
	return source
}
