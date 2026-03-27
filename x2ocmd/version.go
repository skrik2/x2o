package x2ocmd

import (
	"fmt"

	"github.com/skrik2/x2o"
)

func showVersion() {
	fmt.Println(x2o.Version)
	fmt.Println(x2o.GoVersion)
	fmt.Println(x2o.Platform)
}
