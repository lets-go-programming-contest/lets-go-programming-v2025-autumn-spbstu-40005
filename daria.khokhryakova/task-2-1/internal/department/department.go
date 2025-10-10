package department

import (
	"errors"
	"fmt"

	"github.com/DariaKhokhryakova/task-2-1/internal/employee"
)

var ErrInput = errors.New("invalid input")

func ProcessDepartment() ([]int, error) {
	var countDepartment int

	_, err := fmt.Scan(&countDepartment)
	if err != nil || countDepartment < 1 {
		return nil, ErrInput
	}

	results := []int{}

	for range countDepartment {
		departResult, err := ProcessSingleDepartment()
		if err != nil {
			return nil, fmt.Errorf("process single department: %w", err)
		}

		results = append(results, departResult...)
	}

	return results, nil
}

func ProcessSingleDepartment() ([]int, error) {
	var countEmployees int

	_, err := fmt.Scan(&countEmployees)
	if err != nil {
		return nil, fmt.Errorf("read employee count: %w", err)
	}

	if countEmployees < 1 || countEmployees > 1000 {
		return nil, ErrInput
	}

	result, err := employee.ProcessEmployee(countEmployees)
	if err != nil {
		return nil, fmt.Errorf("process employee: %w", err)
	}

	return result, nil
}
