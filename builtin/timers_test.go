package builtin

import (
	//"github.com/robertkrimen/otto"
	"github.com/vijayee/Contract"
	"log"
	"testing"
	"time"
)

func TestSetInterval(t *testing.T) {
	initVm()
	contract.Register(SetInterval)
	err := contract.LoadAll(vm)
	if err != nil {
		t.Error(err)
	}
	script := `
		var result = 0;
		setInterval("result++",200);
	`
	_, err = vm.Run(script)
	if err != nil {
		t.Error(err)
	}

	timeout := time.NewTimer(time.Duration(time.Millisecond * 1500))
	<-timeout.C

	value, err := vm.Get("result")
	if err != nil {
		t.Error(err)
	}

	result, _ := value.ToInteger()
	log.Printf("result: %d", result)

	if result == 0 {
		t.Error("Failed to run a interval Callback")
	}
}
func TestClearInterval(t *testing.T) {
	initVm()
	contract.Register(SetInterval)
	contract.Register(ClearInterval)
	err := contract.LoadAll(vm)
	if err != nil {
		t.Error(err)
	}
	script := `
		var result = 0;
		var id= setInterval("result++",3000);
		console.log(id);
		clearInterval(id);		
	`
	_, err = vm.Run(script)
	if err != nil {
		t.Error(err)
	}

	timeout := time.NewTimer(time.Duration(time.Millisecond * 4000))
	<-timeout.C

	value, err := vm.Get("result")
	if err != nil {
		t.Error(err)
	}

	result, _ := value.ToInteger()
	log.Printf("result: %d", result)

	if result != 0 {
		t.Error("Failed to clear interval Callback")
	}
}
func TestSetTimeout(t *testing.T) {
	initVm()
	contract.Register(SetTimeout)
	err := contract.LoadAll(vm)
	if err != nil {
		t.Error(err)
	}
	script := `
		var result = 0;
		setTimeout("result++",1000);
	`
	_, err = vm.Run(script)
	if err != nil {
		t.Error(err)
	}

	timeout := time.NewTimer(time.Duration(time.Millisecond * 1500))
	<-timeout.C

	value, err := vm.Get("result")
	if err != nil {
		t.Error(err)
	}

	result, _ := value.ToInteger()
	log.Printf("result: %d", result)
	if result == 0 {
		t.Error("Failed to run a timeout Callback")
	}
}
