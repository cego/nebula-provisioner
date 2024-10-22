//go:build !dev
// +build !dev

package webapp

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
)

//go:embed dist/browser/*
var Webapp embed.FS

func WebHandler(l *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	wa, err := fs.Sub(Webapp, "dist/browser")
	if err != nil {
		l.WithError(err).Error("Failed to load webapp from fs")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Support frontend routing
		if r.URL.Path != "/" {
			_, err := fs.Stat(wa, strings.TrimPrefix(path.Clean(r.URL.Path), "/"))
			if err != nil {
				if os.IsNotExist(err) {
					r.URL.Path = "/"
				}
			}
		}
		http.FileServer(http.FS(wa)).ServeHTTP(w, r)
	}
}
