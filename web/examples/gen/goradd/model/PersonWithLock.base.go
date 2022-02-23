package model

// Code generated by goradd. DO NOT EDIT.

import (
	"context"
	"fmt"

	"github.com/goradd/goradd/pkg/orm/broadcast"
	"github.com/goradd/goradd/pkg/orm/db"
	. "github.com/goradd/goradd/pkg/orm/op"
	"github.com/goradd/goradd/pkg/orm/query"
	"github.com/goradd/goradd/pkg/stringmap"
	"github.com/goradd/goradd/web/examples/gen/goradd/model/node"

	//"./node"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"time"

	time2 "github.com/goradd/goradd/pkg/time"
)

// personWithLockBase is a base structure to be embedded in a "subclass" and provides the ORM access to the database.
// Do not directly access the internal variables, but rather use the accessor functions, since this class maintains internal state
// related to the variables.

type personWithLockBase struct {
	id        string
	idIsValid bool
	idIsDirty bool

	firstName        string
	firstNameIsValid bool
	firstNameIsDirty bool

	lastName        string
	lastNameIsValid bool
	lastNameIsDirty bool

	sysTimestamp        time.Time
	sysTimestampIsNull  bool
	sysTimestampIsValid bool
	sysTimestampIsDirty bool

	// Custom aliases, if specified
	_aliases map[string]interface{}

	// Indicates whether this is a new object, or one loaded from the database. Used by Save to know whether to Insert or Update
	_restored bool

	// The original primary key for updates
	_originalPK string
}

const (
	PersonWithLockIDDefault           = ""
	PersonWithLockFirstNameDefault    = ""
	PersonWithLockLastNameDefault     = ""
	PersonWithLockSysTimestampDefault = time2.Zero
)

const (
	PersonWithLock_ID = `ID`

	PersonWithLock_FirstName = `FirstName`

	PersonWithLock_LastName = `LastName`

	PersonWithLock_SysTimestamp = `SysTimestamp`
)

// Initialize or re-initialize a PersonWithLock database object to default values.
func (o *personWithLockBase) Initialize() {

	o.id = ""
	o.idIsValid = false
	o.idIsDirty = false

	o.firstName = ""
	o.firstNameIsValid = false
	o.firstNameIsDirty = false

	o.lastName = ""
	o.lastNameIsValid = false
	o.lastNameIsDirty = false

	o.sysTimestamp = time.Time{}
	o.sysTimestampIsNull = true
	o.sysTimestampIsValid = true
	o.sysTimestampIsDirty = true

	o._restored = false
}

func (o *personWithLockBase) PrimaryKey() string {
	return o.id
}

// ID returns the loaded value of ID.
func (o *personWithLockBase) ID() string {
	return fmt.Sprint(o.id)
}

// IDIsValid returns true if the value was loaded from the database or has been set.
func (o *personWithLockBase) IDIsValid() bool {
	return o._restored && o.idIsValid
}

// FirstName returns the loaded value of FirstName.
func (o *personWithLockBase) FirstName() string {
	if o._restored && !o.firstNameIsValid {
		panic("firstName was not selected in the last query and has not been set, and so is not valid")
	}
	return o.firstName
}

// FirstNameIsValid returns true if the value was loaded from the database or has been set.
func (o *personWithLockBase) FirstNameIsValid() bool {
	return o.firstNameIsValid
}

// SetFirstName sets the value of FirstName in the object, to be saved later using the Save() function.
func (o *personWithLockBase) SetFirstName(v string) {
	o.firstNameIsValid = true
	if o.firstName != v || !o._restored {
		o.firstName = v
		o.firstNameIsDirty = true
	}

}

// LastName returns the loaded value of LastName.
func (o *personWithLockBase) LastName() string {
	if o._restored && !o.lastNameIsValid {
		panic("lastName was not selected in the last query and has not been set, and so is not valid")
	}
	return o.lastName
}

// LastNameIsValid returns true if the value was loaded from the database or has been set.
func (o *personWithLockBase) LastNameIsValid() bool {
	return o.lastNameIsValid
}

