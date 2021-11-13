//go:build !dev
// +build !dev

package webapp

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/sirupsen/logrus"
)

//go:embed dist/*
var Webapp embed.FS

func WebHandler(l *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		wa, err := fs.Sub(Webapp, "dist")
		if err != nil {
			l.WithError(err).Error("Failed to load webapp from fs")
		}
		http.FileServer(http.FS(wa)).ServeHTTP(w, r)
	}
}
