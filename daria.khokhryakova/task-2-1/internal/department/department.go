package department

import (
	"errors"
	"fmt"

	"github.com/DariaKhokhryakova/task-2-1/internal/employee"
)

var ErrorInput = errors.New("Invalid input")

func ProcessDepartment() ([]int, error) {
	var countDepartment int
	_, err := fmt.Scan(&countDepartment)
	if err != nil || countDepartment < 1 {
		return nil, ErrorInput
	}

	results := []int{}
	for i := 0; i < countDepartment; i++ {
		departResult, err := ProcessSingleDepartment()
		if err != nil {
			return nil, err
		}
		for _, value := range departResult {
			results = append(results, value)
		}
	}
	return results, nil
}

func ProcessSingleDepartment() ([]int, error) {
	var countEmployees int
	_, err := fmt.Scan(&countEmployees)
	if err != nil || countEmployees < 1 || countEmployees > 1000 {
		return nil, ErrorInput
	}
	return employee.ProcessEmployee(countEmployees)
}
