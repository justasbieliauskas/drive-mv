package main

import (
	"fmt"
	"os"

	"github.com/justasbieliauskas/drivemv/command"
)

func main() {
	mv := command.New()
	err := mv.Run(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
