// We can use channels to synchronize execution
// across goroutines. Here's an example of using a
// blocking receive to wait for a goroutine to finish.

package main

import "fmt"
import "time"

// This is the function we'll run in a goroutine. The
// `done` channel will be used to notify another
// goroutine that this function's work is done.
func worker(n int, done chan bool) {
	fmt.Println(n, "working...")
	time.Sleep(time.Second)
	fmt.Println(n, "done")

	// Send a value to notify that we're done.
	done <- true
}

func main() {

	go func() {
		c := time.Tick(100 * time. Millisecond)
		for now := range c {
			fmt.Printf("tick %v\n", now)
		}
	}()

	// Start a worker goroutine, giving it the channel to
	// notify on.
	done := make(chan bool, 1)

	for i := 0; i < 3; i++ {
		go worker(i, done)

		// block until message received
		<-done
	}

}
