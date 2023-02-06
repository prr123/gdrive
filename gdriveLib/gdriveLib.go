// gdriveLib.go
// golang library for google drive file management
//
// author: prr, azul software
// date: 30/1/2022
// update: 7/6/2022
// update 8/1/2023 change initGdApi due to google change authorization
// copywrite 2022, 2023 prr, azul software
//

package gdriveLib

import (
        "context"
        "encoding/json"
        "fmt"
		"io"
		"strings"
        "net/http"
        "os"
        "golang.org/x/oauth2"
        "golang.org/x/oauth2/google"
        "google.golang.org/api/drive/v3"
        "google.golang.org/api/option"
		"google.golang.org/api/googleapi"
        "google.golang.org/api/docs/v1"
)

type GdApiObj  struct {
	Ctx context.Context
	GdSvc *drive.Service
	GdocSvc *docs.Service
	Doc *docs.Document
}

type FileInfo struct {
	Id string
	MimeType string
	Name string
	Ext string
	ParentName string
	ParentId string
	SingleParent bool
	ModTime string
	Size int64
}

type cred struct {
    Installed credItems `json:"installed"`
    Web credItems `json:"web"`
}

type credItems struct {
    ClientId string `json:"client_id"`
    ProjectId string `json:"project_id"`
    AuthUri string `json:"auth_uri"`
    TokenUri string `json:"token_uri"`
//  Auth_provider_x509_cert_url string `json:"auth_provider_x509_cert_url"`
    ClientSecret string `json:"client_secret"`
    RedirectUris []string `json:"redirect_uris"`
}

var Gapp = map[string]string {
	"gdoc": "application/vnd.google-apps.document",
	"gsheet": "application/vnd.google-apps.spreadsheet",
	"gdraw": "application/vnd.google-apps.drawing",
	"gscript": "application/vnd.google-apps.script",
	"photo": "application/vnd.google-apps.photo",
	"gslide": "application/vnd.google-apps.presentation",
	"gmap": "application/vnd.google-apps.map",
	"gform": "application/vnd.google-apps.form",
	"folder": "application/vnd.google-apps.folder",
	"file": "application/vnd.google-apps.file",
	"jpg": "image/jpeg",
	"png": "image/png",
	"svg": "image/svg+xml",
	"pdf": "application/pdf",
	"html": "text/html",
	"text": "text/plain",
	"rich": "application/rtf",
	"word": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"excel": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"csv": "text/csv",
	"ppt": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
}

func ListApps() {
// function that lists all available apps
	fmt.Printf("******** Apps *************\n")
 	for k, v := range Gapp {
		fmt.Printf("%-10s %-30s\n", k, v)
	}
}

// Retrieves a token, saves the token, then returns the generated client.
/*
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "/home/peter/go/src/google/gdrive/tokGdrive.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}
*/

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func (gdrive *GdApiObj) CreDumpFile(fid string, filnam string)(err error) {
// function that creates a text file to dump document file

	nfilnam := make([]byte,len(filnam), len(filnam)+5)
	for i:= len(filnam) -1; i > -1; i-- {
		nfilnam[i] = filnam[i]
		if filnam[i] == '.' {
			nfilnam[i] = '_'
		}
	}

	ext := "txt"

	for i:=0; i<len(nfilnam); i++ {
		if nfilnam[i] == ' ' {
			nfilnam[i] = '_'
		}
	}

	// check whether output directory exists
	filinfo, err := os.Stat("output")
	if os.IsNotExist(err) {
		return fmt.Errorf("sub-dir \"output\" does not exist!")
	}
	if err != nil {
		return fmt.Errorf("os.Stat: %v!", err)
	}
	if !filinfo.IsDir() {
		return fmt.Errorf("file \"output\" is not a directory!")
	}

	path:= "output/" + string(nfilnam) + "." + ext
//	fmt.Println("path: ",path)
	outfil, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("OpenFile: %v \n", err)
	}

	// get file attributes
	svc := gdrive.GdSvc
	gfil, err := svc.Files.Get(fid).Do()
	if err != nil {
		return fmt.Errorf("Files.Get: cannot get file with id: %s! %v", fid, err)
	}

	outstr := fmt.Sprintf("File Name: %s Extension: %s Full Ext: %s\n", gfil.Name, gfil.FileExtension, gfil.FullFileExtension)
	outstr += fmt.Sprintf("Mime Type: %s Size: %d\n", gfil.MimeType, gfil.Size)
	outstr += fmt.Sprintf("File Id: %s Version: %d\n", gfil.Id, gfil.Version)
	outstr += fmt.Sprintf("Created: %s\n", gfil.CreatedTime)
	outstr += fmt.Sprintf("Modified: %s\n", gfil.ModifiedTime)
	outstr += fmt.Sprintf("Description: %s\n", gfil.Description)
	outstr += fmt.Sprintf("Original Name: %s \n", gfil.OriginalFilename)
	outstr += fmt.Sprintf("Parents: %d\n", len(gfil.Parents))
	outstr += fmt.Sprintf("Thumbnail: %s\n", gfil.ThumbnailLink)
	outstr += fmt.Sprintf("Web Content Link: %s\n", gfil.WebContentLink)
	outstr += fmt.Sprintf("Web View Link: %s\n", gfil.WebViewLink)

	outfil.WriteString(outstr)

	return nil
}

