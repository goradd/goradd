package welcome

import (
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/sys"
	"github.com/goradd/goradd/web/app"
)

func init() {
	if !config.Release {
		app.RegisterStaticPath("/goradd", sys.SourceDirectory())
	}
}
