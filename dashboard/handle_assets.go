package dashboard

import (
	"github.com/BruteForceFencer/bff/hitcounter"
	"github.com/BruteForceFencer/bff/version"
	"html/template"
	"net/http"
	"path/filepath"
)

// HandleAssets serves the HTML, CSS and JS assets for the dashboard.
func (s *Server) HandleAssets(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		s.serveHomePage(w, r)
	} else {
		assetPath := filepath.Join("assets", r.URL.Path)
		http.ServeFile(w, r, assetPath)
	}
}

func (s *Server) serveHomePage(w http.ResponseWriter, r *http.Request) {
	data := struct {
		ListenAddress string
		ListenType    string
		Version       string
		Directions    map[string]*hitcounter.Direction
	}{
		ListenAddress: s.conf.ListenAddress,
		ListenType:    s.conf.ListenType,
		Version:       version.Version,
		Directions:    s.counter.Directions,
	}

	t, err := template.ParseFiles(filepath.FromSlash("assets/dashboard.html"))
	if err != nil {
		http.Error(w, "Unable to find server files.", http.StatusInternalServerError)
		return
	}

	t.Execute(w, data)
}
