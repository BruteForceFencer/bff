package globals

import (
	"github.com/BruteForceFencer/bff/config"
	"github.com/BruteForceFencer/bff/controlserver"
	"github.com/BruteForceFencer/bff/dashboard"
	"github.com/BruteForceFencer/bff/hitcounter"
)

var (
	Configuration *config.Configuration
	Dashboard     *dashboard.Server
	HitCounter    *hitcounter.HitCounter
	Server        *controlserver.Server
)
