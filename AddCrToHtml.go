package main

import (
	"fmt"
	"os"
	"strings"
)

func main () {

	numArgs := len(os.Args)

	if numArgs < 2 {
		fmt.Printf("insufficient arguments\n  usage is: \"InsCr file\"\n")
		os.Exit(1)
	}
	infilNam := os.Args[1]

	infil, err := os.Open(infilNam)
	if err != nil {
		fmt.Printf("error opening input file! %v\n", err)
		os.Exit(1)
	}
	defer infil.Close()

	idx := strings.Index(infilNam, ".")
//	fmt.Printf("outfil:  %s %d\n", string(infilNam[:idx]), idx)

	outfilNam := infilNam[:idx] + "_cr.html"
	fmt.Printf("outfil name: %s\n", outfilNam)
	outfil, err := os.Create(outfilNam)
	if err != nil {
		fmt.Printf("error opening input file! %v\n", err)
		os.Exit(1)
	}
	defer outfil.Close()

	fileinfo, err := infil.Stat()
	if err != nil {
		fmt.Printf("error getting input file stat! %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("file size: %d\n", fileinfo.Size())
	inByte := make([]byte, fileinfo.Size())
	_, err = infil.Read(inByte)
	if err != nil {
		fmt.Printf("error reading input file! %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("input: %s\n", string(inByte[:15]))

	outByte := make([]byte, fileinfo.Size() + 1024)
	i2 := 0
	for i:=0; i< len(inByte); i++ {
		outByte[i2] = inByte[i]
		i2++
		if inByte[i] == '}' {
			outByte[i2] = '\n'
			i2++
		}
		if inByte[i] == '>' {
			if inByte[i-1] != ' ' {
				outByte[i2] = '\n'
				i2++
			}
		}
	}

	outfil.Write(outByte[:i2])

	fmt.Println("success!")
}
