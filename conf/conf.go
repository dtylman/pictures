package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

const defaultConfFileName = ".pictures"

var Options struct {
	MapQuestAPIKey string   `json:"map_quest_api_key"`
	SourceFolders  []string `json:"source_folders"`
	BackupFolder   string   `json:"backup_folder"`
}

func init() {
	//Options.SourceFolders = make([]string, 0)
	Options.SourceFolders = []string{"one", "two", "three"}
}

//Load loads conf for the current user
func Load() error {
	confFileName, err := getFileName()
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

func getFileName() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	confFileName := filepath.Join(user.HomeDir, defaultConfFileName)
	return confFileName, err
}

func Save() error {
	data, err := json.Marshal(&Options)
	if err != nil {
		return err
	}
	confFileName, err := getFileName()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(confFileName, data, 0755)
}
