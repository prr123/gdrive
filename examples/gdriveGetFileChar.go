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
		"google.golang.org/api/drive/v3"
		gdriveApi "google/gdrive/gdriveApi"
)


func main() {
	var gd gdriveApi.GdriveApiStruct
	var gfil *drive.File


	docId := "1lEodX98Eq6_2elpgct_OOv-5L5Es_iGyZJqrIS2BznY"
    numArgs := len(os.Args)
    if numArgs < 2 {
        fmt.Println("error - no comand line arguments!")
  		fmt.Println("gdrive_about usage is: \"gdriveGetFile id\"\n")
		fmt.Println("Using default docId: ", docId)
    } else {
		docId = os.Args[1]
		fmt.Printf("Using %s as output filename\n", docId)
	}

	err := gd.Init()
	if err != nil {
		fmt.Printf("error main::Init gdriveApi: %v\n", err)
		os.Exit(1)
	}

	gfil, err = gd.GetFileChar(docId)
	if err != nil {
		fmt.Printf("error::GetFileChar: %v\n", err)
		os.Exit(1)
	}

	filnam := gfil.Name
	outfil, err := gd.CreTxtOutFile(filnam, "txt")
	if err != nil {
		fmt.Println("error main::CreTxtOutFile -- cannot open out file: ", err)
		os.Exit(1)
	}


	err = gd.DumpFileChar(gfil, outfil)
	if err != nil {
		fmt.Println("error main::DumpFileChar: ", err)
		os.Exit(1)
	}

	fmt.Println("success GetFile!")
	os.Exit(0)
}
