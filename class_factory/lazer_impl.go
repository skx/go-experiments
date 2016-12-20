package main

import (
	"fmt"
)

type Lazer struct{}

func (s *Lazer) Execute() {
	fmt.Println("LAZER")
}

func init() {
	Register("Lazer", func() Ability {
		return &Lazer{}
	})
}
