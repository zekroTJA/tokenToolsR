package spa

// The most of the code is derived from the gorilla/mux
// readme example:
// https://github.com/gorilla/mux#serving-single-page-applications

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

var FileRx = regexp.MustCompile(`^.*\.(ico|css|js|svg|gif|jpe?g|png)$`)

type SPA struct {
	staticPath string
	index      string

	fileServer http.Handler
}

func NewSPA(staticPath, index string) *SPA {
	return &SPA{
		staticPath: staticPath,
		index:      index,
		fileServer: http.FileServer(http.Dir(staticPath)),
	}
}

func (h *SPA) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(h.staticPath, r.URL.Path)

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.index))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.fileServer.ServeHTTP(w, r)
}
