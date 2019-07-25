package data

import (
	"context"
	"github.com/goradd/goradd/pkg/log"
	"github.com/goradd/goradd/pkg/page"
	"reflect"
)

type DataBinder interface {
	// A DataBinder must be a control so that we can serialize it
	ID() string
	// BindData is called by the data manager to get the data for the control during draw time
	BindData(ctx context.Context, s DataManagerI)
}

// A DataManagerI is the interface for the owner (the embedder) of the DataManager
type DataManagerI interface {
	page.ControlI
	SetDataProvider(b DataBinder)
	// SetData should be passed a slice of data items
	SetData(interface{})
	LoadData(ctx context.Context, owner DataManagerI)
	ResetData()
}

// DataManager is an object designed to be embedded in a control that will help manage the data binding process.
type DataManager struct {
	dataProvider DataBinder

	// data is a temporary copy of the drawing data that is intended to only be loaded during drawing, and then unloaded after drawing.
	data         interface{}
}

func (d *DataManager) SetDataProvider(b DataBinder) {
	d.dataProvider = b
}

func (d *DataManager) HasDataProvider() bool {
	return d.dataProvider != nil
}

// Call SetData to set the data of a control that uses a data binder. You MUST call it with a slice
// of some kind of data.
func (d *DataManager) SetData(data interface{}) {
	kind := reflect.TypeOf(data).Kind()
	if kind != reflect.Slice {
		panic("you must call SetData with a slice")
	}
	d.data = data
}

// ResetData is called by controls that use a data binder to unload the data after it is used.
func (d *DataManager) ResetData() {
	if d.dataProvider != nil {
		d.data = nil
	}
}

// LoadData tells the data binder to load data by calling SetData on the given object. The object should be
// the embedder of the DataManager
func (d *DataManager) LoadData(ctx context.Context, owner DataManagerI) {
	if d.dataProvider != nil && d.data == nil {
		log.FrameworkDebug("Calling BindData")
		d.dataProvider.BindData(ctx, owner) // tell the data binder to call SetData on the given object, or load data some other way
	}
}

// RangeData will call the given function for each item in the data.
// The function should return true to continue, and false to end early.
func (d *DataManager) RangeData(f func(int, interface{}) bool) {
	if d.data == nil {
		return
	}
	listValue := reflect.ValueOf(d.data)
	for i := 0; i < listValue.Len(); i++ {
		itemI := listValue.Index(i).Interface()
		result := f(i, itemI)
		if !result {
			break
		}
	}
}

func (d *DataManager) HasData() bool {
	return d.data != nil
}
