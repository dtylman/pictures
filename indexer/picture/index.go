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

	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	Video = "video"
	Image = "image"
)

type Index struct {
	MD5      string    `json:"md5"`
	MimeType string    `json:"mime_type"`
	Path     string    `json:"path"`
	FileTime time.Time `json:"file_time"`
	Taken    time.Time `json:"taken"`
	Exif     string    `json:"exif"`
	Lat      float64   `json:"lat"`
	Long     float64   `json:"long"`
	Location string    `json:"location"`
	Album    string    `json:"album"`
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
	pic := &Index{}
	pic.Path = path
	folder, _ := filepath.Split(pic.Path)
	pic.Album = filepath.Base(folder)
	pic.FileTime = info.ModTime()
	pic.MimeType = mimemagic.Match("", sig)
	mimeType := strings.Split(pic.MimeType, "/")[0]

	if mimeType != Image && mimeType != Video {
		return nil, errors.New(fmt.Sprintf("File '%s' is '%s' and not '%s' or '%s'", path, mimeType, Image, Video))
	}
	err = pic.populateMD5(file)
	if err != nil {
		return nil, err
	}
	_ = pic.populateExif(file)
	//todo: add error
	return pic, nil
}

func (i *Index) Walk(name exif.FieldName, tag *tiff.Tag) error {
	i.Exif += fmt.Sprintf("%s: %s ", name, tag.String())
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
	i.Exif = ""
	err = x.Walk(i)
	if err != nil {
		//todo: add error
		log.Println(err)
	}
	i.Taken, err = x.DateTime()
	if err != nil {
		//todo: add error
		log.Println(err)
	}
	i.Lat, i.Long, err = x.LatLong()
	if err != nil {
		//todo: add error
		log.Println(err)
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
