package common

import "fmt"

type Config interface {
	Load() error
	Check() error
	NewVariableName(name string) string
}

func NewBasePrefix(appPrefix string, componentPrefix string) string {
	return fmt.Sprintf("%s_%s", appPrefix, componentPrefix)
}
