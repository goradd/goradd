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
	"github.com/goradd/goradd/web/examples/gen/goraddUnit/model/node"

	//"./node"
	"bytes"
	"encoding/gob"
	"encoding/json"
)

// reverseBase is a base structure to be embedded in a "subclass" and provides the ORM access to the database.
// Do not directly access the internal variables, but rather use the accessor functions, since this class maintains internal state
// related to the variables.

type reverseBase struct {
	id        string
	idIsValid bool
	idIsDirty bool

	name        string
	nameIsValid bool
	nameIsDirty bool

	// Reverse reference objects.
	oForwardsAsNotNull      []*Forward          // Objects in the order they were queried
	mForwardsAsNotNull      map[string]*Forward // Objects by PK
	oForwardAsUniqueNotNull *Forward
	oForwardsAsNull         []*Forward          // Objects in the order they were queried
	mForwardsAsNull         map[string]*Forward // Objects by PK
	oForwardAsUniqueNull    *Forward

	// Custom aliases, if specified
	_aliases map[string]interface{}

	// Indicates whether this is a new object, or one loaded from the database. Used by Save to know whether to Insert or Update
	_restored bool
}

const (
	ReverseIDDefault   = ""
	ReverseNameDefault = ""
)

const (
	ReverseID                = `ID`
	ReverseName              = `Name`
	ReverseForwardsAsNotNull = `ForwardsAsNotNull`

	ReverseForwardAsUniqueNotNull = `ForwardAsUniqueNotNull`
	ReverseForwardsAsNull         = `ForwardsAsNull`

	ReverseForwardAsUniqueNull = `ForwardAsUniqueNull`
)

// Initialize or re-initialize a Reverse database object to default values.
func (o *reverseBase) Initialize() {

	o.id = ""
	o.idIsValid = false
	o.idIsDirty = false

	o.name = ""
	o.nameIsValid = false
	o.nameIsDirty = false

	o._restored = false
}

func (o *reverseBase) PrimaryKey() string {
	return o.id
}

// ID returns the loaded value of ID.
func (o *reverseBase) ID() string {
	return fmt.Sprint(o.id)
}

// IDIsValid returns true if the value was loaded from the database or has been set.
func (o *reverseBase) IDIsValid() bool {
	return o._restored && o.idIsValid
}

func (o *reverseBase) Name() string {
	if o._restored && !o.nameIsValid {
		panic("name was not selected in the last query and so is not valid")
	}
	return o.name
}

// NameIsValid returns true if the value was loaded from the database or has been set.
func (o *reverseBase) NameIsValid() bool {
	return o.nameIsValid
}

// SetName sets the value of Name in the object, to be saved later using the Save() function.
func (o *reverseBase) SetName(v string) {
	o.nameIsValid = true
	if o.name != v || !o._restored {
		o.name = v
		o.nameIsDirty = true
	}

}

// GetAlias returns the alias for the given key.
func (o *reverseBase) GetAlias(key string) query.AliasValue {
	if a, ok := o._aliases[key]; ok {
		return query.NewAliasValue(a)
	} else {
		panic("Alias " + key + " not found.")
		return query.NewAliasValue([]byte{})
	}
}

// ForwardAsNotNull returns a single Forward object by primary key, if one was loaded.
// Otherwise, it will return nil.
func (o *reverseBase) ForwardAsNotNull(pk string) *Forward {
	if o.oForwardsAsNotNull == nil || len(o.oForwardsAsNotNull) == 0 {
		return nil
	}
	v, _ := o.mForwardsAsNotNull[pk]
	return v
}

// ForwardsAsNotNull returns a slice of Forward objects if loaded.
func (o *reverseBase) ForwardsAsNotNull() []*Forward {
	if o.oForwardsAsNotNull == nil {
		return nil
	}
	return o.oForwardsAsNotNull
}

// LoadForwardsAsNotNull loads a new slice of Forward objects and returns it.
func (o *reverseBase) LoadForwardsAsNotNull(ctx context.Context, conditions ...interface{}) []*Forward {
	qb := queryForwards(ctx)
	cond := Equal(node.Forward().ReverseNotNullID(), o.PrimaryKey())
	if conditions != nil {
		conditions = append(conditions, cond)
		cond = And(conditions...)
	}

	o.oForwardsAsNotNull = qb.Where(cond).Load(ctx)
	return o.oForwardsAsNotNull
}

