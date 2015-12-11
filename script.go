package contract

import (
	"errors"
	"github.com/robertkrimen/otto"
	"regexp"
)

type Script struct {
	code     string
	scopes   []string
	descopes []string
}

func NewScript() Script {
	return Script{}
}

func (s *Script) SetScopedVariable(variable string) error {
	jsvarre, _ := regexp.Compile("^[^a-zA-Z_$]|[^\\w$]")
	if jsvarre.MatchString(variable) {
		return errors.New("Variable is not a valid javascript name")
	}
	if contains(s.scopes, variable) {
		return errors.New("Variable has already been scoped")
	}
	if contains(s.descopes, variable) {
		return errors.New("Variable has already been descoped")
	}
	s.scopes = append(s.scopes, variable)
	return nil
}

func (s *Script) SetDescopedVariable(variable string) error {
	jsvarre, _ := regexp.Compile("^[^a-zA-Z_$]|[^\\w$]")
	if jsvarre.MatchString(variable) {
		return errors.New("Variable is not a valid javascript name")
	}
	if contains(s.scopes, variable) {
		return errors.New("Variable has already been scoped")
	}
	if contains(s.descopes, variable) {
		return errors.New("Variable has already been descoped")
	}
	s.descopes = append(s.descopes, variable)
	return nil
}

func (s *Script) SetScriptCode(code string) {
	s.code = code
}

func (s *Script) Load(vm *otto.Otto) error {
	if s.code == "" {
		return errors.New("Script contains no code")
	}
	for _, variable := range s.scopes {
		if exists(variable, vm) {
			return errors.New("Global Namespace Conflict for object name: " + variable)
		}
	}
	for _, variable := range s.descopes {
		if exists(variable, vm) {
			return errors.New("Global Namespace Conflict for object name: " + variable)
		}
	}
	_, err := vm.Run(s.code)
	return err
}
