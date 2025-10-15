package temperature

import "errors"

const (
	MaxTemp = 30
	MinTemp = 15
)

type TempCondition struct {
	CurMin, CurMax, CurTemp int
}

func (cond *TempCondition) Change(mode string, parametr int) error {
	switch mode {
	case ">=":
		cond.CurMin = max(cond.CurMin, parametr)
		cond.CurTemp = max(cond.CurTemp, parametr)
	case "<=":
		cond.CurMax = min(cond.CurMax, parametr)
		cond.CurTemp = min(cond.CurTemp, parametr)
	}
	if cond.CurMin > cond.CurMax {
		return errors.New("Invalid changing temperature")
	}
	return nil
}
