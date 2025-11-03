package temperature

import "errors"

var (
	errInvalidOperation = errors.New("invalid operation")
	errInvalidRange     = errors.New("invalid temperature range")
)

const (
	MaxTemp = 30
	MinTemp = 15
)

type TempCondition struct {
	curMin, curMax int
}

func (cond *TempCondition) Change(mode string, parameter int) error {
	if cond.curMin == 0 && cond.curMax == 0 {
		cond.curMin = MinTemp
		cond.curMax = MaxTemp
	}

	switch mode {
	case ">=":
		cond.curMin = max(cond.curMin, parameter)
	case "<=":
		cond.curMax = min(cond.curMax, parameter)
	default:
		return errInvalidOperation
	}

	if cond.curMax < cond.curMin {
		return errInvalidRange
	}

	return nil
}

func (cond *TempCondition) GetCurMin() int {
	if cond.curMin == 0 && cond.curMax == 0 {
		return MinTemp
	}
	return cond.curMin
}

func (cond *TempCondition) GetCurMax() int {
	if cond.curMin == 0 && cond.curMax == 0 {
		return MaxTemp
	}
	return cond.curMax
}

func (cond *TempCondition) HasValidRange() bool {
	return cond.GetCurMin() <= cond.GetCurMax()
}

func (cond *TempCondition) GetTemperatureRange() (int, int) {
	return cond.GetCurMin(), cond.GetCurMax()
}
