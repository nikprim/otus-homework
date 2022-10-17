package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 3 || len(args[1]) == 0 || len(args[2]) == 0 {
		log.Fatal("not valid arguments")
	}

	envs, err := ReadDir(args[1])
	if err != nil {
		log.Fatal(err)
	}

	code := RunCmd(args[2:], envs)
	os.Exit(code)
}
