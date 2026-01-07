package main

import (
	"container/heap"
	"errors"
	"fmt"

	"github.com/aleksey.kurbyko/task-2-2/internal/dishheap"
)

const (
	DishCountMin = 1
	DishCountMax = 10000

	RatingMin = -10000
	RatingMax = 10000
)

var (
	ErrIncorrectDishCount = errors.New("incorrect amount of dishes")
	ErrIncorrectRating    = errors.New("incorrect rating for the dish")
	ErrIncorrectK         = errors.New("incorrect k")
	ErrExpectedInt        = errors.New("expected int")
)

func readInt() (int, error) {
	var value int

	if _, err := fmt.Scan(&value); err != nil {
		return 0, fmt.Errorf("scan int: %w", err)
	}

	return value, nil
}

func run() error {
	dishCount, err := readInt()
	if err != nil {
		return ErrIncorrectDishCount
	}

	if dishCount < DishCountMin || dishCount > DishCountMax {
		return ErrIncorrectDishCount
	}

	ratingsHeap := &dishheap.DishHeap{}
	heap.Init(ratingsHeap)

	for range dishCount {
		rating, err := readInt()
		if err != nil {
			return ErrIncorrectRating
		}

		if rating < RatingMin || rating > RatingMax {
			return ErrIncorrectRating
		}

		heap.Push(ratingsHeap, rating)
	}

	k, err := readInt()
	if err != nil {
		return ErrIncorrectK
	}

	if k < 1 || k > dishCount {
		return ErrIncorrectK
	}

	for range k - 1 {
		heap.Pop(ratingsHeap)
	}

	resultAny := heap.Pop(ratingsHeap)
	result, ok := resultAny.(int)
	if !ok {
		return ErrExpectedInt
	}

	fmt.Println(result)

	return nil
}

func main() {
	if err := run(); err != nil {
		return
	}
}
