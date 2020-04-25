package core

import (
	"fmt"
	"regexp"
)

type Label struct {
	value  string
	parent Labelled
}

type Labelled interface {
	validateLabel(*Label) error
}

func (l Label) Get() string {
	return l.value
}

func (l *Label) Set(label string) error {
	old := l.value
	l.value = label
	if err := l.validate(); err != nil {
		l.value = old
		return err
	}
	return nil
}

func (l *Label) validate() error {
	if l.value == "" {
		return fmt.Errorf("label not set")
	}
	if err := validateIdentifier(l.value); err != nil {
		return err
	}
	return l.parent.validateLabel(l)
}

func validateIdentifier(value string) (err error) {
	// TODO: return typed error
	pattern := "[a-zA-Z]([0-9a-zA-Z_])*"
	regex := regexp.MustCompile(fmt.Sprintf("^%s$", pattern))
	if !regex.MatchString(value) {
		err = fmt.Errorf("Identifier must match %s", pattern)
	}
	return
}

func (l Label) Parent() Labelled {
	return l.parent
}

func (l *Label) hasLabel(other *Label) bool {
	return l != other && l.value == other.value
}
