package control

import "context"

type DataBinder interface {
	BindData(ctx context.Context, s DataSetter)
}

// A DataSetter is the interface for the owner (the embedder) of the DataManager
type DataSetter interface {
	IDer
	SetDataProvider(b DataBinder)
	// SetData should be passed a slice of data items
	SetData([]interface{})
}

// DataManager is an object designed to be embedded in a control that will help manage the data binding process.
type DataManager struct {
	dataProvider DataBinder
	Data   []interface{}
}

func (d *DataManager) SetDataProvider(b DataBinder) {
	d.dataProvider = b
}

func (d *DataManager) SetData(data []interface{}) {
	d.Data = data
}

func (d *DataManager) ResetData() {
	d.Data = nil
}

// GetData tells the data binder to load data by calling SetData on the given object. The object should be
// the embedder of the DataManager
func (d *DataManager) GetData(ctx context.Context, owner DataSetter) {
	if d.dataProvider != nil {
		d.dataProvider.BindData(ctx, owner) // tell the data binder to call SetData on the given object, or load data some other way
	}
}

