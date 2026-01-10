package climate

import (
	"errors"
	"fmt"
)

const (
	initialMin = 15
	initialMax = 30
)

var ErrInvalidOperator = errors.New("invalid operator")

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
			min: initialMin,
			max: initialMax,
		},
	}
}

func (c *Controller) AddConstraint(operator string, temperature int) error {
	switch operator {
	case ">=":
		if temperature > c.bounds.min {
			c.bounds.min = temperature
		}

	case "<=":
		if temperature < c.bounds.max {
			c.bounds.max = temperature
		}

	default:
		return fmt.Errorf("%w: %s", ErrInvalidOperator, operator)
	}

	return nil
}

func (c *Controller) ComfortTemp() int {
	if c.bounds.min <= c.bounds.max {
		return c.bounds.min
	}

	return -1
}
