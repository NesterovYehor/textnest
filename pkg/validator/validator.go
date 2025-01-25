package validator

import "regexp"

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{make(map[string]string)}
}

func (v *Validator) AddErr(key, msg string) {
	if _, exist := v.Errors[key]; !exist {
		v.Errors[key] = msg
	}
}

func (v *Validator) Check(ok bool, key, msg string) {
	if !ok {
		v.AddErr(key, msg)
	}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func Match(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
