package main

import (
	"container/heap"
	"errors"
	"fmt"

	"alexandra.karnauhova/task-2-2/internal/queue"
)

var ErrInvalidDishSelection = errors.New("invalid dish selection")

func readMenu(countDish int) (queue.Queue, error) {
	menu := make(queue.Queue, 0)

	heap.Init(&menu)

	for range countDish {
		var estimation int

		_, err := fmt.Scan(&estimation)
		if err != nil {
			return menu, fmt.Errorf("invalid estimation %w", err)
		}

		heap.Push(&menu, estimation)
	}

	return menu, nil
}

func readPreference(countDish int) int {
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

func chooseDish(menu queue.Queue, countDish int) (int, error) {
	kValue := readPreference(countDish)

	if kValue == 0 {
		return 0, fmt.Errorf("%w: k value cannot be zero", ErrInvalidDishSelection)
	}

	for range kValue - 1 {
		heap.Pop(&menu)
	}

	dishItem := heap.Pop(&menu)

	dish, ok := dishItem.(int)
	if !ok {
		return 0, fmt.Errorf("%w: expected int dish type", ErrInvalidDishSelection)
	}

	return dish, nil
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

	res, err := chooseDish(menu, countDish)
	if err != nil {
		fmt.Printf("Error choosing a dish: %v\n", err)

		return
	}

	fmt.Println(res)
}
