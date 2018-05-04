# go-routine-pool



A way to limit the parallelism while running multiple tasks.

Under the hood it uses Go routines.



### Usage

```go
package main

import (
    pool "github.com/zmalik/go-routine-pool"
)

func parallelTask(x int, y *int) int {
	fmt.Printf("%d\n", x)
	*y = x + *y
	return x
}


func main() {
    routinePool := pool.NewDefaultRoutinePool()

    var result int
    for i := 0; i < 20; i++ {
        routinePool.Run(parallelTask, i, &result)
    }
    fmt.Println(result)
}
```

### Limiting number of Go routines

```go
routinePool := pool.NewRoutinePool(5)
```



## Why?

Power is nothing without control 

Using Go routines one can easily rage the power of routines without limiting how many parallel tasks are running. 

This library helps to control that in easy way