func InitDriveApi() (gdObj *GdApiObj, err error) {
// method that initialises the GdriveApi structure and returns a  service pointer

	var cred cred
	var config oauth2.Config
	var gdobj GdApiObj

	ctx := context.Background()
	gdobj.Ctx = ctx

	credFilNam := "/home/peter/go/src/google/gdoc/loginCred.json"
	credbuf, err := os.ReadFile(credFilNam)
	if err != nil {return nil, fmt.Errorf("os.Read %s: %v!", credFilNam, err)}

	err = json.Unmarshal(credbuf,&cred)
    if err != nil {return nil, fmt.Errorf("json.UnMarshal credbuf: %v\n", err)}

	if len(cred.Installed.ClientId) > 0 {
		config.ClientID = cred.Installed.ClientId
		config.ClientSecret = cred.Installed.ClientSecret
	}
	if len(cred.Web.ClientId) > 0 {
		config.ClientID = cred.Web.ClientId
		config.ClientSecret = cred.Web.ClientSecret
	}

	config.Scopes = make([]string,2)
    config.Scopes[0] = "https://www.googleapis.com/auth/drive"
    config.Scopes[1] = "https://www.googleapis.com/auth/documents"

	config.Endpoint = google.Endpoint

   	tokFile := "tokNew.json"
   	tok, err := tokenFromFile(tokFile)
   	if err != nil {return nil, fmt.Errorf("tokenFromFile: %v!", err)}

	client := config.Client(context.Background(), tok)

	gdsvc, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {return nil, fmt.Errorf("Unable to create Drive Service: %v!", err)}

	gdobj.GdSvc = gdsvc

	gdocsvc, err := docs.NewService(ctx, option.WithHTTPClient(client))
    if err != nil {return nil, fmt.Errorf("Unable to create Doc Service: %v!", err)}

	gdobj.GdocSvc = gdocsvc

	return &gdobj, nil
}

func (gdrive *GdApiObj) GetAbout() (resp *drive.About, err error) {
// method that lists all files in a drive

	fields := []googleapi.Field{"nextPageToken, files(id, name, mimeType, parents, modifiedTime)"}

	svc := gdrive.GdSvc
    resp, err = svc.About.Get().Fields(fields...).Do()
    if err != nil {
        fmt.Println("error svc.about.get:", err)
        return nil, fmt.Errorf("GetAbout: %v", err)
    }
	return resp, nil
}

func (gdrive *GdApiObj) DumpAbout(about *drive.About, outfil *os.File) (err error) {
// method that writes results from an about query into the output file

	var outstr string

	if outfil == nil {
		return fmt.Errorf("error DumpAbout: no outfil initialised!")
	}
	outstr += fmt.Sprintf("User:                %s\n", about.User.DisplayName)
	outstr += fmt.Sprintf("  email:             %s\n", about.User.EmailAddress)
	outstr += fmt.Sprintf("  kind:              %s\n", about.User.Kind)
	outstr += fmt.Sprintf("  me:                %t\n", about.User.Me)
	outstr += fmt.Sprintf("  permission id:     %s\n", about.User.PermissionId)
	outstr += fmt.Sprintf("  photo:             %s\n\n", about.User.PhotoLink)

	outstr += fmt.Sprintf("AppInstalled:        %t\n", about.AppInstalled)
	outstr += fmt.Sprintf("CanCreateDrives:     %t\n", about.CanCreateDrives)
	outstr += fmt.Sprintf("CanCreateTeamDrives: %t\n", about.CanCreateTeamDrives)
	outstr += fmt.Sprintf("Kind:                %s\n", about.Kind)
	outstr += fmt.Sprintf("DriveThemes:         %d\n", len(about.DriveThemes))
	outstr += "Maximum Import Sizes: \n"
	for k, v := range about.MaxImportSizes {
		outstr += fmt.Sprintf("Import: %s Size %s\n",k, v)
	}
	outstr += "Export Formats:\n"
	outstr += fmt.Sprintf("Maximum Upload Size: %d\n", about.MaxUploadSize)

	for k, v := range about.ExportFormats {
		outstr += fmt.Sprintf("format: %s %d\n", k, len(v))
		for i:=0; i<len(v); i++ {
			outstr += fmt.Sprintf("  %d: %s\n",i, v[i])
		}
	}
	outstr += "Import Formats:\n"
	for k, v := range about.ImportFormats {
		outstr += fmt.Sprintf("format: %s %d\n", k, len(v))
		for i:=0; i<len(v); i++ {
			outstr += fmt.Sprintf("  %d: %s\n",i, v[i])
		}
	}

	outstr += "Storage Quota:\n"
	outstr += fmt.Sprintf("  Limit:                %d\n", about.StorageQuota.Limit)
	outstr += fmt.Sprintf("  Usage:                %d\n", about.StorageQuota.Usage)
	outstr += fmt.Sprintf("  Usage in Drive:       %d\n", about.StorageQuota.UsageInDrive)
	outstr += fmt.Sprintf("  Usage in Drive Trash: %d\n", about.StorageQuota.UsageInDriveTrash)

	_, err = outfil.WriteString(outstr)
	if err != nil {
		return fmt.Errorf("error DumpAbout: cannot write to outfil! err: %v", err)
	}
	outfil.Close()
	return nil
}


func (gdrive *GdApiObj) ListFiles() (fileList []*drive.File, err error) {
// method that lists all folders

	fields := []googleapi.Field{"nextPageToken, files(id, name, mimeType, parents, modifiedTime)"}

	pagetoken := ""
	fin := false

	qstr := fmt.Sprintf("mimeType != '%s' and name = 'root'", Gapp["folder"])

//	fmt.Println("qstr: ", qstr)
	for i:=0; i< 3; i++ {
		nfileList, err := gdrive.GdSvc.Files.List().
		Fields(fields...).
		Q(qstr).
		PageToken(pagetoken).Context(gdrive.Ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("ListFile [%d]: %v", i, err)
		}
		fmt.Println("files: ", len(nfileList.Files))
		if len(nfileList.Files) < 1 {
			return nil, fmt.Errorf("ListFile - no files found!")
		}
		for j:=0; j<10; j++ {
			fil := nfileList.Files[j]
			fmt. Println(" name: ", fil.Name, " Mime: ", fil.MimeType, " Parents: ", len(fil.Parents), " :", fil.Parents[0] )
		}

		fileList = append(fileList, nfileList.Files...)

		if len(nfileList.NextPageToken) < 1 {
			fin = true
			break;
		}
	fmt.Printf("call %d files: %d\n", i, len(fileList))
		pagetoken = nfileList.NextPageToken
	}

	if !fin {
		return fileList, fmt.Errorf("ListFile -- too many files > 1000!")
	}
	return fileList, nil
}

