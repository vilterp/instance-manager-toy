package util

import (
	"math"
	"math/rand"
	"time"
)

// Distribution of durations
type DurDist struct {
	Base     time.Duration
	Variance time.Duration
}

func (d DurDist) GetRandom() time.Duration {
	varDur := time.Duration(math.Round(rand.Float64() * float64(d.Variance)))
	return d.Base + varDur
}

func (d DurDist) SleepRandom() {
	time.Sleep(d.GetRandom())
}
