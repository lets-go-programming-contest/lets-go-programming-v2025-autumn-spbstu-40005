package temperature

import "errors"

var errInvalidOperation = errors.New("invalid operation")

type TempCondition struct {
	curMin, curMax int
}

func NewTempCondition(min int, max int) *TempCondition {
	return &TempCondition{
		curMin: min,
		curMax: max,
	}
}

func (cond *TempCondition) Change(mode string, parameter int) error {
	switch mode {
	case ">=":
		cond.curMin = max(cond.curMin, parameter)
	case "<=":
		cond.curMax = min(cond.curMax, parameter)
	default:
		return errInvalidOperation
	}

	return nil
}

func (cond *TempCondition) GetCurrent() (min int, max int) {
	return cond.curMin, cond.curMax
}

func (cond *TempCondition) HasValidRange() bool {
	return cond.curMin <= cond.curMax
}
