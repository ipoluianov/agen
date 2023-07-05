package main

import (
	_ "embed"

	"github.com/ipoluianov/agen/agen"
)

func main() {
	agen.ProcessDirectory(".")
}
