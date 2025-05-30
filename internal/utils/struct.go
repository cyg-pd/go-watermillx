package utils

import (
	"fmt"

	"github.com/cyg-pd/go-reflectx"
)

// StructName returns a normalized name of the passed structure.
func StructName(v any) string {
	if s, ok := v.(fmt.Stringer); ok {
		return s.String()
	}
	return reflectx.Name(v)
}
