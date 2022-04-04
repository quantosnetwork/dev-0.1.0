package fees

// CalculateFees fee formula: (value * 100) / (numbytes * deflation) / 1024 / 2
// CalculateFees per transaction implements fee formula
func CalculateFees(value int64, numbytes int64, deflation float64) float64 {
	f := (float64(value) * 100) / (float64(numbytes) * deflation) / 1024 / 2
	return f
}
