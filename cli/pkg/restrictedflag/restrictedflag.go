package restrictedflag

import (
	"fmt"
	"strings"
)

type Flag struct {
	allowed []string
	value   string
}

func New(value string, allowed ...string) *Flag {
	return &Flag{
		allowed: allowed,
		value:   value,
	}
}

func (f *Flag) Set(value string) error {
	for _, v := range f.allowed {
		if v == value {
			f.value = value
			return nil
		}

		if strings.HasPrefix(value, "go=") {
			f.value = strings.TrimPrefix(value, "go=")
			return nil
		}
	}

	return fmt.Errorf("value %v is not allowed", value)
}

func (f *Flag) Get() string {
	return f.value
}

func (f *Flag) String() string {
	return f.value
}

func (f *Flag) Allowed() []string {
	return f.allowed
}

func (f *Flag) Type() string {
	return "string"
}
