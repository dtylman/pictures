package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

const (
	defaultConfFileName = "conf"
	defaultBleveFolder  = "pictures.db"
	defaultBoltFileName = "bolt.db"
)

var Options struct {
	MapQuestAPIKey string   `json:"map_quest_api_key"`
	SourceFolders  []string `json:"source_folders"`
	BackupFolder   string   `json:"backup_folder"`
}

func init() {
	Options.SourceFolders = make([]string, 0)
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
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	appFolder := filepath.Join(user.HomeDir, ".pictures")
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

//BleveFolder returns bleve path
func BleveFolder() (string, error) {
	return getPathForFile(defaultBleveFolder)
}

//BoltPath bold db file path
func BoltPath() (string, error) {
	return getPathForFile(defaultBoltFileName)
}
