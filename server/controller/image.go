package controller

import (
	"github.com/dtylman/pictures/server/view"
	"net/http"
)

// Image serves the active image
func ActiveImage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, mySearch.ActiveImage.Path)
}

func PrevImage(w http.ResponseWriter, r *http.Request) {
	mySearch.PrevImage()
	v := view.New(r)
	v.Vars["image"] = mySearch.ActiveImage
	v.Name = "image/view"
	v.Render(w)
}

func NextImage(w http.ResponseWriter, r *http.Request) {
	mySearch.NextImage()
	v := view.New(r)
	v.Vars["image"] = mySearch.ActiveImage
	v.Name = "image/view"
	v.Render(w)
}

// Thumb servers image thumbnail file by ID
func Thumb(w http.ResponseWriter, r *http.Request) {
	hit, err := getRouterParamInt(r, "hit")
	if err != nil {
		Error500(w, r, err)
		return
	}
	http.ServeFile(w, r, mySearch.Thumbs[hit].Path)
}

func ImageView(w http.ResponseWriter, r *http.Request) {
	hit, err := getRouterParamInt(r, "hit")
	if err != nil {
		Error500(w, r, err)
		return
	}
	mySearch.SetActiveImage(hit)
	v := view.New(r)
	v.Vars["image"] = mySearch.ActiveImage
	v.Name = "image/view"
	v.Render(w)
}
