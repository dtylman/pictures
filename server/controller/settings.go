package controller

import (
	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/server/view"
	"net/http"
)

// Settings controls the settings tab
func Settings(w http.ResponseWriter, r *http.Request) {
	// Display the view
	v := view.New(r)
	v.Name = "settings/settings"
	if r.Method == http.MethodPost {
		conf.Options.BackupFolder = r.FormValue("backup_folder")
		conf.Options.MapQuestAPIKey = r.FormValue("quest_key")
		sourceFolder := r.FormValue("source_folder")
		if sourceFolder != "" {
			conf.Options.SourceFolders = append(conf.Options.SourceFolders, sourceFolder)

		}
		err := conf.Save()
		if err != nil {
			panic(err)
		}
	}
	v.Vars["backup_folder"] = conf.Options.BackupFolder
	v.Vars["quest_key"] = conf.Options.MapQuestAPIKey
	v.Vars["source_folders"] = conf.Options.SourceFolders
	v.Render(w)
}

func RemoveSourceFolder(w http.ResponseWriter, r *http.Request) {
	folder := getRouterParam(r, "folder")
	if folder != "" {
		conf.RemoveSourceFolder(folder)
	}
	http.Redirect(w, r, "/settings", http.StatusFound)
}
