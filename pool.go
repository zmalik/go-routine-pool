package pool

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
)

const (
	DEFAULT_ROUTINE_NUMBERS = 256
)

type RoutinePool struct {
	parallelism int
	routine     chan int
	running     int64
	runningLock *sync.Mutex
}

// Limit of go routines running in parallel
func NewDefaultRoutinePool() *RoutinePool {
	return NewRoutinePool(DEFAULT_ROUTINE_NUMBERS)
}

// Limit of go routines running in parallel
func NewRoutinePool(parallelism int) *RoutinePool {
	if parallelism <= 0 {
		panic("Parallelism cannot be negative.")
	}
	c := &RoutinePool{
		parallelism: parallelism,
		routine:     make(chan int, parallelism),
		runningLock: &sync.Mutex{},
	}

	for i := 0; i < c.parallelism; i++ {
		c.routine <- i
	}

	return c
}

func (c *RoutinePool) Run(fn interface{}, args ...interface{}) int {
	numRoutine := <-c.routine
	c.runTask(numRoutine, fn, args)
	return numRoutine
}

func (c *RoutinePool) runTask(numberRoutine int, fn interface{}, args []interface{}) {
	atomic.AddInt64(&c.running, 1)
	go func() {
		defer func() {
			c.routine <- numberRoutine
			atomic.AddInt64(&c.running, -1)

		}()

		// run the task
		f := reflect.ValueOf(fn)
		t := f.Type()

		if t.Kind() != reflect.Func {
			panic(fmt.Errorf("expected func, got: %v", t))
		}

		if t.NumIn() > len(args) {
			panic(fmt.Errorf("expected %d args, got %d args", len(args), t.NumIn()))
		}

		inputs := make([]reflect.Value, len(args))
		for i := 0; i < len(args); i++ {
			if args[i] == nil {
				inputs[i] = reflect.Zero(f.Type().In(i))
			} else {
				inputs[i] = reflect.ValueOf(args[i])
			}
		}
		f.Call(inputs)
	}()
}

func (c *RoutinePool) Done() {
	for i := 0; i < c.parallelism; i++ {
		_ = <-c.routine
	}
}

func (c *RoutinePool) Running() int64 {
	c.runningLock.Lock()
	defer c.runningLock.Unlock()
	return c.running
}
