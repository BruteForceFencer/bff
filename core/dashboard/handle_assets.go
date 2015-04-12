package dashboard

import (
	"github.com/BruteForceFencer/core/hitcounter"
	"github.com/BruteForceFencer/core/version"
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

// installPath is the path to the BFF installation.
var installPath string

func init() {
	// Arg[0] is assumed to be an absolute path.
	installPath = filepath.Join(filepath.Dir(os.Args[0]), "..")
}

// HandleAssets serves the HTML, CSS and JS assets for the dashboard.
func (s *Server) HandleAssets(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		s.serveHomePage(w, r)
	} else {
		assetPath := path.Join("assets", r.URL.Path)
		assetPath = filepath.Join(installPath, filepath.FromSlash(assetPath))
		http.ServeFile(w, r, assetPath)
	}
}

func (s *Server) serveHomePage(w http.ResponseWriter, r *http.Request) {
	data := struct {
		ListenAddress string
		ListenType    string
		Version       string
		Directions    []hitcounter.Direction
	}{
		ListenAddress: s.conf.ListenAddress,
		ListenType:    s.conf.ListenType,
		Version:       version.Version,
		Directions:    s.conf.Directions,
	}

	htmlPath := filepath.Join(installPath, filepath.FromSlash("assets/dashboard.html"))
	t, err := template.ParseFiles(htmlPath)
	if err != nil {
		http.Error(w, "Unable to find server files.", http.StatusInternalServerError)
		return
	}

	t.Execute(w, data)
}
