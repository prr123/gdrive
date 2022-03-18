// v1
// v2 start a text debug
// 10/1/22
//
//

package main

import (
        "fmt"
        "os"
		gdrive "google/gdrive/gdriveApi"
)


func main() {
	var gd gdrive.GdriveApiStruct

    numArgs := len(os.Args)
    if numArgs < 3 {
        fmt.Println("error - insufficient comand line arguments!")
  		fmt.Printf(" usage is:\n %s docId mime\n", os.Args[0])
        os.Exit(1)
    }

    docId := os.Args[1]
	fmt.Println("Doc Id: ", docId)

	mime := os.Args[2]
/*
	mimeType, ok := gdrive.Gapp[mime]
	if !ok {
		fmt.Printf("error cvtMime: mime %s is not valid! \n", mime)
		os.Exit(1)
	}
	fmt.Println("mimeType: ",mimeType)
*/

	err := gd.Init()
	if err != nil {
		fmt.Println("error Init:", err)
		os.Exit(1)
	}
	srv := gd.Svc

//	docId := "1pdI_GFPR--q88V3WNKogcPfqa5VFOpzDZASo4alCKrE"

	fil, err := srv.Files.Get(docId).Do()
	if err != nil {
		fmt.Println("error Get: Unable to find file!", err)
		os.Exit(1)
	}

	filnam := fil.Name
	fmt.Println("file name: ", filnam)
	if !(len(filnam) >0) {
		fmt.Println("erro no file name")
		os.Exit(1)
	}

//	mimeType := "application/pdf"
	err = gd.ExportFileById(docId, filnam, mime)
	if err != nil {
		fmt.Println("Unable to download document into app: ", err)
		os.Exit(1)
	}

	fmt.Println("Success!")
	os.Exit(0)
}
