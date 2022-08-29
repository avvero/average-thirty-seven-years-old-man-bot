package dice

import (
	"math/rand"
	"time"
)

type Dice struct {
	sides int
}

func Of(sides int) *Dice {
	return &Dice{sides: sides}
}

func (this Dice) Roll() int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(this.sides)
}
