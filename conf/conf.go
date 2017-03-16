package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

const (
	defaultConfFileName = "conf"
	defaultBleveFolder  = "pictures.db"
	defaultThumbFolder  = "thumbs"
	defaultBoltFileName = "bolt.db"
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
}

func init() {
	Options.SourceFolders = make([]string, 0)
	Options.SearchPageSize = 12
	Options.ThumbX = 300
	Options.ThumbY = 200
	Options.IdleSeconds = 300
	thumbPath, err := ThumbPath()
	if err != nil {
		log.Fatal(err)
	}
	err = os.MkdirAll(thumbPath, 0755)
	if err != nil {
		log.Fatal(err)
	}
}

//Load loads conf for the current user
func Load() error {
	confFileName, err := getPathForFile(defaultConfFileName)
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
	return json.Unmarshal(data, &Options)
}

func getPathForFile(fileName string) (string, error) {
	cuser, err := user.Current()
	if err != nil {
		return "", err
	}
	appFolder := filepath.Join(cuser.HomeDir, ".pictures")
	err = os.MkdirAll(appFolder, 0755)
	if err != nil {
		return "", err
	}
	confFileName := filepath.Join(appFolder, fileName)
	return confFileName, err
}

func Save() error {
	data, err := json.MarshalIndent(&Options, "", "    ")
	if err != nil {
		return err
	}
	confFileName, err := getPathForFile(defaultConfFileName)
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

//BlevePath returns bleve path
func BlevePath() (string, error) {
	return getPathForFile(defaultBleveFolder)
}

//BoltPath bold db file path
func BoltPath() (string, error) {
	return getPathForFile(defaultBoltFileName)
}

func ThumbPath() (string, error) {
	return getPathForFile(defaultThumbFolder)
}
