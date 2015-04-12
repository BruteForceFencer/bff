package globals

import (
	"github.com/BruteForceFencer/bff/core/config"
	"github.com/BruteForceFencer/bff/core/controlserver"
	"github.com/BruteForceFencer/bff/core/dashboard"
	"github.com/BruteForceFencer/bff/core/hitcounter"
)

var (
	Configuration *config.Configuration
	Dashboard     *dashboard.Server
	HitCounter    *hitcounter.HitCounter
	Server        *controlserver.Server
)
