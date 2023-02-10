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
	var path string

    numArgs := len(os.Args)

    if numArgs < 2 {
        fmt.Println("error - no comand line arguments!")
  		fmt.Println("gdrive_about usage is: \"gdriveGetFullPath id\"\n")
//		fmt.Println("Using default docId: ", docId)
	 	os.Exit(0)
    }

	err := gd.Init()
	if err != nil {
		fmt.Printf("error main::Init gdriveApi: %v\n", err)
		os.Exit(1)
	}

	filesInfo, path, err = gd.GetFullPath(os.Args[1])
	if err != nil {
		fmt.Printf("error::GetFileByName: %v\n", err)
		os.Exit(1)
	}

	if filesInfo == nil {
		fmt.Println("no files!")
	} else {
		fmt.Printf("Path: \"%s\" File Infos: %d\n",path, len(*filesInfo))
		FilesInfo := *filesInfo
		for i:=0; i< len(*filesInfo); i++ {
			fmt.Printf("file: %3d name: %-35s mime: %-35s id: %s\n",i+1, FilesInfo[i].Name, FilesInfo[i].MimeType, FilesInfo[i].Id)
		}
	}
	fmt.Println("success GetFullPath!")
	os.Exit(0)
}
