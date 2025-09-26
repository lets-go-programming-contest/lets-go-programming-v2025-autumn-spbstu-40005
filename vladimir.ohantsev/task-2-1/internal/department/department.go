package department

import "errors"

var errInvalidOperation = errors.New("invalid operator")

const (
	UpperLimit = 30
	LowerLimit = 15
)

type ComparisonOperator string

const (
	LessEqual    ComparisonOperator = "<="
	GreaterEqual ComparisonOperator = ">="
)

type ChangeRequest struct {
	Operator    ComparisonOperator
	Temperature int
}

type Department struct {
	minT, maxT int
}

func New() *Department {
	return &Department{
		minT: LowerLimit,
		maxT: UpperLimit,
	}
}

func (d *Department) SetMax(temp int) {
	if d.maxT > temp {
		d.maxT = temp
	}
}

func (d *Department) SetMin(temp int) {
	if d.minT < temp {
		d.minT = temp
	}
}

func (d *Department) Recalculate(req *ChangeRequest) (int, error) {
	switch req.Operator {
	case LessEqual:
		d.SetMax(req.Temperature)
	case GreaterEqual:
		d.SetMin(req.Temperature)
	default:
		return -1, errInvalidOperation
	}

	return d.Optimum(), nil
}

func (d *Department) Optimum() int {
	if d.maxT < d.minT {
		return -1
	}

	return d.minT
}
