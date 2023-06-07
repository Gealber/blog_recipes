package utils

type AnonNumbers struct {
	Name    string
	Value   int
	Address string
	OnSell  bool // if I already put this number on auction
}

const (
	defaultPercentWin = 0.3 // 30 % percent
)

func (n *AnonNumbers) valueToSell() int {
	// fragment charge a 5%
	return int(float64(n.Value) * 1.35)
}
