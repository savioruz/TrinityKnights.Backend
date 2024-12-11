package helper

func StringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func IntOrZero(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}
