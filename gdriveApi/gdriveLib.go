// gdriveLib.go
// golang library for google drive file management
//
// author: prr, azul software
// date: 30/1/2022
// update: 7/6/2022
// copywrite 2022 prr, azul software
//

package gdriveLib

import (
        "context"
        "encoding/json"
        "fmt"
		"io"
        "io/ioutil"
		"strings"
        "net/http"
        "os"
        "golang.org/x/oauth2"
        "golang.org/x/oauth2/google"
        "google.golang.org/api/drive/v3"
        "google.golang.org/api/option"
		"google.golang.org/api/googleapi"
)

type GdriveApiStruct  struct {
	Ctx context.Context
	Svc *drive.Service
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

// Retrieves a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "/home/peter/go/src/google/gdrive/tokGdrive.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Requests a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
        authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
        fmt.Printf("Go to the following link in your browser then type the "+
                "authorization code: \n%v\n", authURL)

        var authCode string
        if _, err := fmt.Scan(&authCode); err != nil {
			fmt.Println("Unable to read authorization code: ", err)
			os.Exit(1)
        }

        tok, err := config.Exchange(oauth2.NoContext, authCode)
        if err != nil {
			fmt.Println("Unable to retrieve token from web: ", err)
			os.Exit(1)
        }
        return tok
}

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

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		fmt.Println("Unable to cache OAuth token: ", err)
		os.Exit(1)
	}
	json.NewEncoder(f).Encode(token)
}

func (gdrive *GdriveApiStruct) CreDumpFile(fid string, filnam string)(err error) {
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
		return fmt.Errorf("error gdrive::CreDumpFile: sub-dir \"output\" does not exist!")
	}
	if err != nil {
		return fmt.Errorf("error gdrive::CreDumpFile: %v \n", err)
	}
	if !filinfo.IsDir() {
		return fmt.Errorf("error gdrive::CreDumpFile -- file \"output\" is not a directory!")
	}

	path:= "output/" + string(nfilnam) + "." + ext
