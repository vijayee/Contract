package contract

import (
	"errors"
	"github.com/robertkrimen/otto"
	//"log"
	"sync"
)

const moduleDeclaration1 = "(function("
const moduleDeclaration2 = "){"
const moduleDeclaration3 = "})("
const moduleDeclaration4 = ");"

//basic form of a golang function to be run in otto
type api func(otto.FunctionCall, Converter) otto.Value

//used to block other otto operations other than type conversion
type Converter struct {
	vm *otto.Otto
}

//used to block other otto operations other than type conversion
func (c *Converter) ToValue(value interface{}) (otto.Value, error) {
	return c.vm.ToValue(value)
}

type API struct {
	name     string
	Function api    //golang function
	Wrapper  Script // javascript wrapper that contain accessors and scoped and descoped variables
}

func NewApi(name string) API {
	return API{name: name}
}

func (a *API) SetFunction(method api) {
	a.Function = method
}

func (a *API) SetWrapper(wrap Script) {
	a.Wrapper = wrap
}

func (a *API) newGoWrapper(vm *otto.Otto) func(otto.FunctionCall) otto.Value {
	conv := Converter{vm: vm}
	funk := a.Function
	return func(call otto.FunctionCall) otto.Value {
		return funk(call, conv)
	}
}

var registry map[string]API
var reglock sync.Mutex

//add to the list of available api's
func Register(function API) error {
	if function.name == "" || function.Function == nil {
		return errors.New("API has missing fields")
	}
	reglock.Lock()
	if registry == nil {
		registry = make(map[string]API)
	}
	registry[function.name] = function
	reglock.Unlock()
	return nil
}
func UnRegister(name string) {
	delete(registry, name)
}
func exists(value string, vm *otto.Otto) bool {
	val, _ := vm.Get(value)
	exist, _ := val.ToString()
	return exist != "undefined"
}

func LoadAll(vm *otto.Otto) error {
	for key, _ := range registry {
		if err := Load(key, vm); err != nil {
			return err
		}
	}
	return nil
}

func Load(name string, vm *otto.Otto) error {
	current, ok := registry[name]
	if ok == false {
		return errors.New("No API registerd by  the name, " + name)
	}
	if exists(name, vm) {
		return errors.New("Global Namespace Conflict for object name: " + name)
	}

	vm.Set(name, current.newGoWrapper(vm))
	current.Wrapper.Load(vm)
	return nil
}

func LoadSet(nameset []string, vm *otto.Otto) error {
	for _, name := range nameset {
		if err := Load(name, vm); err != nil {
			return err
		}
	}
	return nil
}
func commalist(list *string, value string) {
	if *list != "" {
		*list += ", "
	}
	*list += value
}
func contains(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}
func buildParametersInputs(nameset []string) (string, string) {
	inputs := ""
	parameters := ""

	for _, value := range registry {
		if &value.Wrapper != nil {
			for _, scope := range value.Wrapper.scopes {
				commalist(&parameters, scope)
				if contains(nameset, value.name) {
					commalist(&inputs, scope)
				} else {
					commalist(&inputs, "null")
				}
			}
			for _, descope := range value.Wrapper.descopes {
				commalist(&parameters, descope)
				commalist(&inputs, "null")
			}
			commalist(&parameters, value.name)
			if contains(nameset, value.name) {
				commalist(&inputs, value.name)
			} else {
				commalist(&inputs, "null")
			}
		} else {
			commalist(&parameters, value.name)
			if contains(nameset, value.name) {
				commalist(&inputs, value.name)
			} else {
				commalist(&inputs, "null")
			}
		}
	}
	return parameters, inputs
}

func buildModule(script string, nameset []string) string {
	parameters, inputs := buildParametersInputs(nameset)
	module := moduleDeclaration1
	module += parameters
	module += moduleDeclaration2
	module += script
	module += moduleDeclaration3
	module += inputs
	module += moduleDeclaration4
	return module
}
