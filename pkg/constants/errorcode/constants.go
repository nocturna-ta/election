package errorcode

type ErrorDefinition struct {
	Code    int
	Message string
}

var (
	NotFound = ErrorDefinition{
		Code:    40005,
		Message: "Not Found",
	}
	WrongPassword = ErrorDefinition{
		Code:    40006,
		Message: "Wrong Password",
	}
)
