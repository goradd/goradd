package button

import (
	"context"
	"encoding/json"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/log"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/event"
	"github.com/goradd/html5tag"
	"mime/multipart"
	"path"
	"strings"
)

const (
	fileSelectChangeAction = 2000 // This action is just here to update form values
)

const uploadEventName = "upload"

// UploadEvent triggers when an upload has been completed.
// Put an action on this to read the FileHeaders and then copy the uploaded files.
func UploadEvent() *event.Event {
	return event.NewEvent(uploadEventName)
}

// FileInfo is the information regarding each file selected.
// Call FileSelect.FileInfo to retrieve the information.
type FileInfo struct {
	Name        string `json:"name"`
	LasModified int    `json:"lastModified"`
	Size        int    `json:"size"`
	Type        string `json:"type"`
}

type FileSelectI interface {
	page.ControlI
	SetUploadOnChange(bool)
	SetMultiple(bool)
	SetAccept(string)
}

// FileSelect is a standard html file select button.
//
// By default, it will only accept single file. Call SetMultiple to allow multiple selections, and SetAccept
// to control what types of files can be selected based on the file suffixes.
//
// To get the list of files that have been selected, call FileInfo.
// However, to get to the files that are being uploaded, you must do one of the following:
//   - create an Do action on the UploadEvent event, and then call the Upload function, or
//   - create a Server action on a submit button.
//
// In your DoAction handler, call FileHeaders on the FileSelect to retrieve the file headers.
// You then must immediately call Open on each FileHeader to retrieve the information in each file that was uploaded,
// since Go will delete the files after the event is finished processing.
//
// If you would like to trigger an UploadEvent as soon as files are selected, call SetUploadOnChange.
//
// If your application is running behind a web server like apache or nginx, you probably will
// need to set the maximum upload file size in that software.
type FileSelect struct {
	page.ControlBase
	// path taken from the form value. Note that sometimes this has a fake path plus a file name, so likely just the last item is valuable.
	fileInfos []FileInfo

	// uploadOnChange controls whether a change event will trigger an upload
	uploadEventID event.EventID

	// fileHeaders are the fileHeaders right after a file upload. These cannot be
	// serialized, and Go will delete the files after the request, so they have to be dealt with right away.
	fileHeaders []*multipart.FileHeader
}

// NewFileSelect creates a new file select button.
func NewFileSelect(parent page.ControlI, id string) *FileSelect {
	c := &FileSelect{}
	c.Init(c, parent, id)
	return c
}

// Init is called by subclasses of Button to initialize the button control structure.
func (b *FileSelect) Init(self any, parent page.ControlI, id string) {
	b.ControlBase.Init(self, parent, id)
	b.Tag = "input"
	b.SetAttribute("type", "file")
	b.On(event.Change().Private(), action.Do(b.ID(), fileSelectChangeAction))
	b.ParentForm().AddJavaScriptFile(path.Join(config.AssetPrefix, "goradd", "js", "file_select.js"), false, nil)
}

// SetMultiple controls whether the button allows multiple file selections.
func (b *FileSelect) SetMultiple(m bool) {
	b.SetAttribute("multiple", m)
}

// SetAccept is a comma separated list of file endings to filter what is visible in the
// file selection dialog. i.e. ".jqg, .jpeg"
func (b *FileSelect) SetAccept(a string) {
	b.SetAttribute("accept", a)
}

// SetUploadOnChange will set whether the onchange event will immediately trigger the upload process.
// For the upload to happen, you must create an action on the UploadEvent.
func (b *FileSelect) SetUploadOnChange(up bool) {
	if up && b.uploadEventID == event.EventID(0) {
		e := event.Change().Private()
		b.uploadEventID = event.ID(e)
		b.On(e, action.WidgetFunction(b.ID(), "upload"))
	} else if !up && b.uploadEventID != event.EventID(0) {
		b.PrivateOff(b.uploadEventID)
		b.uploadEventID = event.EventID(0)
	}
}

