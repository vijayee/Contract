package builtin

import (
	"github.com/robertkrimen/otto"
	"github.com/vijayee/Contract"
	"time"
)

func SetTimeout(call otto.FunctionCall, conv contract.Converter) otto.Value {
	if len(call.ArgumentList) < 2 {
		return otto.Value{}
	}
	timelen := call.Argument(1).ToInteger()
	timeout := time.NewTimer(time.Millisecond * timelen)
	<-timeout.C
	value, err := call.Argument(0).Call("this")
	if err != nil {
		return otto.Value{}
	}
	return value
}
var clear []chan bool
func SetInterval(call otto.FunctionCall, conv contract.Converter) otto.Value {
	if len(call.ArgumentList) < 2 {
		return otto.Value{}
	}
	timelen := call.Argument(1).ToInteger()
	clear = append(clear, make(chan bool))
	id := len(clear)-1
	for {
		select{
			case <-clear[id]:
				break
			default:
				timeout := time.NewTimer(time.Millisecond * timelen)
				<-timeout.C
				value, err := call.Argument(0).Call("this")
				if err != nil {
					return otto.Value{}
				}			
		}
		
	}
	return value
}
func ClearInterval(call otto.FunctionCall, conv contract.Converter) otto.Value {
	(call otto.FunctionCall, conv contract.Converter) otto.Value {
}
