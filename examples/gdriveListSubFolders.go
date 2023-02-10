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

	gd, err := gdrive.InitDriveApi()
	if err != nil {
		fmt.Printf("Gdrive Api Init error: %v\n", err)
		os.Exit(1)
	}

	azulTestId :="1cRgugvok058kLc8Nxbg5vZaum0HbDAUkkIcrn29WY4P45PLbA-kuW4N_NteQeJZwBICBEppW"

	foldList, err := gd.ListSubFoldersId(azulTestId)
	if err != nil {
		fmt.Println("error ListSubFolders:", err)
//		os.Exit(1)
	}

    gdrive.PrintFileList("folders for AzulTest", foldList)

	folders, err := gd.ListAllSubFolders("AzulTest", azulTestId)
	if err != nil {
		fmt.Println("error ListAllFolders:", err)
//		os.Exit(1)
	}

/*
	numSubDir :=0
	if folders.SubDir != nil {
		numSubDir = len(*((*folders).SubDir))

	}
	fmt.Printf("dir %s: found %d subfolders!\n", folders.Name, numSubDir)
*/
	gdrive.PrintFolderList("AzulTest", folders)

	fmt.Println("Success!")
	os.Exit(0)
}
