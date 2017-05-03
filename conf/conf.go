package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

const (
	defaultConfFileName   = "conf"
	defaultBleveFolder    = "pictures.db"
	defaultFilesPath      = "files"
	defaultBoltFileName   = "bolt.db"
	defaultSQLiteFileName = "sqlite.db"
)

//Options ...
var Options struct {
	//MapQuestAPIKey API key for map quest service
	MapQuestAPIKey string `json:"map_quest_api_key"`
	//SourceFolders folders with pictures to scan
	SourceFolders []string `json:"source_folders"`
	//BackupFolder folder to backup
	BackupFolder string `json:"backup_folder"`
	//SearchPageSize the search page query size
	SearchPageSize int `json:"search_page_size"`
	//ThumbX thumbnail width
	ThumbX uint `json:"thumb_x"`
	//ThumbY thumbnail height
	ThumbY uint `json:"thumb_y"`
	//IdleSeconds wait time before application closes.
	IdleSeconds uint `json:"idle_seconds"`
	//DarknetTimeout, in seconds for wait to allow darknet to scan for objects.
	DarknetTimeout uint `json:"darknet_timeout"`
	//DataFolder is the database folder
	DataFolder string `json:"db_folder"`
}

func init() {
	Options.SourceFolders = make([]string, 0)
	Options.SearchPageSize = 12
	Options.ThumbX = 300
	Options.ThumbY = 200
	Options.IdleSeconds = 5
	Options.DarknetTimeout = 70
	appUserFolder, err := appUserFolder()
	if err != nil {
		panic(err)
	}
	Options.DataFolder = appUserFolder
}

//Load loads conf for the current user
func Load() error {
	confFileName, err := getFileFromUserFolder(defaultConfFileName)
	if err != nil {
		return err
	}
	_, err = os.Stat(confFileName)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return nil
	}
	data, err := ioutil.ReadFile(confFileName)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &Options)
	if err != nil {
		return err
	}
	if Options.DataFolder == "" {
		Options.DataFolder, err = appUserFolder()
		if err != nil {
			return err
		}
	}
	filesPath, err := FilesPath()
	if err != nil {
		return err
	}
	return os.MkdirAll(filesPath, 0755)
}

func Save() error {
	data, err := json.MarshalIndent(&Options, "", "    ")
	if err != nil {
		return err
	}
	confFileName, err := getFileFromUserFolder(defaultConfFileName)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(confFileName, data, 0644)
}

//RemoveSourceFolder removes a folder from the source folders
func RemoveSourceFolder(removeFolder string) {
	list := make([]string, 0)
	for _, folder := range Options.SourceFolders {
		if folder != removeFolder {
			list = append(list, folder)
		}
	}
	Options.SourceFolders = list
}

func AddSourceFolder(path string) {
	for _, folder := range Options.SourceFolders {
		if folder == path {
			return
		}
	}
	Options.SourceFolders = append(Options.SourceFolders, path)
}

//appUserFolder returns the app folder under the user account ~/danny/.pictures
func appUserFolder() (string, error) {
	cuser, err := user.Current()
	if err != nil {
		return "", err
	}
	appFolder := filepath.Join(cuser.HomeDir, ".pictures")
	err = os.MkdirAll(appFolder, 0755)
	if err != nil {
		return "", err
	}
	return appFolder, nil
}

func getFileFromUserFolder(fileName string) (string, error) {
	appFolder, err := appUserFolder()
	if err != nil {
		return "", err
	}
	userFileName := filepath.Join(appFolder, fileName)
	return userFileName, nil
}

func getFileFromDataFolder(fileName string) (string, error) {
	err := os.MkdirAll(Options.DataFolder, 0755)
	if err != nil {
		return "", err
	}
	return filepath.Join(Options.DataFolder, fileName), nil
}

//BlevePath returns bleve path
func BlevePath() (string, error) {
	return getFileFromDataFolder(defaultBleveFolder)
}

//BoltPath bold db file path
func BoltPath() (string, error) {
	return getFileFromDataFolder(defaultBoltFileName)
}

//SqlitePath bold db file path
func SqlitePath() (string, error) {
	return getFileFromDataFolder(defaultSQLiteFileName)
}

//FilesPath is the place where files are stored.
func FilesPath() (string, error) {
	return getFileFromDataFolder(defaultFilesPath)
}
