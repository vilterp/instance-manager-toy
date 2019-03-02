package manager

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func sleepRandom(base time.Duration, variance time.Duration) {
	varDur := time.Duration(math.Round(rand.Float64() * float64(variance)))
	fmt.Println("sleeping for", base, varDur)
	time.Sleep(base + varDur)
}
