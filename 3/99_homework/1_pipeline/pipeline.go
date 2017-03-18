package pipeline

import "sync"

type job func(in, out chan interface{})

func Pipe(funcs ...job) {
	var wg sync.WaitGroup
	wg.Add(len(funcs))
	var channels []chan interface{}
	for i := 0; i < len(funcs)+2; i++ {
		channels = append(channels, make(chan interface{}))
	}
	for i, fun := range funcs {
		go func(i int, f job) {
			curIn := channels[i]
			curOut := channels[i+1]
			defer wg.Done()
			defer close(curOut)
			f(curIn, curOut)
		}(i, fun)
	}
	wg.Wait()
	return
}
