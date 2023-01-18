package ddosy

import (
	"math/rand"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

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
