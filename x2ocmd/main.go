package x2ocmd

import (
	"os"
)

func Main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-c":
			handleConfig()
			return
		case "-v":
			showVersion()
			return
		case "-h":
			showHelp()
			return
		default:
			showHelp()
			return
		}
	}

}
