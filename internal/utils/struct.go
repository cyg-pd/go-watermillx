package utils

import (
	"fmt"
	"reflect"
	"strings"
)

// StructName returns a normalized name of the passed structure.
func StructName(v any) string {
	if s, ok := v.(fmt.Stringer); ok {
		return s.String()
	}
	return strings.TrimPrefix(reflect.TypeOf(v).String(), "*")
}
