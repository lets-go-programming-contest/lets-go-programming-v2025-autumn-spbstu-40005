package main

import (
	"container/heap"
	"fmt"

	"alexandra.karnauhova/task-2-2/internal/queue"
)

func returnNumberDish() int {
	var countDish int
	menu := make(queue.Queue, 0)
	heap.Init(&menu)
	_, err := fmt.Scan(&countDish)
	if err != nil {
		fmt.Println("Invalid count dish")
		return 0
	}
	for i := 0; i < countDish; i++ {
		var estimation int
		_, err = fmt.Scan(&estimation)
		if err != nil {
			fmt.Println("Invalid estimation")
			return 0
		}
		heap.Push(&menu, &queue.Element{Value: estimation, Priority: estimation})
	}
	var k int
	_, err = fmt.Scan(&k)
	if err != nil {
		fmt.Println("Invalid k")
		return 0
	}
	dish := heap.Pop(&menu).(*queue.Element)
	for i := 0; i < k-1; i++ {
		dish = heap.Pop(&menu).(*queue.Element)
	}
	return dish.Value
}

func main() {
	res := returnNumberDish()
	if res == 0 {
		fmt.Errorf("its bad")
		return
	} else {
		fmt.Println(res)
	}
}
