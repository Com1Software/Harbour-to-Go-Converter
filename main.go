package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func BuildApp(sfile string) {
	dfile := strings.TrimSuffix(sfile, filepath.Ext(sfile)) + ".go"
	fmt.Printf("Converting  : %s to %s \n", sfile, dfile)
	prgFile, err := os.Open(sfile)
	goFile := ""
	if err != nil {
		fmt.Println(err)
	}
	defer prgFile.Close()
	fmt.Println("Successfully Opened File")
	bValue, _ := ioutil.ReadAll(prgFile)
	byteValue := strings.ToLower(string(bValue))
	fmt.Printf("Size file %d\n", len(byteValue))
	line := ""
	fmtctl := false
	fun := false
	funname := ""
	for i := 0; i < len(byteValue); i++ {
		if string(byteValue[i:i+1]) != "\n" {

			line = line + string(byteValue[i:i+1])
		}
		if string(byteValue[i:i+1]) == "\n" {
			xline := strings.TrimLeft(line, " ")
			ld := strings.Split(xline, " ")
			lda := strings.Split(ld[0], "(")
			lpn := false
			if len(ld[0]) > 0 {
				if ld[0][0:1] == "?" {
					lpn = true
				}

			}

			fmt.Printf(" Line Segmnets %d %s \n ", len(ld), line)
			switch {
			case ld[0] == "procedure" || ld[0] == "function":
				fun = true
				goFile = goFile + "func " + ld[1] + " {\n"
				ftmp := strings.Split(ld[1], "(")
				funname = ftmp[0]

			case ld[0] == "return" || lda[0] == "return":
				fun = false
				goFile = goFile + "}\n"

			case ld[0] == "?" || lpn == true:
				fmtctl = true
				goFile = goFile + strings.Repeat(" ", 4)
				goFile = goFile + "fmt.Println("
				for ii := 1; ii < len(ld); ii++ {
					goFile = goFile + ld[ii]
				}
				goFile = goFile + ")\n"

			}

			line = ""
		}

	}
	if fun {
		goFile = goFile + "}\n"
	}

	top := "package " + funname + "\n\n"
	if fmtctl {
		top = top + "import (\n"
		f := "fmt"
		top = top + strings.Repeat(" ", 4)
		top = top + fmt.Sprintf("%q\n", f)
		top = top + ")\n\n"
	}
	goFile = top + goFile

	f, err := os.Create(dfile)
	if err != nil {
		fmt.Printf("Error %s\n", err)
	}
	l, err := f.WriteString(goFile)
	if err != nil {
		fmt.Printf("Error %s\n", err)
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Printf("Error %s\n", err)
	}

}

func main() {
	fmt.Printf("Harbour to Go Converter \n")
	fmt.Printf("Operating System : %s\n", runtime.GOOS)
	fmt.Printf("%d Parameters\n", len(os.Args))
	switch {
	//----------------------------
	case len(os.Args) == 1:
		fmt.Printf("Missing File to Convert \n")
		return
	case len(os.Args) == 2:
		sfile := os.Args[1]
		BuildApp(sfile)

	}

}
