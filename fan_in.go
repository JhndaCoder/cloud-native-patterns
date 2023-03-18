package main

import (
	"sync"
)

func Funnel(sources ...<-chan int) <-chan int {
	dest := make(chan int)

	wg := &sync.WaitGroup{}

	wg.Add(len(sources))

	for _, ch := range sources {
		go func(c <-chan int) {
			defer wg.Done()
			for n := range c {
				dest <- n
			}
		}(ch)
	}
	go func() {
		wg.Wait()
		close(dest)
	}()
	return dest
}

// func main() {
// 	sources := make([]<-chan int, 0)
// 	for i := 0; i < 3; i++ {
// 		ch := make(chan int)
// 		sources = append(sources, ch)

// 		go func() {
// 			defer close(ch)
// 			for i := 0; i < 5; i++ {
// 				ch <- i
// 				time.Sleep(time.Second)
// 			}
// 		}()
// 	}

// 	dest := Funnel(sources...)
// 	for d := range dest {
// 		fmt.Println(d)
// 	}
// }