func (gdrive *GdApiObj) ListAllFiles(dirId string) (fileList []*drive.File, err error) {
// method that lists all files in a folder with id dirId

	fields := []googleapi.Field{"nextPageToken, files(id, name, mimeType, parents, modifiedTime)"}

	pagetoken := ""
	fin := false


	qstr := fmt.Sprintf("mimeType != '%s' and '%s' in parent", Gapp["folder"], dirId)

//	fmt.Println("qstr: ", qstr)
	for i:=0; i< 3; i++ {
		nfileList, err := gdrive.GdSvc.Files.List().
		Fields(fields...).
		Q(qstr).
		PageToken(pagetoken).Context(gdrive.Ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("ListFile[%d]: %v", i, err)
		}
		fmt.Println("files: ", len(nfileList.Files))
		if len(nfileList.Files) < 1 {
			return nil, fmt.Errorf("ListFile - no files found!")
		}
		for j:=0; j<10; j++ {
			fil := nfileList.Files[j]
			fmt. Println(" name: ", fil.Name, " Mime: ", fil.MimeType, " Parents: ", len(fil.Parents), " :", fil.Parents[0] )
		}

		fileList = append(fileList, nfileList.Files...)

		if len(nfileList.NextPageToken) < 1 {
			fin = true
			break;
		}
	fmt.Printf("call %d files: %d\n", i, len(fileList))
		pagetoken = nfileList.NextPageToken
	}

	if !fin {
		return fileList, fmt.Errorf("ListFile -- too many files > 1000!")
	}


	return fileList, nil
}


func (gdrive *GdApiObj) ListFilesByName(nam string, dirId string) (filList []*drive.File, err error) {
// method that lists all files with name 'nam' and folder id 'dirId'

	var qstr string

	fields := []googleapi.Field{"nextPageToken, files(id, name, mimeType, parents, modifiedTime)"}

	pagetoken := ""
	fin := false
	if len(dirId) > 0 {
		_, err = gdrive.GdSvc.Files.Get(dirId).Context(gdrive.Ctx).Do()
		if err != nil {return filList, fmt.Errorf("error gdrive::CopyFile: could not find folder with id: %s -- %v", dirId, err)}
		qstr = fmt.Sprintf("(mimeType != '%s' and name = '%s') and '%s' in parents", Gapp["folder"], nam, dirId)
	} else {
		qstr = fmt.Sprintf("mimeType != '%s' and name = '%s'", Gapp["folder"], nam)
	}

//	fmt.Println("qstr: ", qstr)
	for i:=0; i< 3; i++ {
		nfileList, err := gdrive.GdSvc.Files.List().
		Fields(fields...).
		Q(qstr).
		PageToken(pagetoken).Context(gdrive.Ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("ListFile[%d]: %v", i, err)
		}
		fmt.Println("files: ", len(nfileList.Files))
		if len(nfileList.Files) < 1 {
			return nil, fmt.Errorf("ListFile - no files found!")
		}
		for j:=0; j<10; j++ {
			fil := nfileList.Files[j]
			fmt. Println(" name: ", fil.Name, " Mime: ", fil.MimeType, " Parents: ", len(fil.Parents), " :", fil.Parents[0] )
		}

		filList = append(filList, nfileList.Files...)

		if len(nfileList.NextPageToken) < 1 {
			fin = true
			break;
		}
	fmt.Printf("call %d files: %d\n", i, len(filList))
		pagetoken = nfileList.NextPageToken
	}

	if !fin {
		return filList, fmt.Errorf("ListFile -- too many files > 1000!")
	}
	return filList, nil
}

func (gdrive *GdApiObj) ListFoldersByName(nam string) (filList []*drive.File, err error) {
// method that looks for all folders with name 'nam' and returns a splice of file pointers

	fields := []googleapi.Field{"nextPageToken, files(id, name, mimeType, parents, modifiedTime)"}

	pagetoken := ""
	fin := false
	qstr := fmt.Sprintf("mimeType = '%s' and name = '%s'", Gapp["folder"], nam)

//	fmt.Println("qstr: ", qstr)
	for i:=0; i< 3; i++ {
		nfileList, err := gdrive.GdSvc.Files.List().
		Fields(fields...).
		Q(qstr).
		PageToken(pagetoken).Context(gdrive.Ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("ListFile[%d]: %v", i, err)
		}
		fmt.Println("files: ", len(nfileList.Files))
		if len(nfileList.Files) < 1 {
			return nil, fmt.Errorf("ListFile - no files found!")
		}
		for j:=0; j<10; j++ {
			fil := nfileList.Files[j]
			fmt. Println(" name: ", fil.Name, " Mime: ", fil.MimeType, " Parents: ", len(fil.Parents), " :", fil.Parents[0] )
		}

		filList = append(filList, nfileList.Files...)

		if len(nfileList.NextPageToken) < 1 {
			fin = true
			break;
		}
//	fmt.Printf("call %d files: %d\n", i, len(filList))
		pagetoken = nfileList.NextPageToken
	}

	if !fin {
		return filList, fmt.Errorf("ListFile -- too many files > 1000!")
	}
	return filList, nil
}

