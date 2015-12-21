package builtin

import (
	"github.com/robertkrimen/otto"
	"github.com/vijayee/Contract"
	"time"
)

var SetTimeout contract.API
var SetInterval contract.API
var ClearInterval contract.API

func init() {
	SetTimeout = contract.NewApi("setTimeout")
	SetTimeout.SetFunction(setTimeout)
	SetInterval = contract.NewApi("setInterval")
	SetInterval.SetFunction(setInterval)
	ClearInterval = contract.NewApi("clearInterval")
	ClearInterval.SetFunction(clearInterval)
}
func setTimeout(call otto.FunctionCall, conv contract.Converter) otto.Value {
	if len(call.ArgumentList) < 2 {
		return otto.Value{}
	}
	timelen, err := call.Argument(1).ToInteger()
	if err != nil {
		return otto.Value{}
	}
	obj := call.Argument(0)
	go func() {
		timeout := time.NewTimer(time.Duration(int64(time.Millisecond) * timelen))
		<-timeout.C
		switch {
		case obj.IsString():
			runme, err := obj.ToString()
			if err != nil {
				return
			}
			_, err = call.Otto.Run(runme)
			if err != nil {
				return
			}
		case obj.IsFunction():
			_, err := obj.Call(obj)
			if err != nil {
				return
			}
		default:
			return
		}
	}()
	return otto.Value{}
}

var clear map[int64]chan bool
var nextid int64

func setInterval(call otto.FunctionCall, conv contract.Converter) otto.Value {
	if len(call.ArgumentList) < 2 {
		return otto.Value{}
	}
	timelen, err := call.Argument(1).ToInteger()
	if err != nil {
		return otto.Value{}
	}
	obj := call.Argument(0)
	if clear == nil {
		clear = make(map[int64]chan bool)
	}
	id := nextid
	clear[id] = make(chan bool)
	nextid++

	go func() {
		for call.Otto != nil {
			select {
			case <-clear[id]:
				return
			default:
				timeout := time.NewTimer(time.Duration(int64(time.Millisecond) * timelen))
				<-timeout.C
				switch {
				case obj.IsString():
					runme, err := obj.ToString()
					if err != nil {
						return
					}
					_, err = call.Otto.Run(runme)
					if err != nil {
						return
					}
				case obj.IsFunction():
					_, err := obj.Call(obj)
					if err != nil {
						return
					}
				default:
					return
				}
			}

		}
	}()
	returnid, _ := conv.ToValue(id)
	return returnid
}
func clearInterval(call otto.FunctionCall, conv contract.Converter) otto.Value {
	if len(call.ArgumentList) == 0 {
		for _, key := range clear {
			key <- true
		}
	} else {
		key, err := call.Argument(0).ToInteger()
		if err != nil {
			return otto.Value{}
		}
		channel, ok := clear[key]
		if ok == true {
			channel <- true
		}
	}
	return otto.Value{}
}
