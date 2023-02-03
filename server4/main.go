package main

import (
	"fmt"
	"math/rand"
	"time"
)

const MAX = 5

func main() {
	sem := make(chan int, MAX)
	count := 0
	for {
		sem <- count // will block if there is MAX ints in sem
		count++

		go func() {
			rand := 100 + rand.Intn(300)
			t := fmt.Sprintf("%dms", rand)
			ms, _ := time.ParseDuration(t)

			time.Sleep(ms)

			a, ok := <-sem
			if !ok {
				fmt.Println("error")
			}
			fmt.Printf("hello again, world: %d : %s\n", a, t)
		}()
	}
}
