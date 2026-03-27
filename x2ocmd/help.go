package x2ocmd

import "fmt"

func showHelp() {
	fmt.Println("X2O")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  x2o -c config.json        Start the server (default)")
	fmt.Println("  x2o -v                    Show version")
	fmt.Println("  x2o -h                    Show this help message")
	fmt.Println("")
	// TODO: 补上
	fmt.Println("Online Documentation:")
}
