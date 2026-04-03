package handler

import "strings"

// SanitizeString trims and limits length of input string to avoid large inputs
func SanitizeString(s string, max int) string {
	s = strings.TrimSpace(s)
	if len(s) > max {
		return s[:max]
	}
	return s
}

// SanitizeStringSlice applies SanitizeString to each element with provided max length
func SanitizeStringSlice(in []string, max int) []string {
	out := make([]string, 0, len(in))
	for _, v := range in {
		v = SanitizeString(v, max)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}
