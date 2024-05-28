package main

import (
	"fmt"

	"github.com/agungdhewe/dwsqlc"
)

func main() {
	fmt.Println("test")

	cmd := dwsqlc.New()
	cmd.Ready()
}
