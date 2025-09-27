package main

import (
	"container/heap"
	"errors"
	"fmt"

	"sergey.kiselev/task-2-2/internal/maxheap"
)

var (
	errInput  = errors.New("invalid input")
	errNumber = errors.New("this number is not suitable")
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
		fmt.Printf("%s: %v\n", errInput.Error(), err)

		return
	}

	if dishesCount < 1 || dishesCount > 10000 {
		fmt.Println(errNumber.Error())

		return
	}

	nums := make([]int, dishesCount)
	for index := range dishesCount {
		if _, err := fmt.Scan(&nums[index]); err != nil {
			fmt.Printf("%s: %v\n", errInput.Error(), err)

			return
		}

		if nums[index] < -10000 || nums[index] > 10000 {
			fmt.Println(errNumber.Error())

			return
		}
	}

	var number int
	if _, err := fmt.Scan(&number); err != nil {
		fmt.Printf("%s: %v\n", errInput.Error(), err)

		return
	}

	if number < 1 || number > dishesCount {
		fmt.Println(errNumber)

		return
	}

	fmt.Println(findLargest(nums, number))
}
