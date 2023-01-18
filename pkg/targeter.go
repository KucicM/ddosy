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
