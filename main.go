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
	mainctl := false
	mainset := false
	localWarning := 0
	fmtctl := false
	fun := false
	lineExtend := false
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
			case lineExtend == true:
				goFile = goFile + strings.Repeat(" ", 4)
				for ii := 1; ii < len(ld); ii++ {
					if ld[ii] == ";" {
						lineExtend = true
					} else {
						if ld[ii] == ")" {
							goFile = goFile + ld[ii]
						} else {
							dtype := "string"
							goFile = goFile + ld[ii][0:len(ld[ii])-1] + " " + dtype
							if ii < len(ld)-2 {
								goFile = goFile + ", "
							}
						}
						lineExtend = false
					}
				}
				if lineExtend {
					goFile = goFile + "\n"
				} else {
					goFile = goFile + " {\n"
				}
			case ld[0] == "procedure" || ld[0] == "function":
				fun = true
				ftmp := strings.Split(ld[1], "(")
				funname = ftmp[0]
				if funname == "main" {
					mainctl = true
					mainset = true
					lineExtend = false
				}
				goFile = goFile + "func " + funname + "("
				if mainctl {
					goFile = goFile + ")"
				} else {
					for ii := 2; ii < len(ld); ii++ {
						if ld[ii] == ";" {
							lineExtend = true
							goFile = goFile + ", "
						} else {
							if ld[ii] == ")" {
								goFile = goFile + ld[ii]
							} else {
								if ii < len(ld)-1 {
									dtype := "string"
									goFile = goFile + ld[ii][0:len(ld[ii])-1] + " " + dtype
									if ii < len(ld)-2 {
										goFile = goFile + ", "
									}

								}
							}
						}
					}
				}
				if lineExtend {
					goFile = goFile + "\n"
				} else {
					goFile = goFile + " {\n"
				}
			case ld[0] == "return" || lda[0] == "return":
				fun = false
				mainctl = false
				goFile = goFile + "}\n\n"

			case ld[0] == "local":
				if len(ld) > 2 {
					goFile = goFile + strings.Repeat(" ", 4)
					for ii := 1; ii < len(ld); ii++ {
						switch {
						case ld[ii] == ".f.":
							goFile = goFile + "false"
						case ld[ii] == ".t.":
							goFile = goFile + "true"

						default:
							goFile = goFile + ld[ii]
						}
					}
					goFile = goFile + "\n"
				} else {
					localWarning++
					fmt.Printf("Warning  Local variable %S did NOT convert\n", ld[1])

				}
			case ld[0] == "?" || lpn == true:
				fmtctl = true
				goFile = goFile + strings.Repeat(" ", 4)
				goFile = goFile + "fmt.Println("

				for ii := 0; ii < len(ld); ii++ {
					if ld[ii][0:1] == "?" {
						goFile = goFile + ld[ii][1:len(ld[ii])] + " "
					} else {
						goFile = goFile + ld[ii] + " "
					}
				}
				goFile = goFile + ")\n"

			}

			line = ""
		}

	}
	if fun {
		goFile = goFile + "}\n\n"
	}

	top := "package "
	if mainset {
		top = top + "main\n\n"
	} else {
		top = top + funname + "\n\n"
	}
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
	if localWarning > 0 {
		fmt.Printf("Warning %d Local variables did NOT convert\n", localWarning)

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
