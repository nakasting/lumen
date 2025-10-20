package utils

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

func GenerateSlug(text string) string {

	normalized := norm.NFD.String(text)

	var sb strings.Builder

	for _, r := range normalized {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		sb.WriteRune(r)
	}

	noAccents := sb.String()

	noAccents = strings.ToLower(noAccents)

	reNonAlnum := regexp.MustCompile(`[^a-z0-9\s]`)
	noAccents = reNonAlnum.ReplaceAllString(noAccents, "")

	noAccents = strings.ReplaceAll(noAccents, " ", "-")

	reMultiDash := regexp.MustCompile(`-+`)
	slug := reMultiDash.ReplaceAllString(noAccents, "-")

	return slug
}
