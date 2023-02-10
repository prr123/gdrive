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
		gdLib "google/gdrive/gdriveApi"
)


func main() {
	var gd gdLib.GdriveApiStruct
	var gfil *drive.File

    numArgs := len(os.Args)
    if numArgs < 3 {
        fmt.Println("error - insufficient comand line arguments!")
  		fmt.Println("usage is: \"gdriveExportFile id file\"")
		fmt.Println("*** exiting! ***")
		os.Exit(2)
    }

	filId := os.Args[1]
	filNam := os.Args[2]
	fmt.Printf("File Id: %s Out File Name: %s\n", filId, filNam)

	extIdx := 0
	for i:=len(filNam)-1; i>0; i-- {
		if filNam[i] == '.' {
			extIdx = i
			break
		}
	}

	if extIdx == 0 {
		fmt.Printf("error filnam has no extension!\n exiting!\n")
		os.Exit(2)
	}

	fmt.Printf("file name: %s extension %s\n", filNam[:extIdx], filNam[extIdx+1:])

	err := gd.Init()
	if err != nil {
		fmt.Printf("error main::Init gdriveApi: %v\n", err)
		os.Exit(1)
	}

	gfil, err = gd.GetFileChar(filId)
	if err != nil {
		fmt.Printf("error::GetFileChar: %v\n", err)
		os.Exit(1)
	}

	filnam := gfil.Name
	ext := gfil.FileExtension
	fmt.Printf("Found File: %s ext: %s\n", filnam, ext)

	err = gd.CreDumpFile(filId, filnam)
	if err != nil {
		fmt.Println("error CreDumpFile -- cannot create dump file: ", err)
		os.Exit(1)
	}

	dlFilNam := "output/" + filNam + ".jpeg"
	err = gd.DownloadFileById(filId, dfilNam)
	if err != nil {
		fmt.Println("error ExportFile: ", err)
		os.Exit(1)
	}

/*
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
	defer res.Body.Close()
*/
	fmt.Println("success ExportFile!")
	os.Exit(0)
}
