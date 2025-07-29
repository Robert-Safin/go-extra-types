package main

import (
	"fmt"

	"github.com/Robert-Safin/go-extra-types/iter"
)

func main() {

	it := iter.NewIter([]int{1, 4, 3, 4, 4, 4})

	fmt.Println(it)

}
