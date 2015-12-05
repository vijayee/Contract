package contract

import (
	"github.com/robertkrimen/otto"
	"log"
	"os"
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
		console.log("This is a task that is running");
		result++;
	})();	`
	os.Exit(m.Run())

}
func TestTask(t *testing.T) {
	initVm()
	tsk := NewTask()
	tsk.SetVm(vm)
	tsk.Load(script)
	tsk.Run()
	value, err := vm.Get("result")
	if err != nil {
		t.Error(err)
	}

	if result, _ := value.ToInteger(); result != 1 {
		t.Error("Unexpected Execution results")
	}

}

func TestWorkflow(t *testing.T) {
	initVm()
	wrkflw := NewWorkflow()
	var i int64
	for i = 0; i < 3; i++ {
		tsk := NewTask()
		tsk.Load(script)
		wrkflw.Add(tsk)
	}
	wrkflw.SetVm(vm)
	wrkflw.Run()

	t.Logf("count: %d\n", wrkflw.Length())
	t.Logf("current taks: %d\n", wrkflw.Current())
	t.Logf("done: %t\n", wrkflw.Done())

	if wrkflw.Done() != true {
		t.Error("Workflow not reporting completion")
	}
	value, err := vm.Get("result")
	if err != nil {
		t.Error(err)
	}
	result, _ := value.ToInteger()

	if result != i {
		t.Error("Unexpected Execution results")
	}
}

func TestAPI(t *testing.T) {
	initVm()
	testapi1 := NewApi("testApi1")
	testapi2 := NewApi("testApi2")
	testapi3 := NewApi("testApi3")

	testFunc1 := func(call otto.FunctionCall, conv Converter) otto.Value {
		log.Printf("This is test api one\n")
		return otto.Value{}
	}

	testFunc2 := func(call otto.FunctionCall, conv Converter) otto.Value {
		log.Printf("This is test api two\n")
		return otto.Value{}
	}

	testFunc3 := func(call otto.FunctionCall, conv Converter) otto.Value {
		log.Printf("This is test api three\n")
		return otto.Value{}
	}
	apiScript1 := `(function(){
		console.log("stuff is runnin'");
		testApi1();
		testApi2();
		result++;
	})();`

	apiScript2 := `(function(){
		console.log("stuff is runnin'");
		testApi1();
		testApi2();
		testApi3();
		result++;
	})();`

	testapi1.SetFunction(testFunc1)
	testapi2.SetFunction(testFunc2)
	testapi3.SetFunction(testFunc3)

	Register(testapi1)
	Register(testapi2)
	Register(testapi3)

	tsk1 := NewTask()
	tsk1.SetVm(vm)
	tsk1.Load(apiScript1)
	tsk1.requireAPI("testApi1", "testApi2")
	tsk1.Run()
	value, err := vm.Get("result")
	if err != nil {
		t.Error(err)
	}
	result, _ := value.ToInteger()

	if result != 1 {
		t.Error("Unexpected Execution results")
	}

	vm.Set("result", 0)
	tsk2 := NewTask()
	tsk2.SetVm(vm)
	tsk2.Load(apiScript2)
	tsk2.requireAPI("testApi1", "testApi2")
	tsk2.Run()

	value, err = vm.Get("result")
	if err != nil {
		t.Error(err)
	}

	result, _ = value.ToInteger()

	if result != 0 {
		t.Error("Failed to limit module scope")
	}
}
