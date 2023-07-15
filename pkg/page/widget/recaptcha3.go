package widget

import (
	"context"
	"fmt"
	"github.com/goradd/goradd/pkg/page"
	"io"
)

type Recaptcha3 struct {
	page.ControlBase
	SiteKey   string
	SecretKey string
}

// NewRecaptcha3 creates a new recaptcha widget
func NewRecaptcha3(parent page.ControlI, id string) *Recaptcha3 {
	b := new(Recaptcha3)
	b.Init(b, parent, id)
	return b
}

// Init is called by subclasses of Button to initialize the button control structure.
func (r *Recaptcha3) Init(self any, parent page.ControlI, id string) {
	r.ControlBase.Init(self, parent, id)
	r.SetShouldAutoRender(true)

}

func (r *Recaptcha3) Draw(ctx context.Context, w io.Writer) (err error) {
	s := fmt.Sprintf(`<script src="https://www.google.com/recaptcha/api.js?render=%s"></script>
<script>
grecaptcha.ready(function() {
	grecaptcha.execute('%[1]s', {action: 'homepage'}).then(function(token) {
		goradd.setControlValue("%s", "token", token);
	});
});
</script>`, r.SiteKey, r.ID())
	_, err = io.WriteString(w, s)
	return
}
