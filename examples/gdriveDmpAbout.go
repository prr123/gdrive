// v1
// v2 start a text debug
// 10/1/22
//
//

package main

import (
        "fmt"
        "os"
//		"io"
		gdriveApi "google/gdrive/gdriveApi"
)


func main() {
	var gd gdriveApi.GdriveApiStruct

	filnam := "driveDump"
    numArgs := len(os.Args)
    if numArgs < 2 {
        fmt.Println("error - no comand line arguments!")
  		fmt.Println("gdrive_about usage is: \"gdriveDmpAbout outfile\"\n")
		fmt.Println("Using default filename: ",filnam)
    } else {
		filnam = os.Args[1]
		fmt.Printf("Using %s as output filename\n", filnam)
	}

	err := gd.Init()
	if err != nil {
		fmt.Printf("error main:init gdriveApi: %v\n", err)
		os.Exit(1)
	}

	resp, err := gd.GetAbout()
	if err != nil {
		fmt.Println("error main:GetAbout: ", err)
		os.Exit(1)
	}

	outfil, err := gd.CreTxtOutFile(filnam, "txt")
	if err != nil {
		fmt.Println("error main:CreTxtOutFile -- cannot open out file: ", err)
		os.Exit(1)
	}


	err = gd.DumpAbout(resp, outfil)
	if err != nil {
		fmt.Println("error main:DumpAbout: ", err)
		os.Exit(1)
	}

	fmt.Println("success DumpAbout!")
	os.Exit(0)
}