// SetLastName sets the value of LastName in the object, to be saved later using the Save() function.
func (o *personWithLockBase) SetLastName(v string) {
	o.lastNameIsValid = true
	if o.lastName != v || !o._restored {
		o.lastName = v
		o.lastNameIsDirty = true
	}

}

// SysTimestamp returns the loaded value of SysTimestamp.
func (o *personWithLockBase) SysTimestamp() time.Time {
	if o._restored && !o.sysTimestampIsValid {
		panic("sysTimestamp was not selected in the last query and has not been set, and so is not valid")
	}
	return o.sysTimestamp
}

// SysTimestampIsValid returns true if the value was loaded from the database or has been set.
func (o *personWithLockBase) SysTimestampIsValid() bool {
	return o.sysTimestampIsValid
}

// SysTimestampIsNull returns true if the related database value is null.
func (o *personWithLockBase) SysTimestampIsNull() bool {
	return o.sysTimestampIsNull
}

// SysTimestamp_I returns the loaded value of SysTimestamp as an interface.
// If the value in the database is NULL, a nil interface is returned.
func (o *personWithLockBase) SysTimestamp_I() interface{} {
	if o._restored && !o.sysTimestampIsValid {
		panic("sysTimestamp was not selected in the last query and has not been set, and so is not valid")
	} else if o.sysTimestampIsNull {
		return nil
	}
	return o.sysTimestamp
}

func (o *personWithLockBase) SetSysTimestamp(i interface{}) {
	o.sysTimestampIsValid = true
	if i == nil {
		if !o.sysTimestampIsNull {
			o.sysTimestampIsNull = true
			o.sysTimestampIsDirty = true
			o.sysTimestamp = time.Time{}
		}
	} else {
		v := i.(time.Time)
		if o.sysTimestampIsNull ||
			!o._restored ||
			o.sysTimestamp != v {

			o.sysTimestampIsNull = false
			o.sysTimestamp = v
			o.sysTimestampIsDirty = true
		}
	}
}

// GetAlias returns the alias for the given key.
func (o *personWithLockBase) GetAlias(key string) query.AliasValue {
	if a, ok := o._aliases[key]; ok {
		return query.NewAliasValue(a)
	} else {
		panic("Alias " + key + " not found.")
		return query.NewAliasValue([]byte{})
	}
}

// Load returns a PersonWithLock from the database.
// joinOrSelectNodes lets you provide nodes for joining to other tables or selecting specific fields. Table nodes will
// be considered Join nodes, and column nodes will be Select nodes. See Join() and Select() for more info.
func LoadPersonWithLock(ctx context.Context, primaryKey string, joinOrSelectNodes ...query.NodeI) *PersonWithLock {
	return queryPersonWithLocks(ctx).Where(Equal(node.PersonWithLock().ID(), primaryKey)).joinOrSelect(joinOrSelectNodes...).Get()
}

// LoadPersonWithLockByID queries for a single PersonWithLock object by the given unique index values.
// joinOrSelectNodes lets you provide nodes for joining to other tables or selecting specific fields. Table nodes will
// be considered Join nodes, and column nodes will be Select nodes. See Join() and Select() for more info.
// If you need a more elaborate query, use QueryPersonWithLocks() to start a query builder.
func LoadPersonWithLockByID(ctx context.Context, id string, joinOrSelectNodes ...query.NodeI) *PersonWithLock {
	q := queryPersonWithLocks(ctx)
	q = q.Where(Equal(node.PersonWithLock().ID(), id))
	return q.
		joinOrSelect(joinOrSelectNodes...).
		Get()
}

// HasPersonWithLockByID returns true if the
// given unique index values exist in the database.
func HasPersonWithLockByID(ctx context.Context, id string) bool {
	q := queryPersonWithLocks(ctx)
	q = q.Where(Equal(node.PersonWithLock().ID(), id))
	return q.Count(false) == 1
}

// The PersonWithLocksBuilder uses the QueryBuilderI interface from the database to build a query.
// All query operations go through this query builder.
// End a query by calling either Load, Count, or Delete
type PersonWithLocksBuilder struct {
	builder             query.QueryBuilderI
	hasConditionalJoins bool
}

