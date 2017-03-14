package model

import (
	"github.com/dtylman/pictures/conf"
	"github.com/nfnt/resize"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"

	"image"
)

type ThumbItem struct {
	Path string
	MD5  string
}

func (s *Search) buildThumbs() {
	s.Thumbs = make([]ThumbItem, s.Result.Hits.Len())
	for i := 0; i < s.Result.Hits.Len(); i++ {
		s.Thumbs[i].init(s.Result.Hits[i].Fields["md5"].(string), s.Result.Hits[i].Fields["path"].(string))
	}
}

func (t *ThumbItem) init(md5 string, path string) {
	folder, err := conf.ThumbPath()
	if err != nil {
		log.Println(err)
		return
	}
	t.MD5 = md5
	t.Path = filepath.Join(folder, md5)
	_, err = os.Stat(t.Path)
	if err == nil {
		// already exists
		return
	}
	reader, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return
	}
	defer reader.Close()
	i, _, err := image.Decode(reader)
	if err != nil {
		log.Println(err)
		return
	}
	thumb := resize.Thumbnail(conf.Options.ThumbX, conf.Options.ThumbY, i, resize.NearestNeighbor)
	writer, err := os.Create(t.Path)
	if err != nil {
		log.Println(err)
	}
	defer writer.Close()
	err = jpeg.Encode(writer, thumb, nil)
	if err != nil {
		log.Println(err)
	}
}
