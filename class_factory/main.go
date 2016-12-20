package main

import (
	"fmt"
	"sync"
)

type Ability interface {
	Execute()
}

var abilities = struct {
	m map[string]AbilityCtor
	sync.RWMutex
}{m: make(map[string]AbilityCtor)}

type AbilityCtor func() Ability

func Register(id string, newfunc AbilityCtor) {
	abilities.Lock()
	abilities.m[id] = newfunc
	abilities.Unlock()
}

func DumpAbilities() {
	for k, _ := range abilities.m {
		fmt.Println("\t" + k)
	}
}

func GetAbility(id string) (a Ability) {
	abilities.RLock()
	ctor, ok := abilities.m[id]
	abilities.RUnlock()
	if ok {
		a = ctor()
	}
	return
}

func main() {
	fmt.Println("Abilities:")
	DumpAbilities()

	if fireball := GetAbility("Fireball"); fireball != nil {
		fireball.Execute()
	}
}
