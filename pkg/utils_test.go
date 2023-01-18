package ddosy_test

import (
	"testing"

	ddosy "github.com/kucicm/ddosy/pkg"
)

func TestNormDistribution(t *testing.T) {
	dist := []float64{1, 2, 1}
	normDist := ddosy.CreateDistribution(dist)

	expected := []float64{0.25, 0.75, 1}
	for i, actual := range normDist {
		if actual != expected[i] {
			t.Errorf("expected %+v got %+v", expected, actual)
		}
	}
}