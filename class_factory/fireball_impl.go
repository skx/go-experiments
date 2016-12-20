package main

import (
	"fmt"
)

type Fireball struct{}

func (s *Fireball) Execute() {
	fmt.Println("FIREBALL EXECUTED")
}

func init() {
	Register("Fireball", func() Ability {
		return &Fireball{}
	})
}
