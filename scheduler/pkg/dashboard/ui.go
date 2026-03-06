package dashboard

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed all:web/scheduler/dist/*
var uiAssets embed.FS

func GetFileSystem() http.FileSystem {
	subFS, err := fs.Sub(uiAssets, "web/scheduler/dist")
	if err != nil {
		panic(err)
	}
	return http.FS(subFS)
}
