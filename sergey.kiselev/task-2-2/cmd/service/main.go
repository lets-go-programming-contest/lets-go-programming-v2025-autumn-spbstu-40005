package main

import (
	"container/heap"
	"fmt"

	"sergey.kiselev/task-2-2/internal/maxheap"
)

func findLargest(nums []int, number int) int {
	dishesHeap := &maxheap.MaxHeap{}
	heap.Init(dishesHeap)

	for _, num := range nums {
		heap.Push(dishesHeap, num)
	}

	for range number - 1 {
		heap.Pop(dishesHeap)
	}

	return (*dishesHeap)[0]
}

func main() {
	var dishesCount int
	if _, err := fmt.Scan(&dishesCount); err != nil {
		fmt.Printf("invalid input for dishesCount: %v\n", err)

		return
	}

	nums := make([]int, dishesCount)
	for index := range dishesCount {
		if _, err := fmt.Scan(&nums[index]); err != nil {
			fmt.Printf("invalid input for nums: %v\n", err)

			return
		}
	}

	var number int
	if _, err := fmt.Scan(&number); err != nil {
		fmt.Printf("invalid input for number: %v\n", err)

		return
	}

	if number < 1 || number > dishesCount {
		fmt.Println("this number is not suitable")

		return
	}

	fmt.Println(findLargest(nums, number))
}
