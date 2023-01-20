package ddosy

import "math/rand"

func CreateDistribution(ws []float64) []float64 {
	sum := 0.0
	for _, w := range ws {
		sum += w
	}

	dist := make([]float64, len(ws))

	cs := 0.0
	for i, w := range ws {
		cs += w / sum
		dist[i] = cs
	}

	return dist
}


func RandomPick(ws []float64) int {
	p := rand.Float64()
	for i, w := range ws {
		if p < w {
			return i
		}
	}
	return len(ws) - 1
}