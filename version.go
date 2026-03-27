package x2o

import (
	"fmt"
	"runtime"
)

// go build -ldflags "-X 'x2o.Version=value'"
var (
	Version   = ""
	GoVersion = runtime.Version()
	Platform  = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
)
