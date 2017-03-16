package lockfile

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

//Info lockfile information struct
type Info struct {
	Address string `json:"address"`
	PID     int    `json:"pid"`
}

var created bool

const name = "pictures.lock"

//Open opens the lock file
func Open() (*Info, error) {
	data, err := ioutil.ReadFile(fileName())
	if err != nil {
		return nil, err
	}
	var info Info
	err = json.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

//Create creates the lock file
func (i *Info) Create() error {
	data, err := json.Marshal(i)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fileName(), data, 0755)
	if err != nil {
		return err
	}
	created = true
	return nil
}

//URL format the address as URL
func (i *Info) URL() string {
	return fmt.Sprintf("http://%s", i.Address)
}

//Delete removes the lock file
func Delete() error {
	if created {
		return os.Remove(fileName())
	}
	return nil
}

func fileName() string {
	return filepath.Join(os.TempDir(), name)
}
