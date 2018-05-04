package pool

import (
	"fmt"
	"testing"
	"time"
)

func testFunc(x int, y *int) int {
	fmt.Printf("%d\n", x)
	*y = x + *y
	return x
}

func testSleep(x, final int, y *int) int {
	fmt.Printf("%d\n", x)
	*y = x + *y
	for *y != final {
		time.Sleep(1 * time.Second)
	}
	return x
}

func TestPool(t *testing.T) {

	pool := NewRoutinePool(5)

	var result int
	for i := 0; i < 20; i++ {
		pool.Run(testFunc, i, &result)
	}

	pool.Done()
	if result != 190 {
		t.Error("invalid num of result", result)
	}
}


func TestDefaultPool(t *testing.T) {

	pool := NewDefaultRoutinePool()

	var result int
	for i := 0; i < 20; i++ {
		pool.Run(testFunc, i, &result)
	}

	pool.Done()
	if result != 190 {
		t.Error("invalid num of result", result)
	}
}

func TestRunning(t *testing.T) {

	pool := NewRoutinePool(30)

	var result int
	for i := 0; i < 19; i++ {
		pool.Run(testSleep, i, 190, &result)
	}

	fmt.Println("hello: ", pool.Running())

	if pool.Running() != 19 {
		t.Error("invalid num of running routines", pool.Running)
	}
	pool.Run(testSleep, 19, 190, &result)

	pool.Done()
	if result != 190 {
		t.Error("invalid num of result", result)
	}
}
