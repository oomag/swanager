package lib

import (
	"regexp"
	"strings"
)

// IdentifierName formats input string to valid id format
func IdentifierName(s string) string {
	cleanupRegexp := regexp.MustCompile("[^A-Za-z0-9 ]+")

	s = strings.ToLower(s)
	s = cleanupRegexp.ReplaceAllString(s, "")
	s = strings.Trim(s, " \n")
	return strings.Replace(s, " ", "-", -1)
}
