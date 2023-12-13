package validator

func HasElementWithLengthGreaterThanOne(s []string) bool {
	for _, str := range s {
		if len(str) > 1 {
			return true
		}
	}

	return false
}

func IsElementInSlice[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}