func newPersonWithLockBuilder(ctx context.Context) *PersonWithLocksBuilder {
	b := &PersonWithLocksBuilder{
		builder: db.GetDatabase("goradd").NewBuilder(ctx),
	}
	return b.Join(node.PersonWithLock())
}

// Load terminates the query builder, performs the query, and returns a slice of PersonWithLock objects. If there are
// any errors, they are returned in the context object. If no results come back from the query, it will return
// an empty slice
func (b *PersonWithLocksBuilder) Load() (personWithLockSlice []*PersonWithLock) {
	results := b.builder.Load()
	if results == nil {
		return
	}
	for _, item := range results {
		o := new(PersonWithLock)
		o.load(item, o, nil, "")
		personWithLockSlice = append(personWithLockSlice, o)
	}
	return personWithLockSlice
}

// LoadI terminates the query builder, performs the query, and returns a slice of interfaces. If there are
// any errors, they are returned in the context object. If no results come back from the query, it will return
// an empty slice.
func (b *PersonWithLocksBuilder) LoadI() (personWithLockSlice []interface{}) {
	results := b.builder.Load()
	if results == nil {
		return
	}
	for _, item := range results {
		o := new(PersonWithLock)
		o.load(item, o, nil, "")
		personWithLockSlice = append(personWithLockSlice, o)
	}
	return personWithLockSlice
}

// Get is a convenience method to return only the first item found in a query.
// The entire query is performed, so you should generally use this only if you know
// you are selecting on one or very few items.
func (b *PersonWithLocksBuilder) Get() *PersonWithLock {
	results := b.Load()
	if results != nil && len(results) > 0 {
		obj := results[0]
		return obj
	} else {
		return nil
	}
}

// Expand expands an array type node so that it will produce individual rows instead of an array of items
func (b *PersonWithLocksBuilder) Expand(n query.NodeI) *PersonWithLocksBuilder {
	b.builder.Expand(n)
	return b
}

// Join adds a node to the node tree so that its fields will appear in the query. Optionally add conditions to filter
// what gets included. The conditions will be AND'd with the basic condition matching the primary keys of the join.
func (b *PersonWithLocksBuilder) Join(n query.NodeI, conditions ...query.NodeI) *PersonWithLocksBuilder {
	var condition query.NodeI
	if len(conditions) > 1 {
		condition = And(conditions)
	} else if len(conditions) == 1 {
		condition = conditions[0]
	}
	b.builder.Join(n, condition)
	if condition != nil {
		b.hasConditionalJoins = true
	}
	return b
}

// Where adds a condition to filter what gets selected.
func (b *PersonWithLocksBuilder) Where(c query.NodeI) *PersonWithLocksBuilder {
	b.builder.Condition(c)
	return b
}

// OrderBy specifies how the resulting data should be sorted.
func (b *PersonWithLocksBuilder) OrderBy(nodes ...query.NodeI) *PersonWithLocksBuilder {
	b.builder.OrderBy(nodes...)
	return b
}

// Limit will return a subset of the data, limited to the offset and number of rows specified
func (b *PersonWithLocksBuilder) Limit(maxRowCount int, offset int) *PersonWithLocksBuilder {
	b.builder.Limit(maxRowCount, offset)
	return b
}

// Select optimizes the query to only return the specified fields. Once you put a Select in your query, you must
// specify all the fields that you will eventually read out. Be careful when selecting fields in joined tables, as joined
// tables will also contain pointers back to the parent table, and so the parent node should have the same field selected
// as the child node if you are querying those fields.
func (b *PersonWithLocksBuilder) Select(nodes ...query.NodeI) *PersonWithLocksBuilder {
	b.builder.Select(nodes...)
	return b
}

// Alias lets you add a node with a custom name. After the query, you can read out the data using GetAlias() on a
// returned object. Alias is useful for adding calculations or subqueries to the query.
func (b *PersonWithLocksBuilder) Alias(name string, n query.NodeI) *PersonWithLocksBuilder {
	b.builder.Alias(name, n)
	return b
}

// Distinct removes duplicates from the results of the query. Adding a Select() may help you get to the data you want, although
// using Distinct with joined tables is often not effective, since we force joined tables to include primary keys in the query, and this
// often ruins the effect of Distinct.
func (b *PersonWithLocksBuilder) Distinct() *PersonWithLocksBuilder {
	b.builder.Distinct()
	return b
}

