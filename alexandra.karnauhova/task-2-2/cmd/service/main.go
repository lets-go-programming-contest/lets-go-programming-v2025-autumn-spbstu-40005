package main

import (
	"container/heap"
	"fmt"

	"alexandra.karnauhova/task-2-2/internal/queue"
)

func readMenu(countDish int) (queue.Queue, error) {
	menu := make(queue.Queue, 0)

	heap.Init(&menu)

	for range countDish {
		var estimation int

		_, err := fmt.Scan(&estimation)
		if err != nil {
			return menu, err
		}

		heap.Push(&menu, estimation)
	}

	return menu, nil
}

func choosePreference(countDish int) int {
	var kValue int

	_, err := fmt.Scan(&kValue)
	if err != nil {
		return 0
	}

	if kValue <= 0 || kValue > countDish {
		return 0
	}

	return kValue
}

func chooseDish(menu queue.Queue, countDish int) int {
	kValue := choosePreference(countDish)

	if kValue == 0 {
		fmt.Println("Invalid k value")

		return 0
	}

	dishItem := heap.Pop(&menu)

	dish, oke := dishItem.(int)
	if !oke {
		fmt.Println("Invalid dish type")

		return 0
	}

	for range kValue - 1 {
		dishItem = heap.Pop(&menu)

		dish, oke = dishItem.(int)
		if !oke {
			fmt.Println("Invalid dish type")

			return 0
		}
	}

	return dish
}

func main() {
	var countDish int

	_, err := fmt.Scan(&countDish)
	if err != nil {
		fmt.Printf("Invalid count dish: %v\n", err)

		return
	}

	menu, err := readMenu(countDish)
	if err != nil {
		fmt.Printf("Error reading menu: %v\n", err)

		return
	}

	res := chooseDish(menu, countDish)

	if res == 0 {
		fmt.Println("Error choosing a dish")

		return
	} else {
		fmt.Println(res)
	}
}
