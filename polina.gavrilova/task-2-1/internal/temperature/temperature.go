package temperature

import "errors"

var errInvalidOperation = errors.New("invalid operation")

const (
	MaxTemp = 30
	MinTemp = 15
)

type TempCondition struct {
	CurMin, CurMax int
}

func (cond *TempCondition) Change(mode string, parametr int) (bool, error) {
	switch mode {
	case ">=":
		cond.CurMin = max(cond.CurMin, parametr)
	case "<=":
		cond.CurMax = min(cond.CurMax, parametr)
	default:
		return false, errInvalidOperation
	}

	if cond.CurMax < cond.CurMin {
		return false, nil
	}

	return true, nil
}
