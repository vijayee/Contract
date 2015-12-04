package contract

import (
	"github.com/bamzi/jobrunner"
	"github.com/robertkrimen/otto"
	//"log"
	"sync"
	"time"
)

type Contract struct {
	parties []string
	flow    Workflow
	apis    []string
}
type Workflow struct {
	vm             *otto.Otto
	tasks          []Task
	started        bool
	completed      bool
	completionDate time.Time
	startDate      time.Time
	currentTask    int
	errors         []error
	sync.Mutex
}

func NewWorkflow() Workflow {
	return Workflow{}
}

func (w *Workflow) Add(tsk Task) {
	if w.completed == false {
		w.Lock()
		w.tasks = append(w.tasks, tsk)
		w.Unlock()
	}
}
func (w *Workflow) Run() {
	for w.Done() != true {
		w.ExecuteNext()
	}
}

func (w *Workflow) SetVm(vm *otto.Otto) {
	w.vm = vm
}

func (w *Workflow) Length() int {
	return len(w.tasks)
}

func (w *Workflow) Started() bool {
	return w.started
}

func (w *Workflow) Done() bool {
	return w.completed
}

func (w *Workflow) Current() int {
	return w.currentTask
}

func (w *Workflow) ExecuteNext() { //blocking
	if w.vm == nil {
		return
	}
	if w.Length() > 0 && w.completed == false {
		w.Lock()
		if w.currentTask == 0 {
			w.started = true
			w.startDate = time.Now()
		}
		current := w.tasks[w.currentTask]
		jobrunner.Start()
		current.SetVm(w.vm)
		var wg sync.WaitGroup
		current.SetWaitGroup(&wg)
		wg.Add(1)
		/*
			TODO: add Time Lapse scheduling
		*/
		if current.isScheduled() == true {
			jobrunner.Schedule(current.schedule, &current)
		} else {
			jobrunner.Now(&current)
		}
		if current.err != nil {
			w.errors = append(w.errors, current.err)
		}
		wg.Wait()
		w.currentTask++

		if w.currentTask >= w.Length() {
			w.completed = true
			w.completionDate = time.Now()
		}
		w.Unlock()
	}

}

type Task struct {
	vm             *otto.Otto
	err            error
	apis           []string
	code           string
	completed      bool
	completionDate time.Time
	startDate      time.Time
	schedule       string
	wg             *sync.WaitGroup
	sync.Mutex
}

func NewTask() Task {
	return Task{}
}

func (t *Task) Done() bool {
	return t.completed
}

func (t *Task) SetVm(vm *otto.Otto) {
	t.vm = vm
}

func (t *Task) requireAPI(name string) {
	for _, a := range t.apis {
		if a == name {
			return
		}
	}
	t.apis = append(t.apis, name)
}

func (t *Task) SetSchedule(schedule string) {
	t.schedule = schedule
}

func (t *Task) isScheduled() bool {
	return t.schedule != ""
}

func (t *Task) SetWaitGroup(wg *sync.WaitGroup) {
	t.wg = wg
}

func (t *Task) Run() {
	if t.vm != nil {
		t.Lock()
		if err := LoadSet(t.apis, t.vm); err != nil {
			t.err = err
			t.Unlock()
			return
		}
		t.startDate = time.Now()
		t.vm.Run(buildModule(t.code, t.apis))
		t.completionDate = time.Now()
		t.completed = true
		t.Unlock()
		if t.wg != nil {
			t.wg.Done()
		}
	}
}

func (t *Task) Load(code string) {
	if t.completed == false {
		t.Lock()
		t.code = code
		t.Unlock()
	}

}
