package temperature

const (
	MaxTemp = 30
	MinTemp = 15
)

type TempCondition struct {
	CurMin, CurMax, CurTemp int
}

func (cond *TempCondition) Change(mode string, parametr int) {
	switch mode {
	case ">=":
		cond.CurMin = max(cond.CurMin, parametr)
		cond.CurTemp = max(cond.CurTemp, parametr)
	case "<=":
		cond.CurMax = min(cond.CurMax, parametr)
		cond.CurTemp = min(cond.CurTemp, parametr)
	}
}