// ForwardAsUniqueNotNull returns the connected Forward object, if one was loaded
// otherwise, it will return nil.
func (o *reverseBase) ForwardAsUniqueNotNull() *Forward {
	if o.oForwardAsUniqueNotNull == nil {
		return nil
	}
	return o.oForwardAsUniqueNotNull
}

// LoadForwardAsUniqueNotNull returns the connected Forward object, if one was loaded
// otherwise, it will return nil.
func (o *reverseBase) LoadForwardAsUniqueNotNull(ctx context.Context) *Forward {
	if o.oForwardAsUniqueNotNull == nil {
		o.oForwardAsUniqueNotNull = LoadForwardByReverseUniqueNotNullID(ctx, o.ID())
	}
	return o.oForwardAsUniqueNotNull
}

// ForwardAsNull returns a single Forward object by primary key, if one was loaded.
// Otherwise, it will return nil.
func (o *reverseBase) ForwardAsNull(pk string) *Forward {
	if o.oForwardsAsNull == nil || len(o.oForwardsAsNull) == 0 {
		return nil
	}
	v, _ := o.mForwardsAsNull[pk]
	return v
}

// ForwardsAsNull returns a slice of Forward objects if loaded.
func (o *reverseBase) ForwardsAsNull() []*Forward {
	if o.oForwardsAsNull == nil {
		return nil
	}
	return o.oForwardsAsNull
}

// LoadForwardsAsNull loads a new slice of Forward objects and returns it.
func (o *reverseBase) LoadForwardsAsNull(ctx context.Context, conditions ...interface{}) []*Forward {
	qb := queryForwards(ctx)
	cond := Equal(node.Forward().ReverseNullID(), o.PrimaryKey())
	if conditions != nil {
		conditions = append(conditions, cond)
		cond = And(conditions...)
	}

	o.oForwardsAsNull = qb.Where(cond).Load(ctx)
	return o.oForwardsAsNull
}

// ForwardAsUniqueNull returns the connected Forward object, if one was loaded
// otherwise, it will return nil.
func (o *reverseBase) ForwardAsUniqueNull() *Forward {
	if o.oForwardAsUniqueNull == nil {
		return nil
	}
	return o.oForwardAsUniqueNull
}

// LoadForwardAsUniqueNull returns the connected Forward object, if one was loaded
// otherwise, it will return nil.
func (o *reverseBase) LoadForwardAsUniqueNull(ctx context.Context) *Forward {
	if o.oForwardAsUniqueNull == nil {
		o.oForwardAsUniqueNull = LoadForwardByReverseUniqueNullID(ctx, o.ID())
	}
	return o.oForwardAsUniqueNull
}

// Load returns a Reverse from the database.
// joinOrSelectNodes lets you provide nodes for joining to other tables or selecting specific fields. Table nodes will
// be considered Join nodes, and column nodes will be Select nodes. See Join() and Select() for more info.
func LoadReverse(ctx context.Context, primaryKey string, joinOrSelectNodes ...query.NodeI) *Reverse {
	return queryReverses(ctx).Where(Equal(node.Reverse().ID(), primaryKey)).joinOrSelect(joinOrSelectNodes...).Get(ctx)
}

// The ReversesBuilder uses the QueryBuilderI interface from the database to build a query.
// All query operations go through this query builder.
// End a query by calling either Load, Count, or Delete
type ReversesBuilder struct {
	base                query.QueryBuilderI
	hasConditionalJoins bool
}

func newReverseBuilder() *ReversesBuilder {
	b := &ReversesBuilder{
		base: db.GetDatabase("goraddUnit").
			NewBuilder(),
	}
	return b.Join(node.Reverse())
}

// Load terminates the query builder, performs the query, and returns a slice of Reverse objects. If there are
// any errors, they are returned in the context object. If no results come back from the query, it will return
// an empty slice
func (b *ReversesBuilder) Load(ctx context.Context) (reverseSlice []*Reverse) {
	results := b.base.Load(ctx)
	if results == nil {
		return
	}
	for _, item := range results {
		o := new(Reverse)
		o.load(item, !b.hasConditionalJoins, o, nil, "")
		reverseSlice = append(reverseSlice, o)
	}
	return reverseSlice
}

