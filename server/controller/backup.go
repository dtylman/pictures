package controller

import (
	"fmt"
	"github.com/dtylman/pictures/backuper"
	"github.com/dtylman/pictures/server/view"
	"net/http"
)

// Backup displays the backup page
func Backup(w http.ResponseWriter, r *http.Request) {
	// Display the view
	v := view.New(r)
	v.Name = "backup/backup"
	v.Render(w)
}

// Backup displays the backup page
func BackupStatus(w http.ResponseWriter, r *http.Request) {
	obj, _ := backuper.Status.ToJSON()
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, obj)
}
