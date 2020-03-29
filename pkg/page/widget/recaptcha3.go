package widget

import (
	"bytes"
	"context"
	"fmt"
	"github.com/goradd/goradd/pkg/page"
)

type Recaptcha3 struct {
	page.ControlBase
	SiteKey string
	SecretKey string
}

// NewRecaptcha3 creates a new recaptcha widget
func NewRecaptcha3(parent page.ControlI, id string) *Recaptcha3 {
	b := new(Recaptcha3)
	b.Self = b
	b.Init(parent, id)
	return b
}

// Init is called by subclasses of Button to initialize the button control structure.
func (r *Recaptcha3) Init(parent page.ControlI, id string) {
	r.ControlBase.Init(parent, id)
	r.SetShouldAutoRender(true)

}

func (r *Recaptcha3) Draw(ctx context.Context, buf *bytes.Buffer) (err error) {
	s := fmt.Sprintf(`<script src="https://www.google.com/recaptcha/api.js?render=%s"></script>
<script>
grecaptcha.ready(function() {
	grecaptcha.execute('%[1]s', {action: 'homepage'}).then(function(token) {
		goradd.setControlValue("%s", "token", token);
	});
});
</script>`, r.SiteKey, r.ID())
	buf.WriteString(s)
	return
}
