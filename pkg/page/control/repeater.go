package control

import (
	"context"
	"github.com/goradd/goradd/pkg/log"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/html5tag"
	"io"
)

type RepeaterI interface {
	PagedControlI
	DrawItem(ctx context.Context, i int, data interface{}, w io.Writer)
	SetItemHtmler(h RepeaterHtmler) RepeaterI
}

type Repeater struct {
	page.ControlBase
	PagedControl
	DataManager
	itemHtmler   RepeaterHtmler
	itemHtmlerId string // only used for serialization
}

type RepeaterHtmler interface {
	RepeaterHtml(ctx context.Context, r RepeaterI, i int, data interface{}, w io.Writer)
}

// NewRepeater creates a new Repeater
func NewRepeater(parent page.ControlI, id string) *Repeater {
	r := &Repeater{}
	r.Self = r
	r.Init(parent, id)
	return r
}

// Init is an internal function that enables the object-oriented pattern of calling virtual functions used by the
// goradd controls.
func (r *Repeater) Init(parent page.ControlI, id string) {
	r.ControlBase.Init(parent, id)
	r.Tag = "div"
}

// this returns the RepeaterI interface for calling into "virtual" functions. This allows us to call functions defined
// by a subclass.
func (r *Repeater) this() RepeaterI {
	return r.Self.(RepeaterI)
}

// SetItemHtmler sets the htmler that provides the html for each item in the repeater.
func (r *Repeater) SetItemHtmler(h RepeaterHtmler) RepeaterI {
	r.itemHtmler = h
	return r.this()
}

// DrawTag is called by the framework to draw the tag. The Repeater overrides this to call into the DataProvider
// to load the table's data into memory just before drawing. The data will be unloaded after drawing.
func (r *Repeater) DrawTag(ctx context.Context, w io.Writer) {
	log.FrameworkDebug("Drawing repeater tag")
	if r.HasDataProvider() {
		log.FrameworkDebug("Getting repeater data")
		r.this().LoadData(ctx, r.this())
		defer r.ResetData()
	}
	r.ControlBase.DrawTag(ctx, w)
}

// DrawingAttributes is an override to add attributes to the table, including not showing the table at all if there
// is no data to show. This will hide header and footer cells and potentially the outline of the table when there is no
// data in the table.
func (r *Repeater) DrawingAttributes(ctx context.Context) html5tag.Attributes {
	a := r.ControlBase.DrawingAttributes(ctx)
	a.SetData("grctl", "repeater")
	return a
}

// DrawInnerHtml is an override to draw the individual items of the repeater.
func (r *Repeater) DrawInnerHtml(ctx context.Context, w io.Writer) {
	var this = r.this() // Get the sub class so we call into its hooks for drawing

	r.RangeData(func(index int, value interface{}) bool {
		this.DrawItem(ctx, index, value, w)
		return true
	})
	return
}

func (r *Repeater) DrawItem(ctx context.Context, i int, data interface{}, w io.Writer) {
	if r.itemHtmler != nil {
		r.itemHtmler.RepeaterHtml(ctx, r, i, data, w)
	}
	return
}

// MarshalState is an internal function to save the state of the control
func (r *Repeater) MarshalState(m page.SavedState) {
	r.PagedControl.MarshalState(m)
}

// UnmarshalState is an internal function to restore the state of the control
func (r *Repeater) UnmarshalState(m page.SavedState) {
	r.PagedControl.UnmarshalState(m)
}

func (r *Repeater) Serialize(e page.Encoder) {
	r.ControlBase.Serialize(e)
	r.PagedControl.Serialize(e)
	r.DataManager.Serialize(e)

	// If itemHtmler is a control, we will just serialize the control's id, since the control will get
	// serialized elsewhere. Otherwise, we serialize the itemHtmler itself.
	var htmler interface{} = r.itemHtmler
	if ctrl, ok := r.itemHtmler.(page.ControlI); ok {
		htmler = ctrl.ID()
	}
	if err := e.Encode(&htmler); err != nil {
		panic(err)
	}
}

func (r *Repeater) Deserialize(dec page.Decoder) {
	r.ControlBase.Deserialize(dec)
	r.PagedControl.Deserialize(dec)
	r.DataManager.Deserialize(dec)

	var htmler interface{}
	if err := dec.Decode(&htmler); err != nil {
		panic(err)
	}
	if id, ok := htmler.(string); ok {
		r.itemHtmlerId = id
	} else {
		r.itemHtmler = htmler.(RepeaterHtmler)
	}
}

func (r *Repeater) Restore() {
	r.ControlBase.Restore()
	if r.itemHtmlerId != "" {
		r.itemHtmler = r.Page().GetControl(r.itemHtmlerId).(RepeaterHtmler)
	}
	return
}

// RepeaterCreator creates a table that can be paged
type RepeaterCreator struct {
	// ID is the control id
	ID string
	// ItemHtmler is the object that provides the html for each item
	ItemHtmler RepeaterHtmler
	// DataProvider is the data binder for the table. It can be either a control id or a DataBinder
	DataProvider DataBinder
	// DataProviderID is the control id of the data binder for the table.
	DataProviderID string
	// Data is the actual data for the table, and should be a slice of objects
	Data interface{}
	page.ControlOptions
	// PageSize is the number of rows to include in a page
	PageSize int
	// SaveState will cause the table to remember what page it was on
	SaveState bool
}

// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c RepeaterCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewRepeater(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

func (c RepeaterCreator) Init(ctx context.Context, ctrl RepeaterI) {
	if c.ItemHtmler != nil {
		ctrl.SetItemHtmler(c.ItemHtmler)
	}
	if c.DataProvider != nil {
		ctrl.SetDataProvider(c.DataProvider)
	} else if c.DataProviderID != "" {
		provider := ctrl.Page().GetControl(c.DataProviderID).(DataBinder)
		ctrl.SetDataProvider(provider)
	}

	if c.Data != nil {
		ctrl.SetData(c.Data)
	}
	ctrl.ApplyOptions(ctx, c.ControlOptions)
	if c.PageSize != 0 {
		ctrl.SetPageSize(c.PageSize)
	}
	if c.SaveState {
		ctrl.SaveState(ctx, true)
	}
}

// GetRepeater is a convenience method to return the repeater with the given id from the page.
func GetRepeater(c page.ControlI, id string) *Repeater {
	return c.Page().GetControl(id).(*Repeater)
}

func init() {
	page.RegisterControl(&Repeater{})
}
