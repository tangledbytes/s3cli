package restrictedflag

import (
	"fmt"
	"regexp"
)

type Flag struct {
	allowed   []string
	value     string
	validator func(string) error
}

func New(value string, allowed ...string) *Flag {
	return &Flag{
		allowed: allowed,
		value:   value,
	}
}

// SetValidator sets a custom validator
func (f *Flag) SetValidator(validator func(string) error) *Flag {
	f.validator = validator
	return f
}

func (f *Flag) Set(value string) error {
	for _, v := range f.allowed {
		match, err := regexp.MatchString(v, value)
		if err != nil {
			return fmt.Errorf("internal error: regex match failed for pattern %s and value %s: %w", v, value, err)
		}
		if !match {
			continue
		}

		if f.validator != nil {
			if err := f.validator(value); err != nil {
				return err
			}
		}

		f.value = value
		return nil
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
