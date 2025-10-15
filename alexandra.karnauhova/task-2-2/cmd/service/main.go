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

	for range countDish {
		var estimation int

		_, err = fmt.Scan(&estimation)
		if err != nil {
			fmt.Println("Invalid estimation")

			return 0
		}
		heap.Push(&menu, &queue.Element{
			Value:    estimation,
			Priority: estimation,
			Index:    0,
		})
	}

	var kValue int

	_, err = fmt.Scan(&kValue)
	if err != nil {
		fmt.Println("Invalid k")

		return 0
	}

	if kValue <= 0 || kValue > countDish {
		fmt.Println("Invalid k value")

		return 0
	}

	dishItem := heap.Pop(&menu)

	dish, oke := dishItem.(*queue.Element)
	if !oke {
		fmt.Println("Invalid dish type")

		return 0
	}

	for range kValue - 1 {
		dishItem = heap.Pop(&menu)

		dish, oke = dishItem.(*queue.Element)
		if !oke {
			fmt.Println("Invalid dish type")

			return 0
		}
	}

	return dish.Value
}

func main() {
	res := returnNumberDish()

	if res == 0 {
		fmt.Println("its bad")

		return
	} else {
		fmt.Println(res)
	}
}
