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

func fillHeap(h *int_heap.IntHeap, cnt int) error {
	var dishRating int
	for range cnt {
		_, err := fmt.Scan(&dishRating)
		if err != nil {
			return err
		}

		heap.Push(h, dishRating)
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

	h := &int_heap.IntHeap{}
	heap.Init(h)

	err = fillHeap(h, dishesCnt)
	if err != nil {
		fmt.Println(err.Error())
	}

	var desire int
	_, err = fmt.Scan(&desire)
	if err != nil {
		fmt.Println(errInput.Error())

		return
	}

	res := h.GetNth(desire)
	fmt.Println(res)
}
