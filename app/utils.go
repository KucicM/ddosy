package ddosy

import (
	"math/rand"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

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

func NewWeightedTargeter(pattern TrafficPattern) vegeta.Targeter {
	dist := pattern.dist
	rand.Seed(time.Now().UnixNano())
	return func(t *vegeta.Target) error {
		i := RandomPick(dist.weigths)
		t.Body = dist.payloads[i]
		t.URL = pattern.endpoint
		t.Method = pattern.method
		t.Header = pattern.header
		return nil
	}
}
