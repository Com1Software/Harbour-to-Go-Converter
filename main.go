package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	fmt.Printf("HTMLtoFunc Convert HTML Page to Go Function \n")
	fmt.Printf("Operating System : %s\n", runtime.GOOS)
	fmt.Printf("%d Parameters\n", len(os.Args))
	switch {
	//----------------------------
	case len(os.Args) == 2:
		file := os.Args[1]
		fmt.Printf("File : %s\n", file)
	case len(os.Args) == 1:
		file := "data/test.html"
		fmt.Printf("File : %s\n", file)

	}

}
