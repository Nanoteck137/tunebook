package utils

import "strings"

func SplitTagString(s string) []string {
	tags := []string{}
	if s != "" {
		tags = strings.Split(s, ",")
	}

	return tags
}

func StringToPtr(value string) *string {
	if value != "" {
		return &value
	}

	return nil
}
