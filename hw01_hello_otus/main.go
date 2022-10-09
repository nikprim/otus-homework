package main

import (
	"fmt"
	"golang.org/x/example/stringutil"
)

func main() {
	word := "Hello, OTUS!"

	fmt.Println(stringutil.Reverse(word))
}
