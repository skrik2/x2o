package x2ocmd

import (
	"fmt"
	"os"

	"github.com/skrik2/x2o"
)

func handleConfig() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: x2o -c config.json")
		os.Exit(1)
	}

	x2o.Load(os.Args[2])
}
