package data

import (
	"context"
	"github.com/spekary/goradd/page"
	"goradd-project/config"
	"reflect"
	"github.com/spekary/goradd/log"
)

type DataBinder interface {
	BindData(ctx context.Context, s DataManagerI)
}

// A DataManagerI is the interface for the owner (the embedder) of the DataManager
type DataManagerI interface {
	page.ControlI
	SetDataProvider(b DataBinder)
	// SetData should be passed a slice of data items
	SetData(...interface{})
	GetData(ctx context.Context, owner DataManagerI)
	ResetData()
}

// DataManager is an object designed to be embedded in a control that will help manage the data binding process.
type DataManager struct {
	dataProvider DataBinder
	Data         []interface{}
}

func (d *DataManager) SetDataProvider(b DataBinder) {
	d.dataProvider = b
}

func (d *DataManager) HasDataProvider() bool {
	return d.dataProvider != nil
}

// Call SetData to set the data of a control that uses a data binder. Generally, you should call it with an expanded slice,
// and in fact, it will issue a warning if you give it one item that is a slice, because it will assume that you accidentally
// did not expand the array.
func (d *DataManager) SetData(data ...interface{}) {
// We use an expanded interface list here, instead of a slice of interfaces, because there is a subtle difference between the two.
// A slice of interfaces will require just that, a slice of interfaces and nothing else. However, this declaration above lets you
// send in a slice of whatever kind of object you want, and also just list out the objects if needed.

	if (!config.Release) { // This code will only be included in the debug version
		if len(data) == 1 &&
			reflect.TypeOf(data[0]).Kind() == reflect.Slice {

			// Its possible this is legitimate, but not likely, so we issue a warning
			log.Warning("You called SetData with a single entry that is a slice. Did you not expand the slice?")
		}
	}
	d.Data = data
}

func (d *DataManager) ResetData() {
	if d.dataProvider != nil {
		d.Data = nil
	}
}

// GetData tells the data binder to load data by calling SetData on the given object. The object should be
// the embedder of the DataManager
func (d *DataManager) GetData(ctx context.Context, owner DataManagerI) {
	if d.dataProvider != nil && d.Data == nil {
		d.dataProvider.BindData(ctx, owner) // tell the data binder to call SetData on the given object, or load data some other way
	}
}
