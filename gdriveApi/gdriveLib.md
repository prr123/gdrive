<p style="font-size:18pt; text-align:center;">gdriveLib</p>

# Description
   gdriveLib.go   
   golang library for google drive file management   
     
   author: prr, azul software   
   date: 30/1/2022   
   update: 7/6/2022   
   copywrite 2022 prr, azul software   
     


# Types
## GdApiObj    
type GdApiObj  struct     

## FileInfo    
type FileInfo struct     


# Functions
## ListApps    
func ListApps()     

 function that lists all available apps    
## getClient    
func getClient(config *oauth2.Config) *http.Client     

## getTokenFromWeb    
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token     

## tokenFromFile    
func tokenFromFile(file string) (*oauth2.Token, error)     

## saveToken    
func saveToken(path string, token *oauth2.Token)     

## listExt    
func listExt()     

 function that prints all extensions    

# Methods
## gdrive *GdApiOb: CreDumpFile    
func (gdrive *GdApiObj) CreDumpFile(fid string, filnam string)(err error)     

 function that creates a text file to dump document file    
## gdObj *GdApiOb: InitDriveApi    
func (gdObj *GdApiObj) InitDriveApi() (err error)     

 method that initialises the GdriveApi structure and returns a  service pointer    
## gdObj *GdApiOb: Init    
func (gdObj *GdApiObj) Init() (err error)     

 method that initialise the Gdrive Object. Ther service pointer is assigned to the GdApiObj    
 better to use the method InitDriveApi    
## gdrive *GdApiOb: GetAbout    
func (gdrive *GdApiObj) GetAbout() (resp *drive.About, err error)     

 method that lists all files in a drive    
## gdrive *GdApiOb: DumpAbout    
func (gdrive *GdApiObj) DumpAbout(about *drive.About, outfil *os.File) (err error)     

 method that writes results from an about query into the output file    
## gdrive *GdApiOb: ListFiles    
func (gdrive *GdApiObj) ListFiles() (fileList []*drive.File, err error)     

 method that lists all folders    
## gdrive *GdApiOb: ListAllFiles    
func (gdrive *GdApiObj) ListAllFiles(dirId string) (fileList []*drive.File, err error)     

 method that lists all files in a folder with id dirId    
## gdrive *GdApiOb: ListFilesByName    
func (gdrive *GdApiObj) ListFilesByName(nam string, dirId string) (filList []*drive.File, err error)     

 method that lists all files with name 'nam' and folder id 'dirId'    
## gdrive *GdApiOb: ListFoldersByName    
func (gdrive *GdApiObj) ListFoldersByName(nam string) (filList []*drive.File, err error)     

 method that looks for all folders with name 'nam' and returns a splice of file pointers    
## gdrive *GdApiOb: ListFolderByName    
func (gdrive *GdApiObj) ListFolderByName(nam string) (folderList *[]FileInfo, err error)     

 method that looks for all sub folders of folder with name 'nam' and returns a list of folder pointers    
## gdrive *GdApiOb: ListFFByName    
func (gdrive *GdApiObj) ListFFByName(nam string) (filList []*drive.File, err error)     

 method that looks for all files with name 'nam'    
## gdrive *GdApiOb: ListFilesBySize    
func (gdrive *GdApiObj) ListFilesBySize(foldId string, minSize int64) (filList []*drive.File, err error)     

 method that lists all files above a certain size    
## gdrive *GdApiOb: ListTopDir    
func (gdrive *GdApiObj) ListTopDir() (id string, err error)     

 method that provides the id of the root folder    
## gdrive *GdApiOb: CopyFile    
func (gdrive *GdApiObj) CopyFile(filId string, nam string, dirId string) (nfilId string, err error)     

 method that copies a file with id 'filId' to new file. The method returns the id of the new file 'nfilId    
 still todo    
 assign same parent id to new file    
