// gdriveGetFile
// author: prr
// 24/1/2022
//
//

package main

import (
        "fmt"
        "os"
		"io"
		"google.golang.org/api/drive/v3"
		gdLib "google/gdrive/gdriveApi"
)


func main() {
	var gd gdLib.GdriveApiStruct
	var gfil *drive.File


	filId := "1lEodX98Eq6_2elpgct_OOv-5L5Es_iGyZJqrIS2BznY"
    numArgs := len(os.Args)
    if numArgs < 2 {
        fmt.Println("error - no comand line arguments!")
  		fmt.Println("gdrive_about usage is: \"gdriveGetFile id\"\n")
		fmt.Println("Using default docId: ", docId)
    } else {
		docId = os.Args[1]
		fmt.Printf("Using %s as docId\n", docId)
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
	ext := gfil.FileExtension
	fmt.Printf("Found File: %s.%s\n", filnam, ext)
	outfil, err := gd.CreTxtOutFile(filnam, ext)
	if err != nil {
		fmt.Println("error main::CreTxtOutFile -- cannot open out file: ", err)
		os.Exit(1)
	}


	res, err := gd.GetFile(filId)
	if err != nil {
		fmt.Println("error main::GetFile: ", err)
		os.Exit(1)
	}

	 _, err = io.Copy(outfil, res.Body)
    if err != nil {
        fmt.Println("error main -- cannot save downloaded file: ", err)
        os.Exit(1)
	}

	fmt.Println("success GetFile!")
	os.Exit(0)
}