// LoadI terminates the query builder, performs the query, and returns a slice of interfaces. If there are
// any errors, they are returned in the context object. If no results come back from the query, it will return
// an empty slice.
func (b *ReversesBuilder) LoadI(ctx context.Context) (reverseSlice []interface{}) {
	results := b.base.Load(ctx)
	if results == nil {
		return
	}
	for _, item := range results {
		o := new(Reverse)
		o.load(item, !b.hasConditionalJoins, o, nil, "")
		reverseSlice = append(reverseSlice, o)
	}
	return reverseSlice
}

// Get is a convenience method to return only the first item found in a query. It is equivalent to adding
// Limit(1,0) to the query, and then getting the first item from the returned slice.
// Limits with joins do not currently work, so don't try it if you have a join
// TODO: Change this to Load1 to be more descriptive and avoid confusion with other Getters
func (b *ReversesBuilder) Get(ctx context.Context) *Reverse {
	results := b.Limit(1, 0).Load(ctx)
	if results != nil && len(results) > 0 {
		obj := results[0]
		return obj
	} else {
		return nil
	}
}

// Expand expands an array type node so that it will produce individual rows instead of an array of items
func (b *ReversesBuilder) Expand(n query.NodeI) *ReversesBuilder {
	b.base.Expand(n)
	return b
}

// Join adds a node to the node tree so that its fields will appear in the query. Optionally add conditions to filter
// what gets included. The conditions will be AND'd with the basic condition matching the primary keys of the join.
func (b *ReversesBuilder) Join(n query.NodeI, conditions ...query.NodeI) *ReversesBuilder {
	var condition query.NodeI
	if len(conditions) > 1 {
		condition = And(conditions)
	} else if len(conditions) == 1 {
		condition = conditions[0]
	}
	b.base.Join(n, condition)
	if condition != nil {
		b.hasConditionalJoins = true
	}
	return b
}

// Where adds a condition to filter what gets selected.
func (b *ReversesBuilder) Where(c query.NodeI) *ReversesBuilder {
	b.base.Condition(c)
	return b
}

// OrderBy specifies how the resulting data should be sorted.
func (b *ReversesBuilder) OrderBy(nodes ...query.NodeI) *ReversesBuilder {
	b.base.OrderBy(nodes...)
	return b
}

// Limit will return a subset of the data, limited to the offset and number of rows specified
func (b *ReversesBuilder) Limit(maxRowCount int, offset int) *ReversesBuilder {
	b.base.Limit(maxRowCount, offset)
	return b
}

// Select optimizes the query to only return the specified fields. Once you put a Select in your query, you must
// specify all the fields that you will eventually read out. Be careful when selecting fields in joined tables, as joined
// tables will also contain pointers back to the parent table, and so the parent node should have the same field selected
// as the child node if you are querying those fields.
func (b *ReversesBuilder) Select(nodes ...query.NodeI) *ReversesBuilder {
	b.base.Select(nodes...)
	return b
}

// Alias lets you add a node with a custom name. After the query, you can read out the data using GetAlias() on a
// returned object. Alias is useful for adding calculations or subqueries to the query.
func (b *ReversesBuilder) Alias(name string, n query.NodeI) *ReversesBuilder {
	b.base.Alias(name, n)
	return b
}

// Distinct removes duplicates from the results of the query. Adding a Select() may help you get to the data you want, although
// using Distinct with joined tables is often not effective, since we force joined tables to include primary keys in the query, and this
// often ruins the effect of Distinct.
func (b *ReversesBuilder) Distinct() *ReversesBuilder {
	b.base.Distinct()
	return b
}

// GroupBy controls how results are grouped when using aggregate functions in an Alias() call.
func (b *ReversesBuilder) GroupBy(nodes ...query.NodeI) *ReversesBuilder {
	b.base.GroupBy(nodes...)
	return b
}

// Having does additional filtering on the results of the query.
func (b *ReversesBuilder) Having(node query.NodeI) *ReversesBuilder {
	b.base.Having(node)
	return b
}

// Count terminates a query and returns just the number of items selected.
func (b *ReversesBuilder) Count(ctx context.Context, distinct bool, nodes ...query.NodeI) uint {
	return b.base.Count(ctx, distinct, nodes...)
}

// Delete uses the query builder to delete a group of records that match the criteria
func (b *ReversesBuilder) Delete(ctx context.Context) {
	b.base.Delete(ctx)
}

