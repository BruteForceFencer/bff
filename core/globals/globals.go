package globals

import (
	"github.com/BruteForceFencer/core/config"
	"github.com/BruteForceFencer/core/controlserver"
	"github.com/BruteForceFencer/core/dashboard"
	"github.com/BruteForceFencer/core/hitcounter"
)

var (
	Configuration *config.Configuration
	Dashboard     *dashboard.Server
	HitCounter    *hitcounter.HitCounter
	Server        *controlserver.Server
)
