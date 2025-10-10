package main

import (
	"fmt"

	"github.com/DariaKhokhryakova/task-2-1/internal/department"
)

func main() {
	itog, err := department.ProcessDepartment()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, temp := range itog {
		fmt.Println(temp)
	}
}
