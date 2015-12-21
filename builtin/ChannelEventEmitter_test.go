package builtin

import (
	"github.com/robertkrimen/otto"
	"github.com/vijayee/Contract"
	//"log"
	"os"
	"sync"
	"testing"
)

var script string
var vm *otto.Otto

func initVm() {
	vm = otto.New()
	vm.Set("result", 0)
}

func TestMain(m *testing.M) {
	script = `
	(function(){
		var a= new EventEmitter;
		a.on('test', function(msg){console.log(msg);});		
		a.emit('test','omg did it work?');
	})();	`
	os.Exit(m.Run())

}
func TestBroadcast(t *testing.T) {
	initVm()
	contract.Register(ChannelEventEmitter)
	contract.Register(ChannelEventEmitterAnnouncer)
	err := contract.LoadAll(vm)
	if err != nil {
		t.Error(err)
	}
	var received int64
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		received = <-Broadcast
		wg.Done()
	}()

	_, err = vm.Run(script)
	if err != nil {
		t.Error(err)
	}
	wg.Wait()
	if received != 1 {
		t.Error("Unexpected Execution results")
	}
}

func TestEvent(t *testing.T) {
	initVm()
	contract.Register(ChannelEventEmitter)
	contract.Register(ChannelEventEmitterAnnouncer)
	contract.Register(SetTimeout)
	err := contract.LoadAll(vm)

	script = `
	(function(){
		var a= new EventEmitter;
		var meh = function(){		
			a.emit('test','omg did it work?');
		};
		setTimeout(meh,3000);		
	})();`
	if err != nil {
		t.Error(err)
	}
	var e Event

	echan := make(chan Event)
	Subscribe(1, "test", echan)
	var wg sync.WaitGroup
	wg.Add(1)
	//
	go func() {
		_, err = vm.Run(script)
		if err != nil {
			t.Error(err)
		}
		wg.Done()
	}()
	e = <-echan

	wg.Wait()

	if &e == nil {
		t.Error("Unexpected Execution results")
	}

}
