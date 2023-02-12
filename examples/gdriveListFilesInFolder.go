// golang program to test gdrive
// author: prr, azul software
// created: 6/2/2023
// copyright 2023 prr, Peter Riemenschneider, Azul Software
//
//

package main

import (
        "fmt"
        "os"
//		"io"
		gdrive "google/gdrive/examples/gdriveLib"
)


func main() {

	numArgs := len(os.Args)
	if numArgs <2 {
		fmt.Printf("no enough arguments -- need folder name as second argument!\n")
		os.Exit(-1)
	}

	foldNam := os.Args[1]

	gd, err := gdrive.InitDriveApi()
	if err != nil {
		fmt.Printf("Gdrive Api Init error: %v\n", err)
		os.Exit(1)
	}


	azulTestId :="1cRgugvok058kLc8Nxbg5vZaum0HbDAUkkIcrn29WY4P45PLbA-kuW4N_NteQeJZwBICBEppW"
	fmt.Printf("Id: %s\n",azulTestId)

	foldList, err := gd.ListFolderByName(foldNam)
	if err != nil {
		fmt.Println("error ListSubFolders:", err)
//		os.Exit(1)
	}

	fmt.Printf("** Found %d folders with name %s:\n", len(foldList), foldNam)

    gdrive.PrintFileList("folders for AzulTest", foldList)

	dirId := ""
	if len(foldList) > 0 {
		dirId = foldList[0].Id
		fmt.Printf("dirId: %s\n",dirId)
	} else { os.Exit(-1)}

	filList, err := gd.ListFilesInFolder(dirId)
	if err != nil {
		fmt.Println("error ListFiles in Folder:", err)
		os.Exit(-1)
	}

	gdrive.PrintFileList(foldNam, filList)

	fmt.Println("Success!")
	os.Exit(0)
}