func (gdrive *GdApiObj) ListFolderByName(nam string) (folderList *[]FileInfo, err error) {
// method that looks for all sub folders of folder with name 'nam' and returns a list of folder pointers

	if len(nam) < 1 { return nil, fmt.Errorf("error gdrive::ListFolderByName -- no name provided!")}

	fields := []googleapi.Field{"nextPageToken, files(id, name, mimeType, parents, modifiedTime)"}

	pagetoken := ""
	qstr := fmt.Sprintf("mimeType = '%s' and name = '%s'", Gapp["folder"], nam)

//	fmt.Println("qstr: ", qstr)
	nfileList, err := gdrive.GdSvc.Files.List().
		Fields(fields...).
		Q(qstr).
		PageToken(pagetoken).Context(gdrive.Ctx).Do()
	if err != nil {return nil, fmt.Errorf("List: %v", err)}
	numFolders := len(nfileList.Files)
//	fmt.Println("folders: ", numFolders)

	if len(pagetoken) > 0 {return nil, fmt.Errorf("too many folders (>100)!")}

	finfolist := make([]FileInfo, numFolders)
	for i:= 0; i< numFolders; i++ {
		fileinfo, _ := gdrive.CvtToFilInfo(nfileList.Files[i])
		finfolist[i] = *fileinfo
	}
	folderList = &finfolist
	return folderList, nil
}

func (gdrive *GdApiObj) ListFFByName(nam string) (filList []*drive.File, err error) {
// method that looks for all files with name 'nam'

	fields := []googleapi.Field{"nextPageToken, files(id, name, mimeType, parents, modifiedTime)"}

	pagetoken := ""
	fin := false
	qstr := fmt.Sprintf("name = '%s'", nam)

//	fmt.Println("qstr: ", qstr)
	for i:=0; i< 3; i++ {
		nfileList, err := gdrive.GdSvc.Files.List().
		Fields(fields...).
		Q(qstr).
		PageToken(pagetoken).Context(gdrive.Ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("ListFile[%d]: %v", i, err)
		}
		fmt.Println("files: ", len(nfileList.Files))
		if len(nfileList.Files) < 1 {
			return nil, fmt.Errorf("ListFile - no files found!")
		}
		for j:=0; j<10; j++ {
			fil := nfileList.Files[j]
			fmt. Println(" name: ", fil.Name, " Mime: ", fil.MimeType, " Parents: ", len(fil.Parents), " :", fil.Parents[0] )
		}

		filList = append(filList, nfileList.Files...)

		if len(nfileList.NextPageToken) < 1 {
			fin = true
			break;
		}
//	fmt.Printf("call %d files: %d\n", i, len(filList))
		pagetoken = nfileList.NextPageToken
	}

	if !fin {
		return filList, fmt.Errorf("ListFile -- too many files > 1000!")
	}
	return filList, nil
}

func (gdrive *GdApiObj) ListFilesBySize(foldId string, minSize int64) (filList []*drive.File, err error) {
// method that lists all files above a certain size

	var qstr string
	if len(foldId) < 1 {return nil, fmt.Errorf("error gdrive::ListFilesBySize: no foldId provided!")}

	fields := []googleapi.Field{"nextPageToken, files(id, name, mimeType, parents, modifiedTime)"}

	pagetoken := ""
	fin := false

	if len(foldId) > 0 {
		qstr = fmt.Sprintf("size > '%d' and '%s' in parents", minSize, foldId)
	} else {
		qstr = fmt.Sprintf("size > '%d'", minSize)
	}
//	fmt.Println("qstr: ", qstr)
	for i:=0; i< 3; i++ {
		nfileList, err := gdrive.GdSvc.Files.List().
		Fields(fields...).
		Q(qstr).
		PageToken(pagetoken).Context(gdrive.Ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("ListFile[%d]: %v", i, err)
		}
//		fmt.Println("files: ", len(nfileList.Files))
		if len(nfileList.Files) < 1 {
			return nil, fmt.Errorf("ListFile - no files found!")
		}
		for j:=0; j<10; j++ {
			fil := nfileList.Files[j]
			fmt. Println(" name: ", fil.Name, " Mime: ", fil.MimeType, " Parents: ", len(fil.Parents), " :", fil.Parents[0] )
		}

		filList = append(filList, nfileList.Files...)

		if len(nfileList.NextPageToken) < 1 {
			fin = true
			break;
		}
//	fmt.Printf("call %d files: %d\n", i, len(filList))
		pagetoken = nfileList.NextPageToken
	}

	if !fin {
		return filList, fmt.Errorf("ListFile -- too many files > 1000!")
	}

	return filList, nil
}

func (gdrive *GdApiObj) ListTopDir() (id string, err error) {
// method that provides the id of the root folder

	var topDir *drive.File
	topDir, err = gdrive.GdSvc.Files.Get("root").Context(gdrive.Ctx).Do()
	if err != nil {
		return "", fmt.Errorf("ListTopDir: %v", err)
	}
	return topDir.Id, nil
}



func (gdrive *GdApiObj) CopyFile(filId string, nam string, dirId string) (nfilId string, err error) {
// method that copies a file with id 'filId' to new file. The method returns the id of the new file 'nfilId
// still todo
// assign same parent id to new file

	var srcFil, destFil *drive.File
	var parId string
	if len(filId) < 1 {return "", fmt.Errorf("no filId string!")}
	if len(nam) < 1 {return "", fmt.Errorf("no nam string!")}

	srcFil, err = gdrive.GdSvc.Files.Get(filId).Context(gdrive.Ctx).Do()
	if err != nil {return "",fmt.Errorf("no source file with id: %s -- %v", filId, err)}

//	fmt.Printf("fil: \n%v\n", srcFil)

	if len(srcFil.Parents) > 0 {parId = srcFil.Parents[0]}

//	PrintDriveFile("source", srcFil)

	if len(dirId) > 0 {
		_, err = gdrive.GdSvc.Files.Get(dirId).Context(gdrive.Ctx).Do()
		if err != nil {return "",fmt.Errorf("no destination folder with id [%s]: %v", dirId, err)}
		parId = dirId
	}

	if len(parId) > 0 {srcFil.Parents[0] = parId}

	srcFil.Name = nam
	// very important
	srcFil.Id = ""
//fmt.Printf("new name: %s\n", srcFil.Name)
	destFil, err = gdrive.GdSvc.Files.Copy(filId, srcFil).Do()
	if err != nil {return "", fmt.Errorf("Files.Copy: %v", err)}

//	PrintDriveFile("dest", destFil)

	return destFil.Id, nil
}

