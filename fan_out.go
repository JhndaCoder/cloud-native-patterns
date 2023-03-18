package main

func Split(source <-chan int, n int) []<-chan int {
	dests := make([]<-chan int, 0)
	for i := 0; i < n; i++ {
		ch := make(chan int)
		dests = append(dests, ch)

		go func() {
			defer close(ch)
			for val := range source {
				ch <- val
			}
		}()
	}
	return dests
}

// func main() {
// 	source := make(chan int)
// 	dests := Split(source, 5)

// 	go func() {
// 		for i := 0; i <= 10; i++ {
// 			source <- i
// 		}
// 		close(source)
// 	}()

// 	wg := &sync.WaitGroup{}
// 	wg.Add(len(dests))

// 	for i, ch := range dests {
// 		go func(i int, d <-chan int) {
// 			defer wg.Done()

// 			for val := range d {
// 				fmt.Printf("#%d got %d\n", i, val)
// 			}
// 		}(i, ch)
// 	}
// 	wg.Wait()
// }
