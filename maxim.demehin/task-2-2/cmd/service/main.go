package main

import (
	"container/heap"
	"errors"
	"fmt"

	"github.com/TvoyBatyA12343/task-2-2/internal/intheap"
)

var errInput = errors.New("input error")

func fillHeap(heapToFill *intheap.IntHeap, cnt int) error {
	var dishRating int
	for range cnt {
		_, err := fmt.Scan(&dishRating)
		if err != nil {
			return errInput
		}

		heap.Push(heapToFill, dishRating)
	}

	return nil
}

func main() {
	var dishesCnt int

	_, err := fmt.Scan(&dishesCnt)
	if err != nil {
		fmt.Println(errInput.Error())

		return
	}

	dishesHeap := &intheap.IntHeap{}
	heap.Init(dishesHeap)

	err = fillHeap(dishesHeap, dishesCnt)
	if err != nil {
		fmt.Println(err.Error())
	}

	var desire int

	_, err = fmt.Scan(&desire)
	if err != nil {
		fmt.Println(errInput.Error())

		return
	}

	res := dishesHeap.GetNth(desire)
	fmt.Println(res)
}
