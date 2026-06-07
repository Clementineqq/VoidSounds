package utils

import "github.com/gosimple/slug"

func GenerateSlug(text string) string {
	return slug.Make(text)
}