func (gdrive *GdApiObj) CreateFile(pDirId string, nam string) (fileId string, err error) {
// A method that creates a file with the parent id 'pDirId' and name 'nam'. 
// The method return the file id of the created file

	var fil drive.File
	var dir *drive.File
	var par [1]string
	if len(pDirId) < 1 {return "", fmt.Errorf("no pDirId string!")}
	if len(nam) < 1 {return "", fmt.Errorf("no nam string!")}

	// we could check nam for invalid chars

	par[0] = pDirId
	fil.Parents = par[:]
	fil.Name = nam
	fil.MimeType = Gapp["folder"]

	dir, err = gdrive.GdSvc.Files.Create(&fil).Context(gdrive.Ctx).Do()
	if err != nil {
		return "", fmt.Errorf("Files.Create: %v", err)
	}
	return dir.Id, nil
}

func (gdrive *GdApiObj) CreateFolder(pDirId string, nam string) (folderId string, err error) {
// A method that creates a folder under the parent folder with id 'pDirId' and name 'nam'.
// The method returns the  file id 'folderId' of the newly created Folder.

	var fil drive.File
	var dir *drive.File
	var par [1]string
	if len(pDirId) < 1 {return "", fmt.Errorf("no pDirId string!")}
	if len(nam) < 1 {return "", fmt.Errorf("no nam string!")}

	// we could check nam for invalid chars

	par[0] = pDirId
	fil.Parents = par[:]
	fil.Name = nam
	fil.MimeType = Gapp["folder"]

	dir, err = gdrive.GdSvc.Files.Create(&fil).Context(gdrive.Ctx).Do()
	if err != nil {
		return "", fmt.Errorf("Files.Create: %v", err)
	}
	return dir.Id, nil
}

func (gdrive *GdApiObj) DeleteFileById(filId string) (err error) {
// A method that deletes a file with the id 'filId'

	if len(filId)<1 {return fmt.Errorf("no file id provided!")}

	err = gdrive.GdSvc.Files.Delete(filId).Context(gdrive.Ctx).Do()
	if err != nil {
		return fmt.Errorf("Files.Delete: %v", err)
	}
	return nil
}

func (gdrive *GdApiObj) DeleteFileByName(nam string) (err error) {
// A method that deletes a file indendified by name 'nam'
// not finished

	if len(nam)<1 {return fmt.Errorf("no nam provided!")}

	fid:= "abc"

	err = gdrive.GdSvc.Files.Delete(fid).Context(gdrive.Ctx).Do()
	if err != nil {
		return fmt.Errorf("Files.Delete: %v", err)
	}
	return nil
}



func (gdrive *GdApiObj) FetchFileById(fid string) (resp *http.Response, err error) {
// A method that downloads a file identified by id 'fil' which returns the file in the http response body

	if len(fid)<1 {return nil, fmt.Errorf("no nam provided!")}

	resp, err = gdrive.GdSvc.Files.Get(fid).Context(gdrive.Ctx).Download()
	if err != nil {
		return nil, fmt.Errorf("GetFile Download: %v", err)
	}
	return resp, nil
}

func (gdrive *GdApiObj) MoveFileById(filId string, dirId string) (err error) {
// A method that moves the file with file id 'filId' into a dierectory with the id 'dirId'

	var fil, updfil *drive.File
	var parentStr string

	if len(filId)<1 {return fmt.Errorf("no filId provided!")}
	if len(dirId)<1 {return fmt.Errorf("no dirId provided!")}
	fil, err = gdrive.GdSvc.Files.Get(filId).Context(gdrive.Ctx).Do()
	if err != nil {return fmt.Errorf("no file with id [%s]: %v", filId, err)}
	_, err = gdrive.GdSvc.Files.Get(dirId).Context(gdrive.Ctx).Do()
	if err != nil {return fmt.Errorf("no folder with id [%s]: %v", dirId, err)}

	// remove old parents first, if they exist
	updfil = fil
	updFilId := filId
	if len(fil.Parents) > 0 {
		parentStr = fil.Parents[0]
		for i:=1; i < len(fil.Parents); i++ {
			parentStr += "," + fil.Parents[i]
		}

		updfil, err = gdrive.GdSvc.Files.Update(filId, fil).RemoveParents(parentStr).Context(gdrive.Ctx).Do()
		if err != nil {	return fmt.Errorf("RemoveParents: %v", err)}
		updFilId = updfil.Id
	}

	updfil, err = gdrive.GdSvc.Files.Update(updFilId, updfil).AddParents(dirId).Context(gdrive.Ctx).Do()
	if err != nil {
		return fmt.Errorf("Update AddParents: %v", err)
	}

	return nil
}

func (gdrive *GdApiObj) GetFileById(filId string) (fil *drive.File, err error) {
// A method that returns a file with the id 'filId'

	if len(filId)<1 {return nil, fmt.Errorf("no filId provided!")}

	fil, err = gdrive.GdSvc.Files.Get(filId).Context(gdrive.Ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("GetFileById: %v", err)
	}
	return fil, nil
}

func (gdrive *GdApiObj) GetFileInfoById(filId string) (filinfo *FileInfo , err error) {
// A method that returns a pointer to the file info struct 'FileInfo' of the file with id 'filId'

	fields := []googleapi.Field{"id, name, mimeType, parents, fullFileExtension, modifiedTime, size"}

	if len(filId)<1 {return nil, fmt.Errorf("no filId provided!")}

	fil, err := gdrive.GdSvc.Files.Get(filId).Fields(fields...).Context(gdrive.Ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("GetFiles: %v", err)
	}

	filinfo, err = gdrive.CvtToFilInfo(fil)
	if err != nil { return nil, fmt.Errorf("CvtToFileInfo: %v", err)}
	return filinfo, nil
}

