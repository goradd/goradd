package control

import (
	"context"
	"github.com/spekary/goradd/page"
	localpage "goradd/page"
	"goradd/config"
	"github.com/spekary/goradd/orm/db"
	"github.com/spekary/goradd/page/event"
	"github.com/spekary/goradd/page/action"
	"strings"
	"fmt"
)

const (
	databaseProfileAction = iota + 10000
)

type FormI interface {
	localpage.FormBaseI
}


// The Form is the control that all GetForm objects should descend from, and is the master container for all other goradd controls.
type Form struct {
	localpage.FormBase
}

// The methods below are here to prevent import cycles.

func (f *Form) Init(ctx context.Context, self page.FormBaseI, path string, id string) {
	f.FormBase.Init(ctx, self, path, id)

	if config.Mode == config.AppModeDevelopment && db.IsProfiling(ctx) {
		btn := NewButton(f, "grProfileButton")
		btn.SetText("SQL Profile <i class='fas fa-arrow-circle-down' ></i>")
		btn.SetEscapeText(false)
		btn.On(event.Click(), action.Ajax(f.ID(), databaseProfileAction))
		btn.SetShouldAutoRender(true)

		panel := NewPanel(f, "grProfilePanel")
		panel.SetShouldAutoRender(true)
		panel.SetEscapeText(false)
		panel.SetVisible(false)
	}

	/*	TODO: Add a dialog and designer click if in design mode
					if (defined('QCUBED_DESIGN_MODE') && QCUBED_DESIGN_MODE == 1) {
					// Attach custom event to dialog to handle right click menu items sent by form

					$dlg = new Q\ModelConnector\EditDlg ($objClass, 'qconnectoreditdlg');

					$dlg->addAction(
						new Q\Event\On('qdesignerclick'),
						new Q\Action\Ajax ('ctlDesigner_Click', null, null, 'ui')
					);
				}

	*/



}

func (f *Form) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case databaseProfileAction:
		if c := f.Page().GetControl("grProfilePanel"); c != nil{
			if c.IsVisible() {
				c.SetVisible(false)
			} else {
				c.SetVisible(true)
				var s string
				if profiles := db.GetProfiles(ctx); profiles != nil {
					for _, profile := range profiles {
						dif := profile.EndTime.Sub(profile.BeginTime)
						sql := strings.Replace(profile.Sql, "\n", "<br />", -1)
						s += fmt.Sprintf(`<p class="profile"><div>Time: %s Begin: %s End: %s</div><div>%s</div></p>`,
							dif.String(), profile.BeginTime.Format("3:04:05.000"), profile.EndTime.Format("3:04:05.000"), sql)
					}
				}
				c.SetText(s)
			}
		}
	}
}
