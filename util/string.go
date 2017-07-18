package util

import "github.com/c9s/inflect"

func Pluralize(s string, count int) string {
	if count == 1 {
		return inflect.Singularize(s)
	} else {
		return inflect.Pluralize(s)
	}
}
