package temperature

import "errors"

var errInvalidOperation = errors.New("invalid operation")

type TempCondition struct {
	curMin, curMax int
}

func NewTempCondition(minTemp, maxTemp int) *TempCondition {
	return &TempCondition{
		curMin: minTemp,
		curMax: maxTemp,
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

func (cond *TempCondition) GetCurrent() (int, int) {
	return cond.curMin, cond.curMax
}

func (cond *TempCondition) HasValidRange() bool {
	return cond.curMin <= cond.curMax
}