// Subquery uses the query builder to define a subquery within a larger query. You MUST include what
// you are selecting by adding Alias or Select functions on the subquery builder. Generally you would use
// this as a node to an Alias function on the surrounding query builder.
func (b *ReversesBuilder) Subquery() *query.SubqueryNode {
	return b.base.Subquery()
}

// joinOrSelect is a private helper function for the Load* functions
func (b *ReversesBuilder) joinOrSelect(nodes ...query.NodeI) *ReversesBuilder {
	for _, n := range nodes {
		switch n.(type) {
		case query.TableNodeI:
			b.base.Join(n, nil)
		case *query.ColumnNode:
			b.Select(n)
		}
	}
	return b
}

func CountReverseByID(ctx context.Context, id string) uint {
	return queryReverses(ctx).Where(Equal(node.Reverse().ID(), id)).Count(ctx, false)
}

func CountReverseByName(ctx context.Context, name string) uint {
	return queryReverses(ctx).Where(Equal(node.Reverse().Name(), name)).Count(ctx, false)
}

// load is the private loader that transforms data coming from the database into a tree structure reflecting the relationships
// between the object chain requested by the user in the query.
// If linkParent is true we will have child relationships use a pointer back to the parent object. If false, it will create a separate object.
// Care must be taken in the query, as Select clauses might not be honored if the child object has fields selected which the parent object does not have.
// Also, if any joins are conditional, that might affect which child objects are included, so in this situation, linkParent should be false
func (o *reverseBase) load(m map[string]interface{}, linkParent bool, objThis *Reverse, objParent interface{}, parentKey string) {
	if v, ok := m["id"]; ok && v != nil {
		if o.id, ok = v.(string); ok {
			o.idIsValid = true
			o.idIsDirty = false
		} else {
			panic("Wrong type found for id.")
		}
	} else {
		o.idIsValid = false
		o.id = ""
	}

	if v, ok := m["name"]; ok && v != nil {
		if o.name, ok = v.(string); ok {
			o.nameIsValid = true
			o.nameIsDirty = false
		} else {
			panic("Wrong type found for name.")
		}
	} else {
		o.nameIsValid = false
		o.name = ""
	}

	if v, ok := m["ForwardsAsNotNull"]; ok {
		switch oForwardsAsNotNull := v.(type) {
		case []db.ValueMap:
			o.oForwardsAsNotNull = []*Forward{}
			o.mForwardsAsNotNull = map[string]*Forward{}
			for _, v2 := range oForwardsAsNotNull {
				obj := new(Forward)
				obj.load(v2, linkParent, obj, objThis, "ReverseNotNull")
				if linkParent && parentKey == "ForwardsAsNotNull" && obj.reverseNotNullID == objParent.(*Forward).reverseNotNullID {
					obj = objParent.(*Forward)
				}
				o.oForwardsAsNotNull = append(o.oForwardsAsNotNull, obj)
				o.mForwardsAsNotNull[obj.PrimaryKey()] = obj
			}
		case db.ValueMap: // single expansion
			obj := new(Forward)
			obj.load(oForwardsAsNotNull, linkParent, obj, objThis, "ReverseNotNull")
			if linkParent && parentKey == "ForwardsAsNotNull" && obj.reverseNotNullID == objParent.(*Forward).reverseNotNullID {
				obj = objParent.(*Forward)
			}
			o.oForwardsAsNotNull = []*Forward{obj}
		default:
			panic("Wrong type found for oForwardsAsNotNull object.")
		}
	} else {
		o.oForwardsAsNotNull = nil
	}

	if v, ok := m["ForwardAsUniqueNotNull"]; ok {
		if oForwardAsUniqueNotNull, ok2 := v.(db.ValueMap); ok2 {
			o.oForwardAsUniqueNotNull = new(Forward)
			o.oForwardAsUniqueNotNull.load(oForwardAsUniqueNotNull, linkParent, o.oForwardAsUniqueNotNull, objThis, "ReverseUniqueNotNull")
		} else {
			panic("Wrong type found for oForwardAsUniqueNotNull object.")
		}
	} else {
		o.oForwardAsUniqueNotNull = nil
	}

	if v, ok := m["ForwardsAsNull"]; ok {
		switch oForwardsAsNull := v.(type) {
		case []db.ValueMap:
			o.oForwardsAsNull = []*Forward{}
			o.mForwardsAsNull = map[string]*Forward{}
			for _, v2 := range oForwardsAsNull {
				obj := new(Forward)
				obj.load(v2, linkParent, obj, objThis, "ReverseNull")
				if linkParent && parentKey == "ForwardsAsNull" && obj.reverseNullID == objParent.(*Forward).reverseNullID {
					obj = objParent.(*Forward)
				}
				o.oForwardsAsNull = append(o.oForwardsAsNull, obj)
				o.mForwardsAsNull[obj.PrimaryKey()] = obj
			}
		case db.ValueMap: // single expansion
			obj := new(Forward)
			obj.load(oForwardsAsNull, linkParent, obj, objThis, "ReverseNull")
			if linkParent && parentKey == "ForwardsAsNull" && obj.reverseNullID == objParent.(*Forward).reverseNullID {
				obj = objParent.(*Forward)
			}
			o.oForwardsAsNull = []*Forward{obj}
		default:
			panic("Wrong type found for oForwardsAsNull object.")
		}
	} else {
		o.oForwardsAsNull = nil
	}

	if v, ok := m["ForwardAsUniqueNull"]; ok {
		if oForwardAsUniqueNull, ok2 := v.(db.ValueMap); ok2 {
			o.oForwardAsUniqueNull = new(Forward)
			o.oForwardAsUniqueNull.load(oForwardAsUniqueNull, linkParent, o.oForwardAsUniqueNull, objThis, "ReverseUniqueNull")
		} else {
			panic("Wrong type found for oForwardAsUniqueNull object.")
		}
	} else {
		o.oForwardAsUniqueNull = nil
	}

	if v, ok := m["aliases_"]; ok {
		o._aliases = map[string]interface{}(v.(db.ValueMap))
	}
	o._restored = true
}

