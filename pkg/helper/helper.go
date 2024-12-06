package helper

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type ContextHelper struct{}

func NewContextHelper() *ContextHelper {
	return &ContextHelper{}
}
