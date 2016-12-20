package main

import (
	"fmt"
)

type DropKick struct{}

func (s *DropKick) Execute() {
	fmt.Println("DropKick EXECUTED")
}

func init() {
	Register("Dropkick", func() Ability {
		return &DropKick{}
	})
}