func (gdrive *GdApiObj) CvtToFilInfo(fil *drive.File) (filinfoptr *FileInfo , err error) {
// A method that creates a FileIndo structure of the file referenced by 'fil'

	var filinfo FileInfo
	filinfo.Id = fil.Id
	filinfo.MimeType = fil.MimeType
	filinfo.Name = fil.Name
	filinfo.Ext = fil.FullFileExtension
	filinfo.ParentName = ""
	filinfo.ParentId = fil.Parents[0]
	filinfo.SingleParent = true
	if len(fil.Parents) > 1 {filinfo.SingleParent = false}
	filinfo.ModTime = fil.ModifiedTime
	filinfo.Size = fil.Size

	return &filinfo, nil
}

func (gdrive *GdApiObj) GetFileByName(nam string) (filesInfo *[]FileInfo, err error) {
// A method that returns a reference to a slice of FileInfo structures which all have the name 'nam'

	var rfileList *drive.FileList

	fields := []googleapi.Field{"nextPageToken, files(id, name, mimeType, parents, modifiedTime)"}

	if len(nam)<1 {return nil, fmt.Errorf("no name provided!")}

	dirs := strings.Split(nam,"/")
/*	fmt.Println("dirs: ", len(dirs))
	for i:=0; i<len(dirs); i++ {
		fmt.Printf("entry: %d dir: %s length: %d\n", i+1, dirs[i], len(dirs[i]))
	}
*/
	if len(dirs) == 1 {
		qstr := fmt.Sprintf("name = '%s'", nam)
//		fmt.Println("qstr: ", qstr)
		rfileList, err = gdrive.GdSvc.Files.List().
		Fields(fields...).
		Q(qstr).Context(gdrive.Ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("Files List: %v", err)
		}
//		if len(rfileList.Files) > 1 {return nil, fmt.Errorf("error gdrive::GetFileByName: multiple files with same name!")}

		fileIdList := make([]FileInfo, 1)
		fileIdList[0].Id = rfileList.Files[0].Id
		fileIdList[0].MimeType = rfileList.Files[0].MimeType
		fileIdList[0].Name = rfileList.Files[0].Name
		return &fileIdList, nil
	}

	// parse first directory
	firstDirStr := dirs[0]
	if len(dirs[0]) < 1 { firstDirStr = "root" }
	qstr := fmt.Sprintf("name = '%s' and mimeType = '%s'", firstDirStr, Gapp["folder"])
//	fmt.Println("qstr file no dir: ", qstr)
	rfileList, err = gdrive.GdSvc.Files.List().Fields(fields...).Q(qstr).Context(gdrive.Ctx).Do()

	if err != nil {	return nil, fmt.Errorf("Files List [by name]: %v", err)}
	if len(rfileList.Files) > 1 {return nil, fmt.Errorf("Files List [by name]: multiple folders with same name!")}

	dirId := rfileList.Files[0].Id
//	fmt.Printf("found dir level 0: %s id: %s\n",dirs[0], dirId) 

	for i:=1; i< len(dirs)-1; i++ {
//		fmt.Println(i, ": ",dirs[i])

		qstr := fmt.Sprintf("(name = '%s' and '%s' in parents) and mimeType = '%s'", dirs[i], dirId, Gapp["folder"])
//		fmt.Println("qstr dir: ", qstr)

		rfileList, err = gdrive.GdSvc.Files.List().Q(qstr).Fields(fields...).Context(gdrive.Ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("Files List: %v", err)
		}
		if len(rfileList.Files) > 1 {return nil, fmt.Errorf("Files List: multiple folders with same name!")}
		if len(rfileList.Files) == 0 {return nil, fmt.Errorf("Files List: folder %s does not exist!", dirs[i])}

		dirId = rfileList.Files[0].Id
//		fmt.Printf("found dir level %d: %s id: %s\n",i, dirs[i], dirId)
	}

	lastDir := len(dirs)-1
	filnam := dirs[lastDir]
//	fmt.Println("lastDir: ", lastDir, " filnam: ", filnam)
	isDir := false
	if len(dirId) < 1 {dirId = "root"}
	if len(filnam) < 1 {
		qstr = fmt.Sprintf("'%s' in parents", dirId)
	} else {
		qstr = fmt.Sprintf("name = '%s' and '%s' in parents", filnam, dirId)
		isDir = true
	}
//	fmt.Println("qstr file: ", qstr)

	rfileList, err = gdrive.GdSvc.Files.List().Fields(fields...).Q(qstr).Context(gdrive.Ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("Files List: %v", err)
	}


	if isDir {
		if len(rfileList.Files) > 1 { return nil, fmt.Errorf("Folders List: multiple folders with same name!")}
		if len(rfileList.Files) == 0 { return nil, fmt.Errorf("Folders List: no folder with name %s found!", filnam)}

		fileIdList := make([]FileInfo, 1)
		fileIdList[0].Id = rfileList.Files[0].Id
		fileIdList[0].MimeType = rfileList.Files[0].MimeType
		fileIdList[0].Name = rfileList.Files[0].Name
		return &fileIdList, nil

	} else {
//		if len(rfileList.Files) > 1 { return nil, fmt.Errorf("error gdrive::GetFileByName: multiple files with same name!")}
		if len(rfileList.Files) == 0 { return nil, fmt.Errorf("File List: no file in folder %s found!", filnam)}
		fileIdList := make([]FileInfo, len(rfileList.Files))
		for i:= 0; i< len(rfileList.Files); i++ {
			fileIdList[i].Id = rfileList.Files[i].Id
			fileIdList[i].MimeType = rfileList.Files[i].MimeType
			fileIdList[i].Name = rfileList.Files[i].Name
		}
		return &fileIdList, nil
	}
}

