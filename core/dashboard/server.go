// Package dashboard provides a web portal to see the status of the system.
package dashboard

import (
	"github.com/BruteForceFencer/core/config"
	"github.com/BruteForceFencer/core/hitcounter"
	"net/http"
)

// Server is a server that presents the dashboard.
type Server struct {
	http.Server
	mux     *http.ServeMux
	conf    *config.Configuration
	counter *hitcounter.HitCounter
}

// New returns an initialized instance of *Server.
func New(conf *config.Configuration, counter *hitcounter.HitCounter) *Server {
	result := new(Server)
	result.conf = conf
	result.counter = counter
	result.mux = http.NewServeMux()
	result.Server = http.Server{
		Addr:    conf.DashboardAddress,
		Handler: result.mux,
	}

	result.setupRoutes()

	return result
}

// ListenAndServe blocks and listens for requests.
func (s *Server) ListenAndServe() {
	s.Server.ListenAndServe()
}

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/history", s.HandleHistory)
	s.mux.HandleFunc("/", s.HandleAssets)
}