// DrawingAttributes is called by the framework to retrieve the tag's private attributes at draw time.
func (b *FileSelect) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := b.ControlBase.DrawingAttributes(ctx)
	a.SetData("grctl", "fileselect")
	a.SetData("grWidget", "goradd.FileSelect")

	a.Set("name", b.ID()) // needed for posts
	if b.IsRequired() {
		a.Set("required", "")
	}
	return a
}

// UpdateFormValues is used by the framework to cause the control to retrieve its values from the form
func (b *FileSelect) UpdateFormValues(ctx context.Context) {
	id := b.ID()
	grctx := page.GetContext(ctx)

	if v, ok := grctx.FormValue(id); ok {
		d := json.NewDecoder(strings.NewReader(v))
		err := d.Decode(&b.fileInfos)
		if err != nil {
			log.Error(err)
		}
	}
	if grctx.Request.MultipartForm != nil {
		if v, ok := grctx.Request.MultipartForm.File[id]; ok {
			b.fileHeaders = v
		}
	}
}

// Serialize is called by the framework during pagestate serialization.
func (b *FileSelect) Serialize(e page.Encoder) {
	b.ControlBase.Serialize(e)

	if err := e.Encode(b.fileInfos); err != nil {
		panic(err)
	}
	if err := e.Encode(b.uploadEventID); err != nil {
		panic(err)
	}
}

// Deserialize is called by the framework during page state serialization.
func (b *FileSelect) Deserialize(d page.Decoder) {
	b.ControlBase.Deserialize(d)

	if err := d.Decode(&b.fileInfos); err != nil {
		panic(err)
	}
	if err := d.Decode(&b.uploadEventID); err != nil {
		panic(err)
	}
}

// Upload will trigger the upload process.
// Create an action on the UploadEvent() event to respond once the files are uploaded.
func (b *FileSelect) Upload() {
	b.ExecuteWidgetFunction("upload")
}

// FileHeaders returns the file headers for the most recently submitted files.
// The files pointed to by the file headers are only valid when responding to the
// UploadEvent() event. After that, the files they point to will be deleted.
func (b *FileSelect) FileHeaders() []*multipart.FileHeader {
	return b.fileHeaders
}

// FileInfos returns the file information about each file selected.
// FileInfo information is not available for Server events.
func (b *FileSelect) FileInfos() []FileInfo {
	return b.fileInfos
}

// FileSelectCreator is the initialization structure for declarative creation of file selection buttons
type FileSelectCreator struct {
	// ID is the control id
	ID string
	// Multiple controls whether multiple files can be selected
	Multiple bool
	// Accept is a list of command separated file endings that will filter
	// the visible files. i.e. ".jpg, .jpeg"
	Accept string
	// OnChange is an action to take when a file has been selected.
	OnChange action.ActionI
	// UploadOnChange indicates that the upload process should proceed immediately when a file is selected.
	UploadOnChange bool
	// OnUpload is an action to take after the Upload() command has completed. Respond to this action by calling
	// FileHeaders() to get the FileHeaders of the uploaded files and then copying those files to a more permanent place.
	OnUpload action.ActionI

	// ControlOptions are additional options that are common to all controls.
	page.ControlOptions
}

// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c FileSelectCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewFileSelect(parent, c.ID)

	c.Init(ctx, ctrl)
	return ctrl
}

// Init is called by implementations of controls to initialize a control with the
// creator.
func (c FileSelectCreator) Init(ctx context.Context, ctrl FileSelectI) {
	if c.OnChange != nil {
		ctrl.On(event.Change(), c.OnChange)
	}
	if c.OnUpload != nil {
		ctrl.On(UploadEvent(), c.OnUpload)
	}
	if c.Multiple {
		ctrl.SetMultiple(true)
	}
	if c.Accept != "" {
		ctrl.SetAccept(c.Accept)
	}
	if c.UploadOnChange {
		ctrl.SetUploadOnChange(true)
	}
	ctrl.ApplyOptions(ctx, c.ControlOptions)
}

// GetFileSelect is a convenience method to return the button with the given id from the page.
func GetFileSelect(c page.ControlI, id string) *FileSelect {
	return c.Page().GetControl(id).(*FileSelect)
}

func init() {
	page.RegisterControl(&FileSelect{})
}
