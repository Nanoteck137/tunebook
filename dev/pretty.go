package dev

import (
	"fmt"
	"reflect"
)

func Println(a ...any) (n int, err error) {
	return fmt.Println(wrap(a)...)
}

func Sprint(a ...any) string {
	return fmt.Sprint(wrap(a)...)
}

func wrap(a []any) []any {
	w := make([]any, len(a))
	for i, x := range a {
		w[i] = formatter{v: reflect.ValueOf(x), force: true, quote: true}
	}
	return w
}
