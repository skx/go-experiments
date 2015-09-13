/**
 * See if we can hash integers, such that
 * we could create unique "short identifiers"
 * for uploaded files.
 */

package main

import (
	"fmt"
	"github.com/speps/go-hashids"
)

func main() {
	hd := hashids.NewData()
	hd.Salt = "I hope this is secure"
	hd.MinLength = 1
	h := hashids.NewWithData(hd)

	for i := 0; i < 10000; i++ {
		numbers := []int{99}
		numbers[0] = i

		e, _ := h.Encode(numbers)
		fmt.Printf("%d -> %s\n", i, e)
	}
}