func (gdrive *GdApiObj) GetFullPath(filId string) (filesInfo *[]FileInfo, path string, err error) {
// A method that returns a slice of FileInfo structures (one for each folder in the path)
// and the full folder path of a file with file id 'filId'

	var gfil *drive.File
	var folders [10]FileInfo
	var fslice []FileInfo

//	fields := []googleapi.Field{"nextPageToken, files(id, name, mimeType, parents, modifiedTime)"}
	fields := []googleapi.Field{"id, name, mimeType, parents, modifiedTime"}

	if len(filId)<1 {return nil, "", fmt.Errorf("GetFullPath: no id provided!")}

	path = ""
	nfilId := filId
	for i:=0; i< 10; i++ {

//		qstr := fmt.Sprintf("(name = '%s' and '%s' in parents) and mimeType = '%s'", dirs[i], dirId, Gapp["folder"])
//		fmt.Println("qstr dir: ", qstr)

		gfil, err = gdrive.GdSvc.Files.Get(nfilId).Fields(fields...).Context(gdrive.Ctx).Do()
		if err != nil {
			return nil, "", fmt.Errorf("GetFullPath: %v", err)
		}
		if gfil == nil {return nil, "", fmt.Errorf("GetFullPath: folder with id %s does not exist!", nfilId)}
		if len(gfil.Parents) > 1 {return nil, "", fmt.Errorf("GetFullPath: file/folder %s has multiple parents!", nfilId)}

		folders[i].Id = gfil.Id
		folders[i].Name = gfil.Name
		folders[i].MimeType = gfil.MimeType
		if i== 0 {
			path = fmt.Sprintf("%s",gfil.Name)
		} else {
			path = fmt.Sprintf("%s/%s",gfil.Name, path)
		}
		if len(gfil.Parents) == 0 {
			fslice = folders[:i]
			return &fslice, path, nil
		}
		nfilId = gfil.Parents[0]

//		fmt.Printf("found folder level %d: \"%-30s\" id: %-35s\n",i, folders[i].Name, nfilId)
	}

	return nil, path, fmt.Errorf("GetFullPath: file/folder path has too many nest levels (>10)!")
}

func (gdrive *GdApiObj) GetFileChar(fid string) (gfil *drive.File, err error) {
// method that returns a file reference for the file with the id 'fid'

	if len(fid)<1 {return nil, fmt.Errorf("no file id provided!")}

	svc := gdrive.GdSvc
	gfil, err = svc.Files.Get(fid).Do()
	if err != nil {
		return nil, fmt.Errorf("Files.Get: %v",err)
	}
	return gfil, nil
}

func (gdrive *GdApiObj) EmptyTrash() (err error) {
// A method that empties the bin that contains the deleted files

	svc := gdrive.GdSvc
	err = svc.Files.EmptyTrash().Do()
	if err != nil {return err}
	return nil
}

func (gdrive *GdApiObj) DumpFileChar(inGfil *drive.File, outfil *os.File) (err error) {
// A method that writes all the file characteristics of an input file 'inGfil' to an output file 'outfil'

	var outstr string

	if inGfil == nil {return fmt.Errorf("no input drive file provided!")	}
	if outfil == nil {return fmt.Errorf("no outfil provided!")	}

	outstr =  fmt.Sprintf("Id:          %s\n", inGfil.Id)
	outstr += fmt.Sprintf("Name:        %s\n", inGfil.Name)
	outstr += fmt.Sprintf("Version:     %d\n", inGfil.Version)
	outstr += fmt.Sprintf("File Ext:    %s\n", inGfil.FileExtension)
	outstr += fmt.Sprintf("Full Ext:    %s\n", inGfil.FullFileExtension)
	outstr += fmt.Sprintf("Description: %s\n", inGfil.Description)
	outstr += fmt.Sprintf("Create Time: %s\n", inGfil.CreatedTime)
	outstr += fmt.Sprintf("Last Viewed: %s\n", inGfil.ViewedByMeTime)
	outstr += fmt.Sprintf("DriveId:     %s\n", inGfil.DriveId)
	outstr += fmt.Sprintf("WebView:     %s\n", inGfil.WebViewLink)
	outstr += fmt.Sprintf("Shared:      %t\n", inGfil.Shared)
	outstr += fmt.Sprintf("Starred:     %t\n", inGfil.Starred)
	outstr += fmt.Sprintf("Owned by me: %t\n", inGfil.OwnedByMe)
	outstr += fmt.Sprintf("Trashed:     %t\n", inGfil.ExplicitlyTrashed)

	outstr += fmt.Sprintf("Size:        %d\n", inGfil.Size)
	outstr += fmt.Sprintf("Mime Type:   %s\n", inGfil.MimeType)

	outstr += fmt.Sprintf("\nOwners:      %d\n", len(inGfil.Owners))
	for i:=0; i< len(inGfil.Owners); i++ {
		outstr += fmt.Sprintf("  owner %d: %s %s\n", i, inGfil.Owners[i].DisplayName, inGfil.Owners[i].EmailAddress)
	}

	outstr += fmt.Sprintf("\nParents:     %d\n", len(inGfil.Parents))
	for i:=0; i< len(inGfil.Parents); i++ {
		outstr += fmt.Sprintf("  parent %d: %s\n", i, inGfil.Parents[i])
	}

	outstr += fmt.Sprintf("\nPermissions:  %d\n", len(inGfil.Permissions))
	for i:=0; i< len(inGfil.Permissions); i++ {
		permit := inGfil.Permissions[i]
		outstr += fmt.Sprintf("  Permission %d: id: %s %s\n", i, permit.Id, permit.DisplayName)
	}

	outstr += fmt.Sprintf("\nPermission Ids:  %d\n", len(inGfil.PermissionIds))
	for i:=0; i< len(inGfil.PermissionIds); i++ {
		outstr += fmt.Sprintf("  Permission Id %d: %s\n", i, inGfil.PermissionIds[i])
	}

	outstr += fmt.Sprintf("\nSpaces:          %d\n", len(inGfil.Spaces))
	for i:=0; i< len(inGfil.Spaces); i++ {
		outstr += fmt.Sprintf("  Space %d: %s\n", i, inGfil.Spaces[i])
	}

	if inGfil.SharingUser != nil {
		outstr += fmt.Sprintf("Sharing User: %s\n", inGfil.SharingUser.EmailAddress)
	} else {
		outstr += "Sharing User: none\n"
	}
	outstr += fmt.Sprintf("Version:   %d\n", inGfil.Version)

	outfil.WriteString(outstr)
	outfil.Close()
	return nil
}

