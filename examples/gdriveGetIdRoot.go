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

	idStr, err := gd.GetIdOfTopDir()
	if err != nil {
		fmt.Println("error GetIdOfTopDir:", err)
		os.Exit(1)
	}

	fmt.Printf("Root Folder Id: %s\n", idStr)

	fmt.Println("Success!")
	os.Exit(0)
}
