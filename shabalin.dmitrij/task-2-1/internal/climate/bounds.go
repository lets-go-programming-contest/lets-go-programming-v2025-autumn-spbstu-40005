package climate

import "fmt"

const (
	InitialMin = 15
	InitialMax = 30
)

type Bounds struct {
	min int
	max int
}

type Controller struct {
	bounds *Bounds
}

func NewController() *Controller {
	return &Controller{
		bounds: &Bounds{
			min: InitialMin,
			max: InitialMax,
		},
	}
}

func (c *Controller) AddConstraint(op string, temp int) error {
	switch op {
	case ">=":
		if temp > c.bounds.min {
			c.bounds.min = temp
		}

	case "<=":
		if temp < c.bounds.max {
			c.bounds.max = temp
		}

	default:
		return fmt.Errorf("invalid operator: %s", op)
	}

	return nil
}

func (c *Controller) ComfortTemp() int {
	if c.bounds.min <= c.bounds.max {
		return c.bounds.min
	}
	return -1
}
