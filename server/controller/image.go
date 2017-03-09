package controller

import (
	"github.com/blevesearch/bleve"
	"github.com/dtylman/pictures/indexer/db"
	"github.com/dtylman/pictures/server/view"
	"github.com/nfnt/resize"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"path/filepath"
)

// Image servers image file by ID
func Image(w http.ResponseWriter, r *http.Request) {
	imageID := getRouterParam(r, "id")
	imageInfo, err := db.GetImage(imageID)
	if err != nil {
		Error500(w, r, err)
	} else {
		http.ServeFile(w, r, imageInfo.Path)
	}
}

// Thumb servers image thumbnail file by ID
func Thumb(w http.ResponseWriter, r *http.Request) {
	imageID := getRouterParam(r, "id")
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
	imageID := getRouterParam(r, "id")
	req := bleve.NewSearchRequest(bleve.NewDocIDQuery([]string{imageID}))
	req.Fields = []string{"*"}
	sr, err := db.Search(req)
	if err != nil {
		flashError(r, err)
	}

	// Display the view
	v := view.New(r)
	v.Vars["id"] = imageID

	v.Vars["details"] = sr.Hits[0].Fields
	v.Vars["image"] = filepath.Base(sr.Hits[0].Fields["path"].(string))
	v.Name = "image/view"
	v.Render(w)
}
