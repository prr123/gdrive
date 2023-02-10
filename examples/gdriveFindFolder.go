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
//		"google.golang.org/api/drive/v3"
		gdriveApi "google/gdrive/gdriveApi"
)


func main() {
	var gd gdriveApi.GdriveApiStruct
	var filesInfo *[]gdriveApi.FileInfo

    numArgs := len(os.Args)

    if numArgs < 2 {
        fmt.Println("error - no comand line arguments!")
  		fmt.Println("gdrive_about usage is: \"gdriveFindFolder id\"\n")
//		fmt.Println("Using default docId: ", docId)
	 	os.Exit(0)
    }

	err := gd.Init()
	if err != nil {
		fmt.Printf("error main::Init gdriveApi: %v\n", err)
		os.Exit(1)
	}

	filesInfo, nf, err := gd.ListFolderByName(os.Args[1])
	if err != nil {
		fmt.Printf("error::ListFolderByName: %v\n", err)
		os.Exit(1)
	}

	if nf < 1 {
		fmt.Printf("could not find a folder with name %s!", os.Args[1])
		os.Exit(0)
	}

	if nf > 0 {
		fmt.Printf("Found %d Folder(s)\n", nf)
		FilesInfo := *filesInfo
		for i:=0; i< nf; i++ {
			fmt.Printf("Folder: %3d name: %-35s mime: %-35s id: %s parent id: %s\n",i+1, FilesInfo[i].Name, FilesInfo[i].MimeType, FilesInfo[i].Id, FilesInfo[i].ParentId)
		}
	}
	fmt.Println("success ListFolderByName!")
	os.Exit(0)
}
