package stringcheck

func CheckStringLengths(s []string) bool {
	for _, str := range s {
		if len(str) > 1 {
			return false
		}
	}

	return true
}
