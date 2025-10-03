package employee

import "sergey.kiselev/task-2-1/internal/temperature"

type Employee struct {
	operator    string
	temperature int
}

func New(oper string, temp int) *Employee {
	return &Employee{
		operator:    oper,
		temperature: temp,
	}
}
func (employee *Employee) Process(manager *temperature.TemperatureManager) (int, error) {
	if err := manager.Update(employee.operator, employee.temperature); err != nil {
		return 0, err
	}

	return manager.GetComfortTemp(), nil
}