// Save will update or insert the object, depending on the state of the object.
// If it has any auto-generated ids, those will be updated.
func (o *reverseBase) Save(ctx context.Context) {
	if o._restored {
		o.Update(ctx)
	} else {
		o.Insert(ctx)
	}
}

// Update will update the values in the database, saving any changed values.
func (o *reverseBase) Update(ctx context.Context) {
	if !o._restored {
		panic("Cannot update a record that was not originally read from the database.")
	}
	m := o.getModifiedFields()
	if len(m) == 0 {
		return
	}
	d := db.GetDatabase("goraddUnit")
	d.Update(ctx, "reverse", m, "id", fmt.Sprint(o.id))
	o.resetDirtyStatus()
	broadcast.Update(ctx, "goraddUnit", "reverse", o.id, stringmap.SortedKeys(m)...)
}

// Insert forces the object to be inserted into the database. If the object was loaded from the database originally,
// this will create a duplicate in the database.
func (o *reverseBase) Insert(ctx context.Context) {
	m := o.getModifiedFields()
	if len(m) == 0 {
		return
	}
	d := db.GetDatabase("goraddUnit")
	id := d.Insert(ctx, "reverse", m)
	o.id = id
	o.resetDirtyStatus()
	o._restored = true
	broadcast.Insert(ctx, "goraddUnit", "reverse", o.id)
}

func (o *reverseBase) getModifiedFields() (fields map[string]interface{}) {
	fields = map[string]interface{}{}
	if o.idIsDirty {
		fields["id"] = o.id
	}

	if o.nameIsDirty {
		fields["name"] = o.name
	}

	return
}

// Delete deletes the associated record from the database.
func (o *reverseBase) Delete(ctx context.Context) {
	if !o._restored {
		panic("Cannot delete a record that has no primary key value.")
	}
	d := db.GetDatabase("goraddUnit")
	d.Delete(ctx, "reverse", "id", o.id)
	broadcast.Delete(ctx, "goraddUnit", "reverse", o.id)
}

// deleteReverse deletes the associated record from the database.
func deleteReverse(ctx context.Context, pk string) {
	d := db.GetDatabase("goraddUnit")
	d.Delete(ctx, "reverse", "id", pk)
	broadcast.Delete(ctx, "goraddUnit", "reverse", pk)
}

func (o *reverseBase) resetDirtyStatus() {
	o.idIsDirty = false
	o.nameIsDirty = false
}

func (o *reverseBase) IsDirty() bool {
	return o.idIsDirty ||
		o.nameIsDirty
}

