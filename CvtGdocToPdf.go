// v1
// v2 start a text debug
// 10/1/22
//
//

package main

import (
        "fmt"
        "os"
		"io"
		gdriveApi "google/gdrive/gdriveApi"
)


func main() {
	var gd gdriveApi.GdriveApiStruct

    numArgs := len(os.Args)
    if numArgs < 2 {
        fmt.Println("error - no comand line arguments!")
  		fmt.Println("CvtGdocToPdf usage is:\n  CvtGdocToPdf docId\n")
        os.Exit(1)
    }

    docId := os.Args[1]
	fmt.Println("Doc Id: ", docId)

	err := gd.Init()
	if err != nil {
		fmt.Println("Init error: ", err)
		os.Exit(1)
	}
	srv := gd.Svc

//	docId := "1pdI_GFPR--q88V3WNKogcPfqa5VFOpzDZASo4alCKrE"

	fil, err := srv.Files.Get(docId).Do()
	if err != nil {
		fmt.Println("Unable to find file!", err)
		os.Exit(1)
	}

	filnam := fil.Name
	fmt.Println("file name: ", filnam)
	if len(filnam) <1 {
		fmt.Println("not a valid file name")
		os.Exit(1)
	}

	mimeType := "application/pdf"
	res, err := srv.Files.Export(docId, mimeType).Context(gd.Ctx).Download()
	if err != nil {
		fmt.Println("Unable to download document into app: ", err)
		os.Exit(1)
	}

	outfil, err := gd.CreTxtOutFile(filnam, "pdf")
	if err != nil {
		fmt.Println("error main -- cannot open out file: ", err)
		os.Exit(1)
	}

	_, err = io.Copy(outfil, res.Body)
	if err != nil {
		fmt.Println("error main -- cannot convert gdoc file: ", err)
		os.Exit(1)
	}

	outfil.Close()
	fmt.Println("Success!")
	os.Exit(0)
}
