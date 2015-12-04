package contract

import (
	"github.com/robertkrimen/otto"
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
