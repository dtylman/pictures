package controller

import (
	"github.com/dtylman/pictures/db"
	"github.com/dtylman/pictures/server/view"
	"github.com/nfnt/resize"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"path/filepath"
)

// Image servers image file by ID
func Image(w http.ResponseWriter, r *http.Request) {
	imageID := getParamByName(r, "id")
	imageInfo, err := db.GetImage(imageID)
	if err != nil {
		Error500(w, r, err)
	} else {
		http.ServeFile(w, r, imageInfo.Path)
	}
}

// Thumb servers image thumbnail file by ID
func Thumb(w http.ResponseWriter, r *http.Request) {
	imageID := getParamByName(r, "id")
	imageInfo, err := db.GetImage(imageID)
	if err != nil {
		Error500(w, r, err)
		return
	}
	reader, err := os.Open(imageInfo.Path)
	if err != nil {
		Error500(w, r, err)
		return
	}
	i, _, err := image.Decode(reader)
	if err != nil {
		Error500(w, r, err)
		return
	}
	thumb := resize.Thumbnail(300, 200, i, resize.NearestNeighbor)
	err = jpeg.Encode(w, thumb, nil)
	if err != nil {
		Error500(w, r, err)
	}
}

func ImageView(w http.ResponseWriter, r *http.Request) {
	imageID := getParamByName(r, "id")
	doc, err := db.GetImageDocument(imageID)
	if err != nil {
		flash(r, view.FlashError, err.Error())
	}

	// Display the view
	v := view.New(r)

	details := make(map[string]string)
	for _, field := range doc.Fields {
		value := string(field.Value())
		details[field.Name()] = value
		if field.Name() == "path" {
			v.Vars["image"] = filepath.Base(value)
		}
	}
	v.Vars["id"] = imageID
	v.Vars["details"] = details
	v.Name = "image/view"
	v.Render(w)
}
