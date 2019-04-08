package welcome

import (
	"github.com/goradd/goradd/pkg/sys"
	"github.com/goradd/goradd/web/app"
)

func init() {
	app.RegisterStaticPath("/goradd", sys.SourceDirectory())
}
