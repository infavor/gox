# package pool

A pool allows a particular type of work to run in it and limits 
its parallel number and maximum wait queue length.

Usage:

```golang
import "github.com/hetianyi/gox/pool"

// init a task pool which allow 2 tasks run in parallel,
// and max wait task count is 100.
p := pool.New(1, 100)
for i := 0; i < 100; i++ {
    tmp := i
    err := p.Push(func() {
        myfunc(...){}
    })
    if err != nil {
        fmt.Println("Err:", err)
    }
}
```