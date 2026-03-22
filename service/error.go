package service

import "fmt"

type InvalidFilterError struct {
	Service string
	Message string
}

func (e *InvalidFilterError) Error() string {
	return fmt.Sprintf("%s: invalid filter: %s", e.Service, e.Message)
}

type InvalidSortError struct {
	Service string
	Message string
}

func (e *InvalidSortError) Error() string {
	return fmt.Sprintf("%s: invalid sort: %s", e.Service, e.Message)
}
