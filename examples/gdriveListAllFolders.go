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

	folders, err := gd.ListAllFolders()
	if err != nil {
		fmt.Println("error ListAllFolders:", err)
//		os.Exit(1)
	}

//	gdrive.PrintFileList("All Folders", folders)

	fmt.Printf("found %d folders!\n", len(folders))
	fmt.Println("Success!")
	os.Exit(0)
}
