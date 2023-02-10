// gdriveGetFile
// author: prr
// 24/1/2022
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
	var gfiles []*drive.File


//	docId := "1lEodX98Eq6_2elpgct_OOv-5L5Es_iGyZJqrIS2BznY"
//    numArgs := len(os.Args)
/*
    if numArgs < 2 {
        fmt.Println("error - no comand line arguments!")
  		fmt.Println("gdrive_about usage is: \"gdriveListFile id\"\n")
//		fmt.Println("Using default docId: ", docId)
    } else {
		docId = os.Args[1]
		fmt.Printf("Using %s as docId\n", docId)
	}
*/
	err := gd.Init()
	if err != nil {
		fmt.Printf("error main::Init gdriveApi: %v\n", err)
		os.Exit(1)
	}

	gfiles, err = gd.ListFile()
	if err != nil {
		fmt.Printf("error::ListFile: %v\n", err)
		os.Exit(1)
	}

	gfilNum:= len(gfiles)
	fmt.Println("files: ",gfilNum)
/*
	filnam := gfil.Name
	ext := gfil.FileExtension
	fmt.Printf("Found File: %s.%s\n", filnam, ext)
	outfil, err := gd.CreTxtOutFile(filnam, ext)
	if err != nil {
		fmt.Println("error main::CreTxtOutFile -- cannot open out file: ", err)
		os.Exit(1)
	}


	res, err := gd.GetFile(docId)
	if err != nil {
		fmt.Println("error main::GetFile: ", err)
		os.Exit(1)
	}

	 _, err = io.Copy(outfil, res.Body)
    if err != nil {
        fmt.Println("error main -- cannot save downloaded file: ", err)
        os.Exit(1)
	}
*/
	fmt.Println("success ListFile!")
	os.Exit(0)
}
