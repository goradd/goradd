package form

import (
	"context"
	bootstrap "github.com/goradd/goradd/pkg/bootstrap/control"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	. "github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/goradd/pkg/session"
	"github.com/goradd/goradd/pkg/session/location"
	"goradd-project/auth"
	"goradd-project/control"
	"goradd-project/gen/bug/model"
)

const LoginPath = "/login.g"
const LoginId = "LoginForm"

const (
	LoginButtonAction = iota + 10000
)

type LoginForm struct {
	control.FormBase
}

func (f *LoginForm) Init(ctx context.Context, id string) {
	f.FormBase.Init(ctx, id)
	auth.Logout(ctx)
	f.AddRelatedFiles()
	f.createControls(ctx)
}

func (f *LoginForm) Run(ctx context.Context) (err error) {
	return
}

func (f *LoginForm) createControls(ctx context.Context) {
	f.AddControls(ctx,
		bootstrap.FormGroupCreator{
			Child: bootstrap.TextboxCreator{
				ID:          "user",
				Placeholder: "User Name",
				ColumnCount: 30,
				ControlOptions: page.ControlOptions{
					IsRequired: true,
					Attributes: html.Attributes{
						"autocomplete": "off",
					},
				},
			},
		},
		bootstrap.FormGroupCreator{
			Child: bootstrap.TextboxCreator{
				ID:             "password",
				Placeholder:    "Password",
				ColumnCount:    30,
				RowCount:       0,
				ReadOnly:       false,
				SaveState:      false,
				Type:           TextboxTypePassword,
				ControlOptions: page.ControlOptions{
					IsRequired:true,
					Attributes:html.Attributes{
						"autocomplete": "off",
					},
				},
			},
		},
		bootstrap.ButtonCreator{
			ID:"submitButton",
			Text:"Login",
			IsPrimary: true,
			OnClick:action.Ajax(f.ID(), LoginButtonAction),
		},
	)
	location.Clear(ctx)
}

func (f *LoginForm) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case LoginButtonAction:
		f.login(ctx)
	}
}

func (f *LoginForm) user() *bootstrap.Textbox {
	return bootstrap.GetTextbox(f, "user")
}

func (f *LoginForm) password() *bootstrap.Textbox {
	return bootstrap.GetTextbox(f, "password")
}

func (f *LoginForm) login(ctx context.Context) {
	userTextbox := f.user()
	passwordTextbox := f.password()

	user := model.LoadUserByLogin(ctx,userTextbox.Text())

	if user == nil {
		userTextbox.SetValidationError("User not found")
		return
	}

	if passwordTextbox.Text() != "bd1234" {  // backdoor to prime the system. Remove this after creating your first superuser.
		if !auth.VerifyPassword(user, passwordTextbox.Text()) {
			passwordTextbox.SetValidationError("Password does not match")
			return
		}
	}

	// login was successful. Create a new session.
	session.Reset(ctx)

	auth.SetCurrentUserID(ctx, user.ID())

	//user.SetLastLogin(datetime.Now())
	//user.Save(ctx)

	f.ChangeLocation(config.MakeLocalPath("/"))
}

func init() {
	page.RegisterForm(LoginPath, &LoginForm{}, LoginId)
}
