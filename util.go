package main

import (
	"regexp"
	"strings"
)

func SanitizeFilename(filename string) string {
	re := regexp.MustCompile(`[<>:"/\\|?*]+`)
	cleaned := re.ReplaceAllString(filename, "")

	cleaned = strings.ReplaceAll(cleaned, " ", "_")

	maxLength := 255
	if len(cleaned) > maxLength {
		cleaned = cleaned[:maxLength]
	}

	cleaned = strings.Trim(cleaned, " .")

	if cleaned == "" {
		return "default_filename"
	}
	return cleaned
}