func (gdrive *GdApiObj) ExportFile(inGfil *drive.File, outfil *os.File) (err error) {
// method that exports a file to an outfil
// still todo

	return nil
}

func listExt() {
// function that prints all extensions

	fmt.Println("***** Valid Extensions *********")
	fmt.Printf("\"png\"  mime: image/png\n")
	fmt.Printf("\"jpg\"  mime: image/jpeg\n")
	fmt.Printf("\"pdf\"  mime: application/pdf\n")
	fmt.Printf("\"txt\"  mime: text/plain\n")
	fmt.Printf("\"html\" mime: text/html\n")
	fmt.Printf("\"rtf\"  mime: application/rtf\n")
	fmt.Printf("\"svg\"  mime: image/svg+xml\n")
	fmt.Printf("\"docx\" mime: application/vnd.openxmlformats-officedocument.wordprocessingml.document\n")
	fmt.Printf("\"xlsx\" mime: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet\n")
	fmt.Printf("\"epub\" mime: application/epub+zip\n")
	fmt.Printf("\"pptx\" mime: application/vnd.openxmlformats-officedocument.presentationml.presentation\n")
	fmt.Printf("\"csv\"  mime: text/csv\n")
}

func (gdrive *GdApiObj) ExportFileById(filId string, fileName string, ext string) (err error) {
// A method that exports a file with file id 'filId' to a file with name 'fileName' and extension 'ext'
// The extensions have to be be part of a table

	var mime string

	if !(len(filId) > 0) {return fmt.Errorf("no filId provided!")}
	if !(len(fileName) > 0) {return fmt.Errorf("no fileName provided!")}

	mime, ok := Gapp[ext]
	if !ok {return fmt.Errorf("ext not a valid mime type!")}
/*
	switch ext {
	case "png":
		mime = "image/png"

	case "jpg":
		mime = "image/jpeg"

	case "pdf":
		mime = "application/pdf"

	case "txt":
		mime = "text/plain"

	case "html":
		mime = "text/html"

	case "rtf":
		mime = "application/rtf"

	case "svg":
		mime = "image/svg+xml"

	case "docx":
		mime = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"

	case "xlsx":
		mime = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

	case "epub":
		mime = "application/epub+zip"

	case "pptx":
		mime = "application/vnd.openxmlformats-officedocument.presentationml.presentation"

	case "csv":
		mime = "text/csv"

	default:
		return fmt.Errorf("error ExportFileById: unknown file type: %s!", ext)
	}

*/
	httpResp, err := gdrive.GdSvc.Files.Export(filId, mime).Context(gdrive.Ctx).Download()
	if err != nil {
		return fmt.Errorf("gdrive Export: %v", err)
	}

    defer httpResp.Body.Close()

    if httpResp.StatusCode != 200 {
        return fmt.Errorf("response code invalid: %d", httpResp.StatusCode)
    }

    //Create a empty file
	outFilNam := fmt.Sprintf("%s.%s", fileName, ext)
    outfile, err := os.Create(outFilNam)
    if err != nil {
        return fmt.Errorf("os.Create: %v\n", err)
    }
    defer outfile.Close()

    //Write the bytes to the fiel
    _, err = io.Copy(outfile, httpResp.Body)
    if err != nil {
        return fmt.Errorf("io.Copy: %v\n", err)
    }

	return nil
}

func (gdrive *GdApiObj) DownloadFileById(filId string, fileName string) (err error) {
// A method that downloads a file with the id 'filid' to a file with the name 'fileName'

	if !(len(filId) > 0) {return fmt.Errorf("error gdrive::DownloadFileById: no id provided!")}
	if !(len(fileName) > 0) {return fmt.Errorf("error gdrive::DownloadFileById: no file name provided!")}


	httpResp, err := gdrive.GdSvc.Files.Get(filId).Context(gdrive.Ctx).Download()
	if err != nil {
		return fmt.Errorf("error gdrive::DownloadFileById: %v", err)
	}
	fmt.Println("downloading!")
    defer httpResp.Body.Close()

    if httpResp.StatusCode != 200 {
        return fmt.Errorf("error gdrive::DownloadFileById: Received non 200 response code")
    }
	fmt.Printf("resp: %d\n", httpResp.StatusCode)
    //Create a empty file
//	path := "output/" + fileName
    outfil, err := os.Create(fileName)
    if err != nil {
        return fmt.Errorf("error gdrive::DownloadFileById: could not create dest file! %v\n", err)
    }
    defer outfil.Close()
	fmt.Printf("created file %s!\n", fileName)
    //Write the bytes to the file
    _, err = io.Copy(outfil, httpResp.Body)
    if err != nil {
        return fmt.Errorf("error gdrive::DownloadFileById: could not copy to dest file! %v\n", err)
    }
	fmt.Println("success copied!")
	return nil
}

func PrintDriveFile(title string, fil *drive.File) {

	fmt.Printf("******* file details: %s *********\n", title)
	fmt.Printf("Name:  %s\n", fil.Name)
	fmt.Printf("Mime:  %s\n", fil.MimeType)
	fmt.Printf("Id:    %s\n", fil.Id)
	fmt.Printf("Parents: %d\n", len(fil.Parents))
	for i:=0; i< len(fil.Parents); i++ {
		fmt.Printf(" parent [%d]: %s\n", i, fil.Parents[i])
	}
	fmt.Printf("Size:    %d\n", fil.Size)
	fmt.Println("******* end file details *********")

}
