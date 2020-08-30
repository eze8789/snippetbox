package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors errors
}

var EmailRegExp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		val := f.Get(field)
		if strings.TrimSpace(val) == "" {
			f.Errors.Add(field, fmt.Sprintf("The field %s cannot be blank", field))
		}
	}
}

func (f *Form) Length(field string, min, max int) {
	val := f.Get(field)
	if val == "" {
		return
	}
	if utf8.RuneCountInString(val) > max || utf8.RuneCountInString(val) < min {
		f.Errors.Add(field, fmt.Sprintf("The field %s should have between %d and %d characters long", field, min, max))
	}
}

func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid")
}

func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("The field %s should have at least %d characters long", field, d))
	}
}

func (f *Form) MatchPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value) {
		f.Errors.Add(field, "Invalid Field")
	}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
