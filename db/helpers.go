package db

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func computeName(displayName string) string {
	// remove diacritic marks
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	name, _, _ := transform.String(t, displayName)

	// remove apostrophes
	name = strings.ReplaceAll(name, "'", "")

	// set to lower case
	name = strings.ToLower(name)

	// restrict to alphanumeric
	reg := regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
	name = reg.ReplaceAllString(name, " ")

	// trim leading/trailing whitespace
	name = strings.Trim(name, " ")

	// replace each block of spaces with a hyphen
	reg = regexp.MustCompile(`\s+`)
	name = reg.ReplaceAllString(name, "-")

	return name
}
