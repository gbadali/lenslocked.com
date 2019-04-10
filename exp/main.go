package main

import (
	"fmt"

	"github.com/gbadali/lenslocked.com/hash"
)

func main() {
	hmac := hash.NewHMAC("my-secret-key")
	// this should print out:
	fmt.Println(hmac.Hash("this is my string to hash"))
}
