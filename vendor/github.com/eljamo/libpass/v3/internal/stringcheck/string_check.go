package stringcheck

func HasElementWithLengthGreaterThanOne(s []string) bool {
	for _, str := range s {
		if len(str) > 1 {
			return true
		}
	}

	return false
}