// GroupBy controls how results are grouped when using aggregate functions in an Alias() call.
func (b *PersonWithLocksBuilder) GroupBy(nodes ...query.NodeI) *PersonWithLocksBuilder {
	b.builder.GroupBy(nodes...)
	return b
}

// Having does additional filtering on the results of the query.
func (b *PersonWithLocksBuilder) Having(node query.NodeI) *PersonWithLocksBuilder {
	b.builder.Having(node)
	return b
}

// Count terminates a query and returns just the number of items selected.
//
// distinct wll count the number of distinct items, ignoring duplicates.
//
// nodes will select individual fields, and should be accompanied by a GroupBy.
func (b *PersonWithLocksBuilder) Count(distinct bool, nodes ...query.NodeI) uint {
	return b.builder.Count(distinct, nodes...)
}

// Delete uses the query builder to delete a group of records that match the criteria
func (b *PersonWithLocksBuilder) Delete() {
	b.builder.Delete()
	broadcast.BulkChange(b.builder.Context(), "goradd", "person_with_lock")
}

// Subquery uses the query builder to define a subquery within a larger query. You MUST include what
// you are selecting by adding Alias or Select functions on the subquery builder. Generally you would use
// this as a node to an Alias function on the surrounding query builder.
func (b *PersonWithLocksBuilder) Subquery() *query.SubqueryNode {
	return b.builder.Subquery()
}

// joinOrSelect is a private helper function for the Load* functions
func (b *PersonWithLocksBuilder) joinOrSelect(nodes ...query.NodeI) *PersonWithLocksBuilder {
	for _, n := range nodes {
		switch n.(type) {
		case query.TableNodeI:
			b.builder.Join(n, nil)
		case *query.ColumnNode:
			b.Select(n)
		}
	}
	return b
}

func CountPersonWithLockByID(ctx context.Context, id string) uint {
	return queryPersonWithLocks(ctx).Where(Equal(node.PersonWithLock().ID(), id)).Count(false)
}

func CountPersonWithLockByFirstName(ctx context.Context, firstName string) uint {
	return queryPersonWithLocks(ctx).Where(Equal(node.PersonWithLock().FirstName(), firstName)).Count(false)
}

func CountPersonWithLockByLastName(ctx context.Context, lastName string) uint {
	return queryPersonWithLocks(ctx).Where(Equal(node.PersonWithLock().LastName(), lastName)).Count(false)
}

func CountPersonWithLockBySysTimestamp(ctx context.Context, sysTimestamp time.Time) uint {
	return queryPersonWithLocks(ctx).Where(Equal(node.PersonWithLock().SysTimestamp(), sysTimestamp)).Count(false)
}

// load is the private loader that transforms data coming from the database into a tree structure reflecting the relationships
// between the object chain requested by the user in the query.
// Care must be taken in the query, as Select clauses might not be honored if the child object has fields selected which the parent object does not have.
func (o *personWithLockBase) load(m map[string]interface{}, objThis *PersonWithLock, objParent interface{}, parentKey string) {
	if v, ok := m["id"]; ok && v != nil {
		if o.id, ok = v.(string); ok {
			o.idIsValid = true
			o.idIsDirty = false
			o._originalPK = o.id
		} else {
			panic("Wrong type found for id.")
		}
	} else {
		o.idIsValid = false
		o.id = ""
	}

	if v, ok := m["first_name"]; ok && v != nil {
		if o.firstName, ok = v.(string); ok {
			o.firstNameIsValid = true
			o.firstNameIsDirty = false
		} else {
			panic("Wrong type found for first_name.")
		}
	} else {
		o.firstNameIsValid = false
		o.firstName = ""
	}

	if v, ok := m["last_name"]; ok && v != nil {
		if o.lastName, ok = v.(string); ok {
			o.lastNameIsValid = true
			o.lastNameIsDirty = false
		} else {
			panic("Wrong type found for last_name.")
		}
	} else {
		o.lastNameIsValid = false
		o.lastName = ""
	}

	if v, ok := m["sys_timestamp"]; ok {
		if v == nil {
			o.sysTimestamp = time.Time{}
			o.sysTimestampIsNull = true
			o.sysTimestampIsValid = true
			o.sysTimestampIsDirty = false
		} else if o.sysTimestamp, ok = v.(time.Time); ok {
			o.sysTimestampIsNull = false
			o.sysTimestampIsValid = true
			o.sysTimestampIsDirty = false
		} else {
			panic("Wrong type found for sys_timestamp.")
		}
	} else {
		o.sysTimestampIsValid = false
		o.sysTimestampIsNull = true
		o.sysTimestamp = time.Time{}
	}

	if v, ok := m["aliases_"]; ok {
		o._aliases = map[string]interface{}(v.(db.ValueMap))
	}
	o._restored = true
}

