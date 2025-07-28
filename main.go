package main

import (
	"fmt"

	"github.com/Robert-Safin/go-lib/option"
)

type user struct {
	name string
}

func main() {
	var bob string
	opt1 := option.NewInfer(bob, true)
	fmt.Println(opt1)

}