## gdrive *GdApiOb: CreateFile    
func (gdrive *GdApiObj) CreateFile(pDirId string, nam string) (fileId string, err error)     

 A method that creates a file with the parent id 'pDirId' and name 'nam'.     
 The method return the file id of the created file    
## gdrive *GdApiOb: CreateFolder    
func (gdrive *GdApiObj) CreateFolder(pDirId string, nam string) (folderId string, err error)     

 A method that creates a folder under the parent folder with id 'pDirId' and name 'nam'.    
 The method returns the  file id 'folderId' of the newly created Folder.    
## gdrive *GdApiOb: DeleteFileById    
func (gdrive *GdApiObj) DeleteFileById(filId string) (err error)     

 A method that deletes a file with the id 'filId'    
## gdrive *GdApiOb: DeleteFileByName    
func (gdrive *GdApiObj) DeleteFileByName(nam string) (err error)     

 A method that deletes a file indendified by name 'nam'    
 not finished    
## gdrive *GdApiOb: FetchFileById    
func (gdrive *GdApiObj) FetchFileById(fid string) (resp *http.Response, err error)     

 A method that downloads a file identified by id 'fil' which returns the file in the http response body    
## gdrive *GdApiOb: MoveFileById    
func (gdrive *GdApiObj) MoveFileById(filId string, dirId string) (err error)     

 A method that moves the file with file id 'filId' into a dierectory with the id 'dirId'    
## gdrive *GdApiOb: GetFileById    
func (gdrive *GdApiObj) GetFileById(filId string) (fil *drive.File, err error)     

 A method that returns a file with the id 'filId'    
## gdrive *GdApiOb: GetFileInfoById    
func (gdrive *GdApiObj) GetFileInfoById(filId string) (filinfo *FileInfo , err error)     

 A method that returns a pointer to the file info struct 'FileInfo' of the file with id 'filId'    
## gdrive *GdApiOb: CvtToFilInfo    
func (gdrive *GdApiObj) CvtToFilInfo(fil *drive.File) (filinfoptr *FileInfo , err error)     

 A method that creates a FileIndo structure of the file referenced by 'fil'    
## gdrive *GdApiOb: GetFileByName    
func (gdrive *GdApiObj) GetFileByName(nam string) (filesInfo *[]FileInfo, err error)     

 A method that returns a reference to a slice of FileInfo structures which all have the name 'nam'    
## gdrive *GdApiOb: GetFullPath    
func (gdrive *GdApiObj) GetFullPath(filId string) (filesInfo *[]FileInfo, path string, err error)     

 A method that returns a slice of FileInfo structures (one for each folder in the path)    
 and the full folder path of a file with file id 'filId'    
## gdrive *GdApiOb: GetFileChar    
func (gdrive *GdApiObj) GetFileChar(fid string) (gfil *drive.File, err error)     

 method that returns a file reference for the file with the id 'fid'    
## gdrive *GdApiOb: EmptyTrash    
func (gdrive *GdApiObj) EmptyTrash() (err error)     

 A method that empties the bin that contains the deleted files    
## gdrive *GdApiOb: DumpFileChar    
func (gdrive *GdApiObj) DumpFileChar(inGfil *drive.File, outfil *os.File) (err error)     

 A method that writes all the file characteristics of an input file 'inGfil' to an output file 'outfil'    
## gdrive *GdApiOb: ExportFile    
func (gdrive *GdApiObj) ExportFile(inGfil *drive.File, outfil *os.File) (err error)     

 method that exports a file to an outfil    
 still todo    
## gdrive *GdApiOb: ExportFileById    
func (gdrive *GdApiObj) ExportFileById(filId string, fileName string, ext string) (err error)     

 A method that exports a file with file id 'filId' to a file with name 'fileName' and extension 'ext'    
 The extensions have to be be part of a table    
## gdrive *GdApiOb: DownloadFileById    
func (gdrive *GdApiObj) DownloadFileById(filId string, fileName string) (err error)     

 A method that downloads a file with the id 'filid' to a file with the name 'fileName'    
