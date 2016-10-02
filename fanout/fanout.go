package main

import (
	"fmt"
	"sync"
)

type inmessage struct {
	line   string
	number int
}

type outmessage struct {
	line   string
	numstr string
}

func main() {
	out := make(chan outmessage)

	go print(out)

	c := gen(2, 3, 4, 5, 6)

	c1 := sq(1, c)
	c2 := sq(2, c)

	for n := range merge(c1, c2) {
		//fmt.Println(n)
		out <- n
	}

	close(out)
}

func sq(number int, in <-chan inmessage) <-chan outmessage {
	out := make(chan outmessage)

	go func() {
		for n := range in {
			outmsg := outmessage{
				line:   fmt.Sprintf("Sq (%v): %v", number, n.line),
				numstr: fmt.Sprintf("Sq (%v): numstr: %v", number, n.number),
			}
			out <- outmsg
		}
		close(out)
	}()

	return out
}

func gen(nums ...int) <-chan inmessage {
	out := make(chan inmessage)

	go func() {
		for _, n := range nums {
			msg := inmessage{
				line:   fmt.Sprintf("Line (%v)", n),
				number: n,
			}
			out <- msg
		}
		close(out)
	}()

	return out
}

func print(in <-chan outmessage) {
	for n := range in {
		fmt.Printf("%v: %v\n", n.line, n.numstr)
	}
}

func merge(cs ...<-chan outmessage) <-chan outmessage {
	var wg sync.WaitGroup

	out := make(chan outmessage)

	output := func(c <-chan outmessage) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
