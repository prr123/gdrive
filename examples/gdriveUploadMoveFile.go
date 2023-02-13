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

/*
	numArgs := len(os.Args)
	if numArgs <2 {
		fmt.Printf("no enough arguments -- need folder name as second argument!\n")
		os.Exit(-1)
	}

	foldNam := os.Args[1]
*/
	gd, err := gdrive.InitDriveApi()
	if err != nil {
		fmt.Printf("Gdrive Api Init error: %v\n", err)
		os.Exit(1)
	}

// folder
	azulTestId :="1cRgugvok058kLc8Nxbg5vZaum0HbDAUkkIcrn29WY4P45PLbA-kuW4N_NteQeJZwBICBEppW"

	docId := "1GjD9109eAAfufreM6Oj1EZpyvi6BHtCX0ihQSsKrloU"
	fmt.Printf("Id: %s\n",docId)

	err = gd.ExportFileByIdDl(docId, "pdftest", "pdf")
	if err != nil {
		fmt.Printf("error Export: %v\n", err)
		os.Exit(-1)
	}

	fmt.Println("Success download!")

	filId, err := gd.UploadFile("pdftest.pdf", "pdf")
	if err != nil {
		fmt.Printf("error Upload: %v\n", err)
		os.Exit(-1)
	}

	fmt.Printf("uploaded file Id: %s\n", filId)

	err = gd.MoveFile(filId, azulTestId)
	if err != nil {
		fmt.Printf("error Move File: %v\n", err)
		os.Exit(-1)
	}

	fmt.Println("Success!")
//	os.Exit(0)
}