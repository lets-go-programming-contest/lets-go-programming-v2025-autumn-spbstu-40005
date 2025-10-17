package main

import (
	"fmt"

	"github.com/DariaKhokhryakova/task-2-1/internal/department"
)

func main() {
	var countDepartment int
	_, err := fmt.Scan(&countDepartment)

	if err != nil || countDepartment < 1 {
		fmt.Println("failed to read countDepartment:", err)

		return
	}

	err = department.ProcessDepartment(countDepartment)
	if err != nil {
		fmt.Println(err.Error())

		return
	}
}
