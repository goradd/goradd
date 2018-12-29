package data

import (
	"context"
	"github.com/goradd/goradd/pkg/page"
	"reflect"
)

type DataBinder interface {
	BindData(ctx context.Context, s DataManagerI)
}

// A DataManagerI is the interface for the owner (the embedder) of the DataManager
type DataManagerI interface {
	page.ControlI
	SetDataProvider(b DataBinder)
	// SetData should be passed a slice of data items
	SetData(interface{})
	GetData(ctx context.Context, owner DataManagerI)
	ResetData()
}

// DataManager is an object designed to be embedded in a control that will help manage the data binding process.
type DataManager struct {
	dataProvider DataBinder
	Data         interface{}
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

// RangeData will call the given function for each item in the data.
// The function should return true to continue, and false to end early.
func (d *DataManager) RangeData(f func( int, interface{}) bool) {
	listValue := reflect.ValueOf(d.Data)
	for i := 0;  i < listValue.Len(); i++ {
		itemI := listValue.Index(i).Interface()
		result := f(i, itemI)
		if !result {
			break
		}
	}
}