// Save will update or insert the object, depending on the state of the object.
// If it has any auto-generated ids, those will be updated.
func (o *personWithLockBase) Save(ctx context.Context) {
	if o._restored {
		o.update(ctx)
	} else {
		o.insert(ctx)
	}
}

// update will update the values in the database, saving any changed values.
func (o *personWithLockBase) update(ctx context.Context) {
	var modifiedFields map[string]interface{}
	d := Database()
	db.ExecuteTransaction(ctx, d, func() {

		if !o._restored {
			panic("Cannot update a record that was not originally read from the database.")
		}

		modifiedFields = o.getModifiedFields()
		if len(modifiedFields) != 0 {
			d.Update(ctx, "person_with_lock", modifiedFields, "id", o._originalPK)
		}

	}) // transaction
	o.resetDirtyStatus()
	if len(modifiedFields) != 0 {
		broadcast.Update(ctx, "goradd", "person_with_lock", o._originalPK, stringmap.SortedKeys(modifiedFields)...)
	}
}

// insert will insert the item into the database. Related items will be saved.
func (o *personWithLockBase) insert(ctx context.Context) {
	d := Database()
	db.ExecuteTransaction(ctx, d, func() {

		if !o.firstNameIsValid {
			panic("a value for FirstName is required, and there is no default value. Call SetFirstName() before inserting the record.")
		}

		if !o.lastNameIsValid {
			panic("a value for LastName is required, and there is no default value. Call SetLastName() before inserting the record.")
		}

		m := o.getValidFields()

		id := d.Insert(ctx, "person_with_lock", m)
		o.id = id
		o._originalPK = id

	}) // transaction
	o.resetDirtyStatus()
	o._restored = true
	broadcast.Insert(ctx, "goradd", "person_with_lock", o.PrimaryKey())
}

func (o *personWithLockBase) getModifiedFields() (fields map[string]interface{}) {
	fields = map[string]interface{}{}
	if o.idIsDirty {

		fields["id"] = o.id

	}
	if o.firstNameIsDirty {

		fields["first_name"] = o.firstName

	}
	if o.lastNameIsDirty {

		fields["last_name"] = o.lastName

	}
	if o.sysTimestampIsDirty {

		if o.sysTimestampIsNull {
			fields["sys_timestamp"] = nil
		} else {
			fields["sys_timestamp"] = o.sysTimestamp
		}

	}
	return
}

func (o *personWithLockBase) getValidFields() (fields map[string]interface{}) {
	fields = map[string]interface{}{}
	if o.firstNameIsValid {

		fields["first_name"] = o.firstName

	}
	if o.lastNameIsValid {

		fields["last_name"] = o.lastName

	}
	if o.sysTimestampIsValid {

		if o.sysTimestampIsNull {
			fields["sys_timestamp"] = nil
		} else {
			fields["sys_timestamp"] = o.sysTimestamp
		}

	}
	return
}

// Delete deletes the associated record from the database.
func (o *personWithLockBase) Delete(ctx context.Context) {
	if !o._restored {
		panic("Cannot delete a record that has no primary key value.")
	}
	d := Database()
	d.Delete(ctx, "person_with_lock", "id", o.id)
	broadcast.Delete(ctx, "goradd", "person_with_lock", fmt.Sprint(o.id))
}

