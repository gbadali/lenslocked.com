package main

import (
	"fmt"

	"github.com/gbadali/lenslocked.com/rand"
)

func main() {
	fmt.Println(rand.String(10))
	fmt.Println(rand.RememberToken())
}
