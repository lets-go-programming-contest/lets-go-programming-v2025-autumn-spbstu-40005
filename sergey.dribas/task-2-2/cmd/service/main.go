package main

import (
	"fmt"

	"sergey.dribas/task-2-2/internal/intheap"
)

func main() {
	var (
		size, number, predict int
		dish                  = &intheap.IntHeap{}
	)

	if _, err := fmt.Scan(&size); err != nil {
		return
	}

	for range size {
		if _, err := fmt.Scan(&number); err != nil {
			return
		}

		dish.Push(number)
	}

	if _, err := fmt.Scan(&predict); err != nil {
		return
	}

	if result, err := intheap.FindKthSmallest(dish, predict); err != nil {
		fmt.Println(result)
	}
}
