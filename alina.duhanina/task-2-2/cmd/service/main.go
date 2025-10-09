package main

import (
	"fmt"
	"log"
	"os"

	"dish-preference/heap"
	"dish-preference/service"
)

var (
	ErrEmptyDishes      = errors.New("список блюд не может быть пустым")
	ErrInvalidDishCount = errors.New("количество блюд должно быть от 1 до 10000")
	ErrInvalidK         = errors.New("k должно быть в диапазоне от 1 до N")
	ErrInvalidRating    = errors.New("рейтинг блюда должен быть в диапазоне от -10000 до 10000")
)

func ValidateInput(N int, dishes []int, k int) error {
	if len(dishes) == 0 {
		return ErrEmptyDishes
	}
	if N < 1 || N > 10000 {
		return ErrInvalidDishCount
	}
	if len(dishes) != N {
		return fmt.Errorf("ожидалось %d блюд, получено %d", N, len(dishes))
	}
	if k < 1 || k > N {
		return ErrInvalidK
	}
	for i, rating := range dishes {
		if rating < -10000 || rating > 10000 {
			return fmt.Errorf("недопустимый рейтинг блюда %d: %d (должен быть от -10000 до 10000)", i+1, rating)
		}
	}
	return nil
}

func FindKthPreference(dishes []int, k int) (int, error) {
	if err := ValidateInput(len(dishes), dishes, k); err != nil {
		return 0, err
	}

	h := &heap.IntHeap{}
	heap.Init(h)
	for _, dish := range dishes {
		heap.Push(h, dish)
	}
	for i := 0; i < k-1; i++ {
		heap.Pop(h)
	}
	result := heap.Pop(h).(int)
	return result, nil
}

func main() {
	if err := run(); err != nil {
		log.Printf("Ошибка: %v", err)
		os.Exit(1)
	}
}

func run() error {
	var N, k int
	_, err := fmt.Scan(&N)
	if err != nil {
		return fmt.Errorf("ошибка чтения количества блюд: %w", err)
	}
	dishes := make([]int, N)
	for i := 0; i < N; i++ {
		_, err := fmt.Scan(&dishes[i])
		if err != nil {
			return fmt.Errorf("ошибка чтения рейтинга блюда %d: %w", i+1, err)
		}
	}
	_, err = fmt.Scan(&k)
	if err != nil {
		return fmt.Errorf("ошибка чтения значения k: %w", err)
	}
	result, err := service.FindKthPreference(dishes, k)
	if err != nil {
		return fmt.Errorf("ошибка обработки данных: %w", err)
	}
	fmt.Println(result)
	return nil
}
