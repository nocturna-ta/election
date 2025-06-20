package common

import "regexp"

// IsNotUUID checks if a string is not a valid UUID format
func IsNotUUID(s string) bool {
	uuidRegex := `^[a-f0-9]{8}-[a-f0-9]{4}-[1-5][a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$`
	re := regexp.MustCompile(uuidRegex)
	return !re.MatchString(s)
}
