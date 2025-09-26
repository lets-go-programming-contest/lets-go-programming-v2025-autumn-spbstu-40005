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

func findLargest(nums []int, k int) int {
	dishesHeap := &maxheap.MaxHeap{}
	heap.Init(dishesHeap)

	for _, num := range nums {
		heap.Push(dishesHeap, num)
	}

	for range k - 1 {
		heap.Pop(dishesHeap)
	}

	return (*dishesHeap)[0]
}

func main() {
	var dishesCount int
	if _, err := fmt.Scan(&dishesCount); err != nil {
		fmt.Println(errInput.Error())

		return
	}

	if dishesCount < 1 || dishesCount > 10000 {
		fmt.Println(errNumber.Error())

		return
	}

	nums := make([]int, dishesCount)
	for i := range dishesCount {
		if _, err := fmt.Scan(&nums[i]); err != nil {
			fmt.Println(errInput.Error())

			return
		}

		if nums[i] < -10000 || nums[i] > 10000 {
			fmt.Println(errNumber.Error())

			return
		}
	}

	var k int
	if _, err := fmt.Scan(&k); err != nil {
		fmt.Println(errInput.Error())

		return
	}

	if k < 1 || k > dishesCount {
		fmt.Println(errNumber)

		return
	}

	fmt.Println(findLargest(nums, k))
}
