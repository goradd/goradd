package control

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
	DataManagerEmbedder
}

// DataManagerEmbedder is the interface to include in embedded control interfaces
// Currently go does not allow interface conflicts, but that is scheduled to change
type DataManagerEmbedder interface {
	SetDataProvider(b DataBinder)
	HasDataProvider() bool
	// SetData should be passed a slice of data items
	SetData(data interface{})
	SetDataWithOffset(data interface{}, offset int)
	LoadData(ctx context.Context, owner DataManagerI)
	ResetData()
}


// DataManager is an object designed to be embedded in a control that will help manage the data binding process.
type DataManager struct {
	dataProviderID string

	// data is a temporary copy of the drawing data that is intended to only be loaded during drawing, and then unloaded after drawing.
	data         interface{}
	// dataOffset is the first row number represented by the data
	dataOffset   int
}

func (d *DataManager) SetDataProvider(b DataBinder) {
	d.dataProviderID = b.ID()
}

func (d *DataManager) HasDataProvider() bool {
	return d.dataProviderID != ""
}

// Call SetData to set the data of a control that uses a data binder. You MUST call it with a slice
// of some kind of data.
func (d *DataManager) SetData(data interface{}) {
	kind := reflect.TypeOf(data).Kind()
	if kind != reflect.Slice {
		panic("you must call SetData with a slice")
	}
	d.data = data
	d.dataOffset = 0
}

// Use SetDataWithOffset with controls that show a window onto a bigger data set.
//
// offset is the first row number that is represented by the data.
func (d *DataManager) SetDataWithOffset(data interface{}, offset int) {
	kind := reflect.TypeOf(data).Kind()
	if kind != reflect.Slice {
		panic("you must call SetData with a slice")
	}
	d.data = data
	d.dataOffset = offset
}


// ResetData is called by controls that use a data binder to unload the data after it is used.
func (d *DataManager) ResetData() {
	if d.HasDataProvider() {
		d.data = nil
	}
}

// LoadData tells the data binder to load data by calling SetData on the given object. The object should be
// the embedder of the DataManager
func (d *DataManager) LoadData(ctx context.Context, owner DataManagerI) {
	if d.HasDataProvider() && // load data if we have a data provider
		!d.HasData() { // We might have already been told to load the data so that another related control
		               // can access information in this control. For example, a paged control and a pager.
		               // This MANDATES that the control then unload the data after drawing

		log.FrameworkDebug("Calling BindData")
		dataProvider := owner.Page().GetControl(d.dataProviderID).(DataBinder)
		dataProvider.BindData(ctx, owner) // tell the data binder to call SetData on the given object, or load data some other way
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
		result := f(i + d.dataOffset, itemI)
		if !result {
			break
		}
	}
}

func (d *DataManager) HasData() bool {
	return d.data != nil
}

type encodedDataManager struct {
	DataProviderID string
	Data         interface{}
}

func (d *DataManager) Serialize(e page.Encoder) (err error) {
	enc := encodedDataManager{
		DataProviderID: d.dataProviderID,
		Data:           d.data,
	}
	if err = e.Encode(enc); err != nil {
		panic (err)
	}
	return
}

func (d *DataManager) Deserialize(dec page.Decoder) (err error) {
	enc := encodedDataManager{}

	if err = dec.Decode(&enc); err != nil {
		panic(err)
	}

	d.dataProviderID = enc.DataProviderID
	d.data = enc.Data
	return
}

