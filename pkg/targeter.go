package ddosy

// import (
// 	"math/rand"
// 	"time"

// 	vegeta "github.com/tsenart/vegeta/v12/lib"
// )

// func NewWeightedTargeter(patterns []TrafficPattern) vegeta.Targeter {
// 	distribution := createDistribution(patterns)
// 	rand.Seed(time.Now().UnixNano())
// 	return func(t *vegeta.Target) error {
// 		i := pickRandom(distribution)
// 		pattern := patterns[i]
// 		t.Body = []byte(pattern.Payload) // TODO 
// 		t.URL = pattern // TODO
// 		t.Header = pattern
// 	}
// }

// func createDistribution(patterns []TrafficPattern) []float64 {
// 	sum := 0.0
// 	for _, pattern := range patterns {
// 		sum += pattern.Weight
// 	}

// 	distribution := make([]float64, len(patterns))
// 	for i, pattern := range patterns {
// 		distribution[i] = pattern.Weight / sum
// 	}

// 	return distribution
// }

// func pickRandom(distribution []float64) int {
// 	p := rand.Float64()
// 	for i, pd := range distribution {
// 		if p < pd {
// 			return i
// 		}
// 	}
// 	return len(distribution) - 1
// }