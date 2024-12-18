package helper

import "strings"

type TicketType struct {
	Short string
	Long  string
}

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

func UintOrZero(u *uint) uint {
	if u == nil {
		return 0
	}
	return *u
}

func TicketUpper(s string) TicketType {
	switch strings.ToUpper(s) {
	case "VIP":
		return TicketType{Short: "VIP", Long: "VIP"}
	case "REGULAR":
		return TicketType{Short: "REG", Long: "REGULAR"}
	default:
		return TicketType{}
	}
}