//	fmt.Println("path: ",path)
	outfil, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("error gdrive::CreDumpFile: cannot open output file: %v \n", err)
	}

	// get file attributes
	svc := gdrive.Svc
	gfil, err := svc.Files.Get(fid).Do()
	if err != nil {
		return fmt.Errorf("error gdrive::CreDumpFile: cannot get file with id: %s! %v", fid, err)
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

//nnn
func (gdObj *GdriveApiStruct) InitDriveApi() (svc *drive.Service, err error) {
        ctx := context.Background()
		gdObj.Ctx = ctx
        b, err := ioutil.ReadFile("/home/peter/go/src/google/gdrive/credGdrive.json")
        if err != nil {
			return fmt.Errorf("Unable to read client secret file: %v!", err)
		}

        // If modifying these scopes, delete your previously saved token.json.
        config, err := google.ConfigFromJSON(b, drive.DriveScope)
        if err != nil {
			return fmt.Errorf("Unable to parse client secret file to config: %v!", err)
        }

        client := getClient(config)

        svc, err := drive.NewService(ctx, option.WithHTTPClient(client))
        if err != nil {
			return fmt.Errorf("Unable to retrieve Drive client: %v !", err)
        }
	return svc, nil
}

func (gdrive *GdriveApiStruct) Init() (err error) {
        ctx := context.Background()
		gdrive.Ctx = ctx
        b, err := ioutil.ReadFile("/home/peter/go/src/google/gdrive/credGdrive.json")
        if err != nil {
			return fmt.Errorf("Unable to read client secret file: %v!", err)
		}

        // If modifying these scopes, delete your previously saved token.json.
        config, err := google.ConfigFromJSON(b, drive.DriveScope)
        if err != nil {
			return fmt.Errorf("Unable to parse client secret file to config: %v!", err)
        }

        client := getClient(config)

        svc, err := gdrive.NewService(ctx, option.WithHTTPClient(client))
        if err != nil {
			return fmt.Errorf("Unable to retrieve Drive client: %v !", err)
        }
		gdrive.Svc = svc
	return nil
}

func (gdrive *GdriveApiStruct) GetAbout() (resp *drive.About, err error) {

	fields := []googleapi.Field{"nextPageToken, files(id, name, mimeType, parents, modifiedTime)"}

	svc := gdrive.Svc
    resp, err = svc.About.Get().Fields(fields...).Do()
    if err != nil {
        fmt.Println("error svc.about.get:", err)
        return nil, fmt.Errorf("error GetAbout: %v", err)
    }

	return resp, nil
}

func (gdrive *GdriveApiStruct) DumpAbout(about *drive.About, outfil *os.File) (err error) {
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


func (gdrive *GdriveApiStruct) ListFiles() (fileList []*drive.File, err error) {

	fields := []googleapi.Field{"nextPageToken, files(id, name, mimeType, parents, modifiedTime)"}

	pagetoken := ""
	fin := false

	qstr := fmt.Sprintf("mimeType != '%s' and name = 'root'", Gapp["folder"])

//	fmt.Println("qstr: ", qstr)
	for i:=0; i< 3; i++ {
		nfileList, err := gdrive.Svc.Files.List().
		Fields(fields...).
		Q(qstr).
		PageToken(pagetoken).Context(gdrive.Ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("error gdrive::ListFile %d %v", i, err)
		}
		fmt.Println("files: ", len(nfileList.Files))
		if len(nfileList.Files) < 1 {
			return nil, fmt.Errorf("error gdrive:: ListFile - no files found!")
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
		return fileList, fmt.Errorf("error gdrive::ListFile -- too many files > 1000!")
	}
	return fileList, nil
}

func (gdrive *GdriveApiStruct) ListAllFiles(dirId string) (fileList []*drive.File, err error) {

	fields := []googleapi.Field{"nextPageToken, files(id, name, mimeType, parents, modifiedTime)"}

	pagetoken := ""
	fin := false


	qstr := fmt.Sprintf("mimeType != '%s' and '%s' in parent", Gapp["folder"], dirId)

//	fmt.Println("qstr: ", qstr)
	for i:=0; i< 3; i++ {
		nfileList, err := gdrive.Svc.Files.List().
		Fields(fields...).
		Q(qstr).
		PageToken(pagetoken).Context(gdrive.Ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("error gdrive::ListFile %d %v", i, err)
		}
		fmt.Println("files: ", len(nfileList.Files))
		if len(nfileList.Files) < 1 {
			return nil, fmt.Errorf("error gdrive:: ListFile - no files found!")
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
		return fileList, fmt.Errorf("error gdrive::ListFile -- too many files > 1000!")
	}


	return fileList, nil
}


func (gdrive *GdriveApiStruct) ListFilesByName(nam string, dirId string) (filList []*drive.File, err error) {
//	var filInfList []FileInfo
	var qstr string

	fields := []googleapi.Field{"nextPageToken, files(id, name, mimeType, parents, modifiedTime)"}

	pagetoken := ""
	fin := false
	if len(dirId) > 0 {
		_, err = gdrive.Svc.Files.Get(dirId).Context(gdrive.Ctx).Do()
		if err != nil {return filList, fmt.Errorf("error gdrive::CopyFile: could not find folder with id: %s -- %v", dirId, err)}
		qstr = fmt.Sprintf("(mimeType != '%s' and name = '%s') and '%s' in parents", Gapp["folder"], nam, dirId)
	} else {
		qstr = fmt.Sprintf("mimeType != '%s' and name = '%s'", Gapp["folder"], nam)
	}

//	fmt.Println("qstr: ", qstr)
	for i:=0; i< 3; i++ {
		nfileList, err := gdrive.Svc.Files.List().
		Fields(fields...).
		Q(qstr).
		PageToken(pagetoken).Context(gdrive.Ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("error gdrive::ListFile %d %v", i, err)
		}
		fmt.Println("files: ", len(nfileList.Files))
		if len(nfileList.Files) < 1 {
			return nil, fmt.Errorf("error gdrive:: ListFile - no files found!")
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
		return filList, fmt.Errorf("error gdrive::ListFile -- too many files > 1000!")
	}
	return filList, nil
}

func (gdrive *GdriveApiStruct) ListFoldersByName(nam string) (filList []*drive.File, err error) {

	fields := []googleapi.Field{"nextPageToken, files(id, name, mimeType, parents, modifiedTime)"}

	pagetoken := ""
	fin := false
	qstr := fmt.Sprintf("mimeType = '%s' and name = '%s'", Gapp["folder"], nam)

//	fmt.Println("qstr: ", qstr)
	for i:=0; i< 3; i++ {
		nfileList, err := gdrive.Svc.Files.List().
		Fields(fields...).
		Q(qstr).
		PageToken(pagetoken).Context(gdrive.Ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("error gdrive::ListFile %d %v", i, err)
		}
		fmt.Println("files: ", len(nfileList.Files))
		if len(nfileList.Files) < 1 {
			return nil, fmt.Errorf("error gdrive:: ListFile - no files found!")
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
		return filList, fmt.Errorf("error gdrive::ListFile -- too many files > 1000!")
	}
	return filList, nil
}

func (gdrive *GdriveApiStruct) ListFolderByName(nam string) (folderList *[]FileInfo, err error) {

	if len(nam) < 1 { return nil, fmt.Errorf("error gdrive::ListFolderByName -- no name provided!")}

	fields := []googleapi.Field{"nextPageToken, files(id, name, mimeType, parents, modifiedTime)"}

	pagetoken := ""
	qstr := fmt.Sprintf("mimeType = '%s' and name = '%s'", Gapp["folder"], nam)

//	fmt.Println("qstr: ", qstr)
	nfileList, err := gdrive.Svc.Files.List().
		Fields(fields...).
		Q(qstr).
		PageToken(pagetoken).Context(gdrive.Ctx).Do()
	if err != nil {return nil, -1, fmt.Errorf("error gdrive::ListFolderByName: get list: %v", err)}
	numFolders = len(nfileList.Files)
//	fmt.Println("folders: ", numFolders)

	if len(pagetoken) > 0 {return nil, fmt.Errorf("error gdrive::ListFolderByName: too many folders (>100)!")}

	finfolist := make([]FileInfo, numFolders)
	for i:= 0; i< numFolders; i++ {
		fileinfo, _ := gdrive.CvtToFilInfo(nfileList.Files[i])
		finfolist[i] = *fileinfo
	}
	folderList = &finfolist
	return folderList, nil
}

func (gdrive *GdriveApiStruct) ListFFByName(nam string) (filList []*drive.File, err error) {

	fields := []googleapi.Field{"nextPageToken, files(id, name, mimeType, parents, modifiedTime)"}

	pagetoken := ""
	fin := false
	qstr := fmt.Sprintf("name = '%s'", nam)

//	fmt.Println("qstr: ", qstr)
	for i:=0; i< 3; i++ {
		nfileList, err := gdrive.Svc.Files.List().
		Fields(fields...).
		Q(qstr).
		PageToken(pagetoken).Context(gdrive.Ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("error gdrive::ListFile %d %v", i, err)
		}
		fmt.Println("files: ", len(nfileList.Files))
		if len(nfileList.Files) < 1 {
			return nil, fmt.Errorf("error gdrive:: ListFile - no files found!")
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
		return filList, fmt.Errorf("error gdrive::ListFile -- too many files > 1000!")
	}
	return filList, nil
}

func (gdrive *GdriveApiStruct) ListFilesBySize(foldId string, minSize int64) (filList []*drive.File, err error) {
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
		nfileList, err := gdrive.Svc.Files.List().
		Fields(fields...).
		Q(qstr).
		PageToken(pagetoken).Context(gdrive.Ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("error gdrive::ListFile %d %v", i, err)
		}
		fmt.Println("files: ", len(nfileList.Files))
		if len(nfileList.Files) < 1 {
			return nil, fmt.Errorf("error gdrive:: ListFile - no files found!")
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
		return filList, fmt.Errorf("error gdrive::ListFile -- too many files > 1000!")
	}

	return filList, nil
}

func (gdrive *GdriveApiStruct) ListTopDir() (id string, err error) {
	var topDir *drive.File
	topDir, err = gdrive.Svc.Files.Get("root").Context(gdrive.Ctx).Do()
	if err != nil {
		return "", fmt.Errorf("error gdrive::ListTopDir: %v", err)
	}
	return topDir.Id, nil
}



func (gdrive *GdriveApiStruct) CopyFile(filId string, nam string, dirId string) (nfilId string, err error) {
// needs work
	var fil, nfil *drive.File
	var par [1]string
	if len(filId) < 1 {return "", fmt.Errorf("error gdrive::CopyFile -- no filId string!")}
	if len(nam) < 1 {return "", fmt.Errorf("error gdrive::CopyFile -- no nam string!")}

	if len(dirId) > 0 {
		_, err = gdrive.Svc.Files.Get(dirId).Context(gdrive.Ctx).Do()
		if err != nil {return "",fmt.Errorf("error gdrive::CopyFile: could not find folder with id: %s -- %v", dirId, err)}
	}

	fil, err = gdrive.Svc.Files.Get(filId).Context(gdrive.Ctx).Do()
	if err != nil {return "",fmt.Errorf("error gdrive::CopyFile: could not find file with id: %s -- %v", filId, err)}
	par[0] = dirId
	fil.Parents = par[:]
	fil.Name = nam
	nfil, err = gdrive.Svc.Files.Copy(filId, fil).Context(gdrive.Ctx).Do()
	if err != nil {
		return "", fmt.Errorf("error gdrive::CopyFile: %v", err)
	}
	return nfil.Id, nil
}

func (gdrive *GdriveApiStruct) CreateFile(pDirId string, nam string) (fileId string, err error) {
// needs work
	var fil drive.File
	var dir *drive.File
	var par [1]string
	if len(pDirId) < 1 {return "", fmt.Errorf("error gdrive::CreateDir -- no pDirId string!")}
	if len(nam) < 1 {return "", fmt.Errorf("error gdrive::CreateDir -- no nam string!")}

	// we could check nam for invalid chars

	par[0] = pDirId
	fil.Parents = par[:]
	fil.Name = nam
	fil.MimeType = Gapp["folder"]

	dir, err = gdrive.Svc.Files.Create(&fil).Context(gdrive.Ctx).Do()
	if err != nil {
		return "", fmt.Errorf("error gdrive::CreateDir: %v", err)
	}
	return dir.Id, nil
}

func (gdrive *GdriveApiStruct) CreateFolder(pDirId string, nam string) (folderId string, err error) {
	var fil drive.File
	var dir *drive.File
	var par [1]string
	if len(pDirId) < 1 {return "", fmt.Errorf("error gdrive::CreateDir -- no pDirId string!")}
	if len(nam) < 1 {return "", fmt.Errorf("error gdrive::CreateDir -- no nam string!")}

	// we could check nam for invalid chars

	par[0] = pDirId
	fil.Parents = par[:]
	fil.Name = nam
	fil.MimeType = Gapp["folder"]

	dir, err = gdrive.Svc.Files.Create(&fil).Context(gdrive.Ctx).Do()
	if err != nil {
		return "", fmt.Errorf("error gdrive::CreateDir: %v", err)
	}
	return dir.Id, nil
}

func (gdrive *GdriveApiStruct) DeleteFileById(filId string) (err error) {

	if len(filId)<1 {return fmt.Errorf("error gdrive::DeleteFileById: no file id provided!")}

	err = gdrive.Svc.Files.Delete(filId).Context(gdrive.Ctx).Do()
	if err != nil {
		return fmt.Errorf("error gdrive::DeleteFolderById: %v", err)
	}
	return nil
}

func (gdrive *GdriveApiStruct) DeleteFileByName(nam string) (err error) {

	if len(nam)<1 {return fmt.Errorf("error gdrive::DeleteFileByName: no name provided!")}

	fid:= "abc"

	err = gdrive.Svc.Files.Delete(fid).Context(gdrive.Ctx).Do()
	if err != nil {
		return fmt.Errorf("error gdrive::DeleteFolderById: %v", err)
	}
	return nil
}



func (gdrive *GdriveApiStruct) FetchFileById(fid string) (resp *http.Response, err error) {

	if len(fid)<1 {return nil, fmt.Errorf("error gdrive::GetFiles: no nam provided!")}

	resp, err = gdrive.Svc.Files.Get(fid).Context(gdrive.Ctx).Download()
	if err != nil {
		return nil, fmt.Errorf("error gdrive::GetFile Download: %v", err)
	}
	return resp, nil
}

func (gdrive *GdriveApiStruct) MoveFileById(filId string, dirId string) (err error) {
	var fil, updfil *drive.File
	var parentStr string

	if len(filId)<1 {return fmt.Errorf("error gdrive::MoveFileById: no filId provided!")}
	if len(dirId)<1 {return fmt.Errorf("error gdrive::MoveFileById: no dirId provided!")}
	fil, err = gdrive.Svc.Files.Get(filId).Context(gdrive.Ctx).Do()
	if err != nil {return fmt.Errorf("error gdrive::MoveFileById: could not find file with id: %s -- %v", filId, err)}
	_, err = gdrive.Svc.Files.Get(dirId).Context(gdrive.Ctx).Do()
	if err != nil {return fmt.Errorf("error gdrive::MoveFileById: could not find folder with id: %s -- %v", dirId, err)}

	// remove old parents first, if they exist
	updfil = fil
	updFilId := filId
	if len(fil.Parents) > 0 {
		parentStr = fil.Parents[0]
		for i:=1; i < len(fil.Parents); i++ {
			parentStr += "," + fil.Parents[i]
		}

		updfil, err = gdrive.Svc.Files.Update(filId, fil).RemoveParents(parentStr).Context(gdrive.Ctx).Do()
		if err != nil {	return fmt.Errorf("error gdrive::GetFilebyId: could not remove parents: %v", err)}
		updFilId = updfil.Id
	}

	updfil, err = gdrive.Svc.Files.Update(updFilId, updfil).AddParents(dirId).Context(gdrive.Ctx).Do()
	if err != nil {
		return fmt.Errorf("error gdrive::GetFilebyId: could not remove parents: %v", err)
	}

	return nil
}

func (gdrive *GdriveApiStruct) GetFileById(filId string) (fil *drive.File, err error) {

	if len(filId)<1 {return nil, fmt.Errorf("error gdrive::GetFiles: no nam provided!")}

	fil, err = gdrive.Svc.Files.Get(filId).Context(gdrive.Ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("error gdrive::GetFile Download: %v", err)
	}
	return fil, nil
}

func (gdrive *GdriveApiStruct) GetFileInfoById(filId string) (filinfo *FileInfo , err error) {

	fields := []googleapi.Field{"id, name, mimeType, parents, fullFileExtension, modifiedTime, size"}

	if len(filId)<1 {return nil, fmt.Errorf("error gdrive::GetFileInfoById: no filId provided!")}

	fil, err := gdrive.Svc.Files.Get(filId).Fields(fields...).Context(gdrive.Ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("error gdrive::GetFile Download: %v", err)
	}

	filinfo, err = gdrive.CvtToFilInfo(fil)
	if err != nil { return nil, fmt.Errorf("error gdrive::GetFileInfoById: GetFilInfo -- %v", err)}
	return filinfo, nil
}

func (gdrive *GdriveApiStruct) CvtToFilInfo(fil *drive.File) (filinfoptr *FileInfo , err error) {

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

func (gdrive *GdriveApiStruct) GetFileByName(nam string) (filesInfo *[]FileInfo, err error) {
	var rfileList *drive.FileList

	fields := []googleapi.Field{"nextPageToken, files(id, name, mimeType, parents, modifiedTime)"}

	if len(nam)<1 {return nil, fmt.Errorf("error gdrive::GetFileByName: no name provided!")}

	dirs := strings.Split(nam,"/")
/*	fmt.Println("dirs: ", len(dirs))
	for i:=0; i<len(dirs); i++ {
		fmt.Printf("entry: %d dir: %s length: %d\n", i+1, dirs[i], len(dirs[i]))
	}
*/
	if len(dirs) == 1 {
		qstr := fmt.Sprintf("name = '%s'", nam)
//		fmt.Println("qstr: ", qstr)
		rfileList, err = gdrive.Svc.Files.List().
		Fields(fields...).
		Q(qstr).Context(gdrive.Ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("error gdrive::GetFilebyName: %v", err)
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
	rfileList, err = gdrive.Svc.Files.List().Fields(fields...).Q(qstr).Context(gdrive.Ctx).Do()

	if err != nil {	return nil, fmt.Errorf("error gdrive::GetFilebyName: %v", err)}
	if len(rfileList.Files) > 1 {return nil, fmt.Errorf("error gdrive::GetFileByName: multiple folders with same name!")}

	dirId := rfileList.Files[0].Id
//	fmt.Printf("found dir level 0: %s id: %s\n",dirs[0], dirId) 

	for i:=1; i< len(dirs)-1; i++ {
//		fmt.Println(i, ": ",dirs[i])

		qstr := fmt.Sprintf("(name = '%s' and '%s' in parents) and mimeType = '%s'", dirs[i], dirId, Gapp["folder"])
//		fmt.Println("qstr dir: ", qstr)

		rfileList, err = gdrive.Svc.Files.List().Q(qstr).Fields(fields...).Context(gdrive.Ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("error gdrive::GetFilebyName: %v", err)
		}
		if len(rfileList.Files) > 1 {return nil, fmt.Errorf("error gdrive::GetFileByName: multiple folders with same name!")}
		if len(rfileList.Files) == 0 {return nil, fmt.Errorf("error gdrive::GetFileByName: folder %s does not exist!", dirs[i])}

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

	rfileList, err = gdrive.Svc.Files.List().Fields(fields...).Q(qstr).Context(gdrive.Ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("error gdrive::GetFilebyName: %v", err)
	}


	if isDir {
		if len(rfileList.Files) > 1 { return nil, fmt.Errorf("error gdrive::GetFileByName: multiple folders with same name!")}
		if len(rfileList.Files) == 0 { return nil, fmt.Errorf("error gdrive::GetFileByName: no folder with name %s found!", filnam)}

		fileIdList := make([]FileInfo, 1)
		fileIdList[0].Id = rfileList.Files[0].Id
		fileIdList[0].MimeType = rfileList.Files[0].MimeType
		fileIdList[0].Name = rfileList.Files[0].Name
		return &fileIdList, nil

	} else {
//		if len(rfileList.Files) > 1 { return nil, fmt.Errorf("error gdrive::GetFileByName: multiple files with same name!")}
		if len(rfileList.Files) == 0 { return nil, fmt.Errorf("error gdrive::GetFileByName: no file in folder %s found!", filnam)}
		fileIdList := make([]FileInfo, len(rfileList.Files))
		for i:= 0; i< len(rfileList.Files); i++ {
			fileIdList[i].Id = rfileList.Files[i].Id
			fileIdList[i].MimeType = rfileList.Files[i].MimeType
			fileIdList[i].Name = rfileList.Files[i].Name
		}
		return &fileIdList, nil
	}
}

func (gdrive *GdriveApiStruct) GetFullPath(filId string) (filesInfo *[]FileInfo, path string, err error) {
//	var rfileList *drive.FileList
	var gfil *drive.File
	var folders [10]FileInfo
	var fslice []FileInfo

//	fields := []googleapi.Field{"nextPageToken, files(id, name, mimeType, parents, modifiedTime)"}
	fields := []googleapi.Field{"id, name, mimeType, parents, modifiedTime"}

	if len(filId)<1 {return nil, "", fmt.Errorf("error gdrive::GetFullPath: no id provided!")}

	path = ""
	nfilId := filId
	for i:=0; i< 10; i++ {

//		qstr := fmt.Sprintf("(name = '%s' and '%s' in parents) and mimeType = '%s'", dirs[i], dirId, Gapp["folder"])
//		fmt.Println("qstr dir: ", qstr)

		gfil, err = gdrive.Svc.Files.Get(nfilId).Fields(fields...).Context(gdrive.Ctx).Do()
		if err != nil {
			return nil, "", fmt.Errorf("error gdrive::GetFullPath: %v", err)
		}
		if gfil == nil {return nil, "", fmt.Errorf("error gdrive::GetFullPath: folder with id %s does not exist!", nfilId)}
		if len(gfil.Parents) > 1 {return nil, "", fmt.Errorf("error gdrive::GetFullPath: file/folder %s has multiple parents!", nfilId)}

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

	return nil, path, fmt.Errorf("error gdrive::GetFullPath: file/folder path has too many levels (>10)!")
}

func (gdrive *GdriveApiStruct) GetFileChar(fid string) (gfil *drive.File, err error) {

	if len(fid)<1 {return nil, fmt.Errorf("error gdrive::GetFileChar: no file id provided!")}

	svc := gdrive.Svc
	gfil, err = svc.Files.Get(fid).Do()
	if err != nil {
		return nil, fmt.Errorf("error gdrive::GetFileChar: %v",err)
	}
	return gfil, nil
}

func (gdrive *GdriveApiStruct) EmptyTrash() (err error) {

	svc := gdrive.Svc
	err = svc.Files.EmptyTrash().Do()
	if err != nil {
		return fmt.Errorf("error gdrive::EmptyTrash %v",err)
	}
	return nil

}

func (gdrive *GdriveApiStruct) DumpFileChar(inGfil *drive.File, outfil *os.File) (err error) {
	var outstr string

	if inGfil == nil {return fmt.Errorf("error gdrive::DumpFileChar: no input drive file provided!")	}
	if outfil == nil {return fmt.Errorf("error gdrive::DumpFileChar: no outfil provided!")	}

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

func (gdrive *GdriveApiStruct) ExportFile(inGfil *drive.File, outfil *os.File) (err error) {

	return nil
}

func (gdrive *GdriveApiStruct) ExportFileById(filId string, fileName string, ext string) (err error) {
	var mime string

	if !(len(filId) > 0) {return fmt.Errorf("error gdrive::GetFiles: no id provided!")}
	if !(len(fileName) > 0) {return fmt.Errorf("error gdrive::GetFiles: no file name provided!")}

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


	httpResp, err := gdrive.Svc.Files.Export(filId, mime).Context(gdrive.Ctx).Download()
	if err != nil {
		return fmt.Errorf("error gdrive::GetFile Download: %v", err)
	}

    defer httpResp.Body.Close()

    if httpResp.StatusCode != 200 {
        return fmt.Errorf("error gdrive::ExportFile: Received non 200 response code")
    }

    //Create a empty file
	outFilNam := fmt.Sprintf("%s.%s", fileName, ext)
    outfile, err := os.Create(outFilNam)
    if err != nil {
        return fmt.Errorf("error gdrive::ExportFile: could not create dest file! %v\n", err)
    }
    defer outfile.Close()

    //Write the bytes to the fiel
    _, err = io.Copy(outfile, httpResp.Body)
    if err != nil {
        return fmt.Errorf("error gdrive::ExportFile: could not copy to dest file! %v\n", err)
    }

	return nil
}

func (gdrive *GdriveApiStruct) DownloadFileById(filId string, fileName string) (err error) {

	if !(len(filId) > 0) {return fmt.Errorf("error gdrive::DownloadFileById: no id provided!")}
	if !(len(fileName) > 0) {return fmt.Errorf("error gdrive::DownloadFileById: no file name provided!")}


	httpResp, err := gdrive.Svc.Files.Get(filId).Context(gdrive.Ctx).Download()
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
