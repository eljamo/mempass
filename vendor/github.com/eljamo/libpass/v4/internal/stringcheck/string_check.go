package stringcheck

func HasElementWithLengthGreaterThanOne(s []string) bool {
	for _, str := range s {
		if len(str) > 1 {
			return true
		}
	}

	return false
}

func IsElementInSlice(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}
