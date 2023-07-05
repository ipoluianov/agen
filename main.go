package main

import (
	_ "embed"
	"fmt"

	"github.com/ipoluianov/agen/agen"
)

func main() {
	fmt.Println("Processing ...")
	err := agen.ProcessDirectory(".")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Complete")
}
