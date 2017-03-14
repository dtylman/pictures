package thumbs

import (
	"github.com/dtylman/pictures/conf"
	"github.com/nfnt/resize"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
)

//MakeThumb creates a thumbnail from a picture and returns a path to the thumbnail
func MakeThumb(srcPath string, md5 string, overwrite bool) (string, error) {
	folder, err := conf.ThumbPath()
	if err != nil {
		return srcPath, err
	}
	thumbFile := filepath.Join(folder, md5)
	_, err = os.Stat(thumbFile)
	if err == nil && !overwrite {
		// already exists
		return thumbFile, nil
	}

	reader, err := os.Open(srcPath)
	if err != nil {
		return srcPath, err
	}
	defer reader.Close()
	i, _, err := image.Decode(reader)
	if err != nil {
		return srcPath, err
	}
	thumb := resize.Thumbnail(conf.Options.ThumbX, conf.Options.ThumbY, i, resize.NearestNeighbor)
	writer, err := os.Create(thumbFile)
	if err != nil {
		return srcPath, err
	}
	defer writer.Close()
	err = jpeg.Encode(writer, thumb, nil)
	if err != nil {
		return srcPath, err
	}
	return thumbFile, nil
}
