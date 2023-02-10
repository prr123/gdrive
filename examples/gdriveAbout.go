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
    if numArgs < 2 {
        fmt.Println("error - no comand line arguments!")
  		fmt.Println("gdrive_about usage is: \"gdriveAbout\"\n")
        os.Exit(1)
    }

	gd, err := gdrive.InitDriveApi()
	if err != nil {
		fmt.Printf("Gdrive Api Init error: %v\n", err)
		os.Exit(1)
	}

	resp, err := gd.GetAbout()
	if err != nil {
		fmt.Println("error svc.about.get:", err)
		os.Exit(1)
	}


	fmt.Println("Success!")
	os.Exit(0)
}
