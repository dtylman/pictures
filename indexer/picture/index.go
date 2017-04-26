package picture

import (
	"bitbucket.org/taruti/mimemagic"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"io"

	"os"
	"path/filepath"
	"strings"
	"time"
	"github.com/dtylman/pictures/tasklog"
)

const (
	Video = "video"
	Image = "image"
)

type Index struct {
	MD5      string
	MimeType string
	Path     string
	FileTime time.Time
	Taken    time.Time
	Lat      float64
	Long     float64
	Location string `json:"location"`
	Album    string `json:"album"`
	Objects  string `json:"objects"`
	Faces    string `json:"faces"`
	Exif     map[string]string
}

func NewIndex(path string, info os.FileInfo) (*Index, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	sig := make([]byte, 1024)
	_, err = file.Read(sig)
	if err != nil {
		return nil, err
	}
	idx := &Index{}
	idx.Path = path
	folder, _ := filepath.Split(idx.Path)
	idx.Album = filepath.Base(folder)
	idx.FileTime = info.ModTime()
	idx.MimeType = mimemagic.Match("", sig)
	if !MimeIs(idx.MimeType, Image, Video) {
		return nil, errors.New(fmt.Sprintf("File '%s' is '%s' and not '%s' or '%s'", path, idx.MimeType, Image, Video))
	}
	err = idx.populateMD5(file)
	if err != nil {
		return nil, err
	}
	err = idx.populateExif(file)
	if err != nil {
		tasklog.Error(err)
	}
	return idx, nil
}

func (i *Index) Walk(name exif.FieldName, tag *tiff.Tag) error {
	i.Exif[string(name)] = tag.String()
	return nil
}

func (i *Index) populateExif(file *os.File) error {
	_, err := file.Seek(0, 0)
	if err != nil {
		return err
	}
	x, err := exif.Decode(file)
	if err != nil {
		return err
	}
	i.Exif = make(map[string]string)
	err = x.Walk(i)
	if err != nil {
		tasklog.Error(err)
	}
	i.Taken, err = x.DateTime()
	if err != nil {
		tasklog.Error(err)
	}
	i.Lat, i.Long, err = x.LatLong()
	if err != nil {
		tasklog.Error(err)
	}
	return nil
}

func (i *Index) populateMD5(file *os.File) error {
	_, err := file.Seek(0, 0)
	if err != nil {
		return err
	}
	h := md5.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return err
	}
	i.MD5 = hex.EncodeToString(h.Sum(nil))
	return nil
}

func (i*Index) ExifString() string {
	return fmt.Sprintf("%v", i.Exif)
}
//MimeIs return true if mime type is one of the provided array
func MimeIs(mimeType string, pictureType ...string) bool {
	base := strings.Split(mimeType, "/")[0]
	for _, wanted := range pictureType {
		if wanted == base {
			return true
		}
	}
	return false
}
