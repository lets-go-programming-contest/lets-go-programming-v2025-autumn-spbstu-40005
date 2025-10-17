package department

import (
	"errors"
	"fmt"

	"github.com/DariaKhokhryakova/task-2-1/internal/employee"
)

var ErrInput = errors.New("invalid input")

func ProcessDepartment(countDepartment int) error {
	for range countDepartment {
		err := ProcessSingleDepartment()
		if err != nil {
			return fmt.Errorf("process single department: %w", err)
		}
	}

	return nil
}

func ProcessSingleDepartment() error {
	var countEmployees int

	_, err := fmt.Scan(&countEmployees)
	if err != nil {
		return fmt.Errorf("read employee count: %w", err)
	}

	if countEmployees < 1 || countEmployees > 1000 {
		return ErrInput
	}

	err := employee.ProcessEmployee(countEmployees)
	if err != nil {
		return fmt.Errorf("process employee: %w", err)
	}

	return nil
}