// deletePersonWithLock deletes the associated record from the database.
func deletePersonWithLock(ctx context.Context, pk string) {
	d := db.GetDatabase("goradd")
	d.Delete(ctx, "person_with_lock", "id", pk)
	broadcast.Delete(ctx, "goradd", "person_with_lock", fmt.Sprint(pk))
}

func (o *personWithLockBase) resetDirtyStatus() {
	o.idIsDirty = false
	o.firstNameIsDirty = false
	o.lastNameIsDirty = false
	o.sysTimestampIsDirty = false

}

func (o *personWithLockBase) IsDirty() bool {
	return o.idIsDirty ||
		o.firstNameIsDirty ||
		o.lastNameIsDirty ||
		o.sysTimestampIsDirty

}

// Get returns the value of a field in the object based on the field's name.
// It will also get related objects if they are loaded.
// Invalid fields and objects are returned as nil
func (o *personWithLockBase) Get(key string) interface{} {

	switch key {
	case "ID":
		if !o.idIsValid {
			return nil
		}
		return o.id

	case "FirstName":
		if !o.firstNameIsValid {
			return nil
		}
		return o.firstName

	case "LastName":
		if !o.lastNameIsValid {
			return nil
		}
		return o.lastName

	case "SysTimestamp":
		if !o.sysTimestampIsValid {
			return nil
		}
		return o.sysTimestamp

	}
	return nil
}

// MarshalBinary serializes the object into a buffer that is deserializable using UnmarshalBinary.
// It should be used for transmitting database objects over the wire, or for temporary storage. It does not send
// a version number, so if the data format changes, its up to you to invalidate the old stored objects.
// The framework uses this to serialize the object when it is stored in a control.
func (o *personWithLockBase) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)

	if err := encoder.Encode(o.id); err != nil {
		return nil, err
	}
	if err := encoder.Encode(o.idIsValid); err != nil {
		return nil, err
	}
	if err := encoder.Encode(o.idIsDirty); err != nil {
		return nil, err
	}

	if err := encoder.Encode(o.firstName); err != nil {
		return nil, err
	}
	if err := encoder.Encode(o.firstNameIsValid); err != nil {
		return nil, err
	}
	if err := encoder.Encode(o.firstNameIsDirty); err != nil {
		return nil, err
	}

	if err := encoder.Encode(o.lastName); err != nil {
		return nil, err
	}
	if err := encoder.Encode(o.lastNameIsValid); err != nil {
		return nil, err
	}
	if err := encoder.Encode(o.lastNameIsDirty); err != nil {
		return nil, err
	}

	if err := encoder.Encode(o.sysTimestamp); err != nil {
		return nil, err
	}
	if err := encoder.Encode(o.sysTimestampIsNull); err != nil {
		return nil, err
	}
	if err := encoder.Encode(o.sysTimestampIsValid); err != nil {
		return nil, err
	}
	if err := encoder.Encode(o.sysTimestampIsDirty); err != nil {
		return nil, err
	}

	if o._aliases == nil {
		if err := encoder.Encode(false); err != nil {
			return nil, err
		}
	} else {
		if err := encoder.Encode(true); err != nil {
			return nil, err
		}
		if err := encoder.Encode(o._aliases); err != nil {
			return nil, err
		}
	}

	if err := encoder.Encode(o._restored); err != nil {
		return nil, err
	}
	if err := encoder.Encode(o._originalPK); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (o *personWithLockBase) UnmarshalBinary(data []byte) (err error) {

	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	var isPtr bool

	_ = isPtr

	if err = dec.Decode(&o.id); err != nil {
		return
	}
	if err = dec.Decode(&o.idIsValid); err != nil {
		return
	}
	if err = dec.Decode(&o.idIsDirty); err != nil {
		return
	}

	if err = dec.Decode(&o.firstName); err != nil {
		return
	}
	if err = dec.Decode(&o.firstNameIsValid); err != nil {
		return
	}
	if err = dec.Decode(&o.firstNameIsDirty); err != nil {
		return
	}

	if err = dec.Decode(&o.lastName); err != nil {
		return
	}
	if err = dec.Decode(&o.lastNameIsValid); err != nil {
		return
	}
	if err = dec.Decode(&o.lastNameIsDirty); err != nil {
		return
	}

	if err = dec.Decode(&o.sysTimestamp); err != nil {
		return
	}
	if err = dec.Decode(&o.sysTimestampIsNull); err != nil {
		return
	}
	if err = dec.Decode(&o.sysTimestampIsValid); err != nil {
		return
	}
	if err = dec.Decode(&o.sysTimestampIsDirty); err != nil {
		return
	}

	if err = dec.Decode(&isPtr); err != nil {
		return
	}
	if isPtr {
		if err = dec.Decode(&o._aliases); err != nil {
			return
		}
	}

	if err = dec.Decode(&o._restored); err != nil {
		return
	}
	if err = dec.Decode(&o._originalPK); err != nil {
		return
	}

	return
}

