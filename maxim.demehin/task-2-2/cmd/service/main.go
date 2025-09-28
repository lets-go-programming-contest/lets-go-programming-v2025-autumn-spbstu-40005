package main

import (
	"container/heap"
	"errors"
	"fmt"

	"github.com/TvoyBatyA12343/task-2-2/internal/int_heap"
)

var (
	errInput = errors.New("input error")
)

func main() {
	var dishesCnt int
	_, err := fmt.Scan(&dishesCnt)
	if err != nil {
		fmt.Println(errInput.Error())

		return
	}

	h := &int_heap.IntHeap{}
	heap.Init(h)

	var dishRating int
	for range dishesCnt {
		_, err = fmt.Scan(&dishRating)
		if err != nil {
			fmt.Println(errInput.Error())

			return
		}

		heap.Push(h, dishRating)
	}

	var (
		desire int
		result int
	)

	_, err = fmt.Scan(&desire)
	if err != nil {
		fmt.Println(errInput.Error())

		return
	}

	for range desire {
		result = heap.Pop(h).(int)
	}

	fmt.Println(result)

}
