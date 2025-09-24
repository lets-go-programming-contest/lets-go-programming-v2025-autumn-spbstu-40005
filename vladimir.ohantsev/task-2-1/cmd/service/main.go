package main

import (
	"errors"
	"fmt"

	"github.com/P3rCh1/task-2-1/internal/department"
)

type comparisonOperator string

const (
	LessEqual    comparisonOperator = "<="
	GreaterEqual comparisonOperator = ">="
)

var errInvalidOperation = errors.New("invalid operation")

func processWorker(dep *department.Department) error {
	var operator comparisonOperator
	if _, err := fmt.Scan(&operator); err != nil {
		return fmt.Errorf("failed to scan operator: %w", err)
	}

	var temp int
	if _, err := fmt.Scan(&temp); err != nil {
		return fmt.Errorf("failed to scan temperature: %w", err)
	}

	switch operator {
	case LessEqual:
		dep.SetMax(temp)
	case GreaterEqual:
		dep.SetMin(temp)
	default:
		return errInvalidOperation
	}

	return nil
}

func processDepartment() error {
	var countWorkers int
	if _, err := fmt.Scan(&countWorkers); err != nil {
		return fmt.Errorf("failed to scan count workers: %w", err)
	}

	dep := department.New()
	for range countWorkers {
		if err := processWorker(dep); err != nil {
			return fmt.Errorf("process worker fail: %w", err)
		}

		fmt.Println(dep.Optimum())
	}

	return nil
}

func main() {
	var countDepartments int
	if _, err := fmt.Scan(&countDepartments); err != nil {
		fmt.Printf("failed to scan count departments: %s\n", err)

		return
	}

	for range countDepartments {
		if err := processDepartment(); err != nil {
			fmt.Printf("process department fail: %s\n", err)

			return
		}
	}
}
