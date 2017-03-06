package controller

import (
	"github.com/dtylman/pictures/db"
	"github.com/dtylman/pictures/server/view"
	"net/http"
	"path/filepath"
)

// Image servers image file by ID
func Image(w http.ResponseWriter, r *http.Request) {
	imageID := getParamByName(r, "id")
	path, err := db.PathForImage(imageID)
	if err != nil {
		flash(r, view.FlashError, err.Error())
	}
	http.ServeFile(w, r, path)
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
