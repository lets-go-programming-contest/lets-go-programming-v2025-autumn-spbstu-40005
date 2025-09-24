package department

const (
	UpperLimit = 30
	LowerLimit = 15
)

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

func (d *Department) Optimum() int {
	if d.maxT < d.minT {
		return -1
	}
	return d.minT
}
