// golang program to test gdrive
// author: prr, azul software
// created: 6/2/2023
// copyright 2023 prr, Peter Riemenschneider, Azul Software
//
//

package main

import (
        "fmt"
//        "os"
//		"io"
		gdrive "google/gdrive/examples/gdriveLib"
)


func main() {

	gdrive.ListExt()

	fmt.Println("Success!")
//	os.Exit(0)
}
