package main

import (
	"container/heap"
	"fmt"

	"github.com/TvoyBatyA12343/task-2-2/internal/intheap"
)

func fillHeap(heapToFill *intheap.IntHeap, cnt int) error {
	var dishRating int
	for range cnt {
		_, err := fmt.Scan(&dishRating)
		if err != nil {
			return fmt.Errorf("failed to read dish rating: %w", err)
		}

		heap.Push(heapToFill, dishRating)
	}

	return nil
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recover", r)
		}
	}()

	var dishesCnt int

	_, err := fmt.Scan(&dishesCnt)
	if err != nil {
		fmt.Printf("invalid dishes count: %v\n", err)

		return
	}

	dishesHeap := &intheap.IntHeap{}
	heap.Init(dishesHeap)

	err = fillHeap(dishesHeap, dishesCnt)
	if err != nil {
		fmt.Printf("failed to fill heap: %v\n", err)
	}

	var desire int

	_, err = fmt.Scan(&desire)
	if err != nil {
		fmt.Printf("failed to read desire: %v\n", err)

		return
	}

	res, err := dishesHeap.GetNth(desire)
	if err != nil {
		fmt.Printf("failed getting desired dish: %v\n", err)

		return
	}

	fmt.Println(res)
}