// MarshalJSON serializes the object into a JSON object.
// Only valid data will be serialized, meaning, you can control what gets serialized by using Select to
// select only the fields you want when you query for the object. Another way to control the output
// is to call MarshalStringMap, modify the map, then encode the map.
func (o *personWithLockBase) MarshalJSON() (data []byte, err error) {
	v := o.MarshalStringMap()
	return json.Marshal(v)
}

// MarshalStringMap serializes the object into a string map of interfaces.
// Only valid data will be serialized, meaning, you can control what gets serialized by using Select to
// select only the fields you want when you query for the object. The keys are the same as the json keys.
func (o *personWithLockBase) MarshalStringMap() map[string]interface{} {
	v := make(map[string]interface{})

	if o.idIsValid {
		v["id"] = o.id
	}

	if o.firstNameIsValid {
		v["firstName"] = o.firstName
	}

	if o.lastNameIsValid {
		v["lastName"] = o.lastName
	}

	if o.sysTimestampIsValid {
		if o.sysTimestampIsNull {
			v["sysTimestamp"] = nil
		} else {
			v["sysTimestamp"] = o.sysTimestamp
		}
	}

	for _k, _v := range o._aliases {
		v[_k] = _v
	}
	return v
}

// UnmarshalJSON unmarshalls the given json data into the personWithLock. The personWithLock can be a
// newly created object, or one loaded from the database.
//
// After unmarshalling, the object is not  saved. You must call Save to insert it into the database
// or update it.
//
// Unmarshalling of sub-objects, as in objects linked via foreign keys, is not currently supported.
//
// The fields it expects are:
//   "id" - string

//   "firstName" - string

//   "lastName" - string

//   "sysTimestamp" - time.Time, nullable

func (o *personWithLockBase) UnmarshalJSON(data []byte) (err error) {
	var v map[string]interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	return o.UnmarshalStringMap(v)
}

// UnmarshalStringMap will load the values from the stringmap into the object.
//
// Override this in personWithLock to modify the json before sending it here.
func (o *personWithLockBase) UnmarshalStringMap(m map[string]interface{}) (err error) {
	for k, v := range m {
		switch k {
		case "firstName":
			{
				if v == nil {
					return fmt.Errorf("json field %s cannot be null", k)
				}
				if s, ok := v.(string); !ok {
					return fmt.Errorf("json field %s must be a string", k)
				} else {
					o.SetFirstName(s)
				}
			}
		case "lastName":
			{
				if v == nil {
					return fmt.Errorf("json field %s cannot be null", k)
				}
				if s, ok := v.(string); !ok {
					return fmt.Errorf("json field %s must be a string", k)
				} else {
					o.SetLastName(s)
				}
			}
		case "sysTimestamp":
			{
				if v == nil {
					o.SetSysTimestamp(v)
					continue
				}
				switch d := v.(type) {
				case float64:
					// a numeric value, which for JSON, means milliseconds since epoc
					o.SetSysTimestamp(time.UnixMilli(int64(d)).UTC())
				case string:
					// an ISO8601 string (hopefully)
					var t time.Time
					err = t.UnmarshalJSON([]byte(`"` + d + `"`))
					if err != nil {
						return fmt.Errorf("JSON format error for time field %s: %w", k, err)
					}
					t = t.UTC()
					o.SetSysTimestamp(t)
				default:
					return fmt.Errorf("json field %s must be a number or a string", k)
				}
			}

		}
	}
	return
}

// Custom functions. See goradd/codegen/templates/orm/modelBase.