// Get returns the value of a field in the object based on the field's name.
// It will also get related objects if they are loaded.
// Invalid fields and objects are returned as nil
func (o *reverseBase) Get(key string) interface{} {

	switch key {
	case "ID":
		if !o.idIsValid {
			return nil
		}
		return o.id

	case "Name":
		if !o.nameIsValid {
			return nil
		}
		return o.name

	case "ForwardsAsNotNull":
		return o.ForwardsAsNotNull()

	case "ForwardAsUniqueNotNull":
		return o.ForwardAsUniqueNotNull()

	case "ForwardsAsNull":
		return o.ForwardsAsNull()

	case "ForwardAsUniqueNull":
		return o.ForwardAsUniqueNull()

	}
	return nil
}

// MarshalBinary serializes the object into a buffer that is deserializable using UnmarshalBinary.
// It should be used for transmitting database object over the wire, or for temporary storage. It does not send
// a version number, so if the data format changes, its up to you to invalidate the old stored objects.
// The framework uses this to serialize the object when it is stored in a control.
func (o *reverseBase) MarshalBinary() ([]byte, error) {
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

	if err := encoder.Encode(o.name); err != nil {
		return nil, err
	}
	if err := encoder.Encode(o.nameIsValid); err != nil {
		return nil, err
	}
	if err := encoder.Encode(o.nameIsDirty); err != nil {
		return nil, err
	}

	if err := encoder.Encode(o.oForwardsAsNotNull); err != nil {
		return nil, err
	}

	if err := encoder.Encode(o.oForwardAsUniqueNotNull); err != nil {
		return nil, err
	}
	if err := encoder.Encode(o.oForwardsAsNull); err != nil {
		return nil, err
	}

	if err := encoder.Encode(o.oForwardAsUniqueNull); err != nil {
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

	return buf.Bytes(), nil
}

func (o *reverseBase) UnmarshalBinary(data []byte) (err error) {

	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	if err = dec.Decode(&o.id); err != nil {
		return
	}
	if err = dec.Decode(&o.idIsValid); err != nil {
		return
	}
	if err = dec.Decode(&o.idIsDirty); err != nil {
		return
	}

	if err = dec.Decode(&o.name); err != nil {
		return
	}
	if err = dec.Decode(&o.nameIsValid); err != nil {
		return
	}
	if err = dec.Decode(&o.nameIsDirty); err != nil {
		return
	}

	if err = dec.Decode(&o.oForwardsAsNotNull); err != nil {
		return
	}
	if len(o.oForwardsAsNotNull) > 0 {
		o.mForwardsAsNotNull = make(map[string]*Forward)
		for _, p := range o.oForwardsAsNotNull {
			o.mForwardsAsNotNull[p.PrimaryKey()] = p
		}
	}
	if err = dec.Decode(&o.oForwardAsUniqueNotNull); err != nil {
		return
	}
	if err = dec.Decode(&o.oForwardsAsNull); err != nil {
		return
	}
	if len(o.oForwardsAsNull) > 0 {
		o.mForwardsAsNull = make(map[string]*Forward)
		for _, p := range o.oForwardsAsNull {
			o.mForwardsAsNull[p.PrimaryKey()] = p
		}
	}
	if err = dec.Decode(&o.oForwardAsUniqueNull); err != nil {
		return
	}

	var hasAliases bool
	if err = dec.Decode(&hasAliases); err != nil {
		return
	}
	if hasAliases {
		if err = dec.Decode(&o._aliases); err != nil {
			return
		}
	}

	if err = dec.Decode(&o._restored); err != nil {
		return
	}

	return
}

// MarshalJSON serializes the object into a JSON object.
// Only valid data will be serialized, meaning, you can control what gets serialized by using Select to
// select only the fields you want when you query for the object.
func (o *reverseBase) MarshalJSON() (data []byte, err error) {
	v := make(map[string]interface{})

	if o.idIsValid {
		v["id"] = o.id
	}

	if o.nameIsValid {
		v["name"] = o.name
	}

	if val := o.ForwardsAsNotNull(); val != nil {
		v["reverseNotNull"] = val
	}

	if val := o.ForwardAsUniqueNotNull(); val != nil {
		v["reverseUniqueNotNull"] = val
	}

	if val := o.ForwardsAsNull(); val != nil {
		v["reverseNull"] = val
	}

	if val := o.ForwardAsUniqueNull(); val != nil {
		v["reverseUniqueNull"] = val
	}

	return json.Marshal(v)
}
