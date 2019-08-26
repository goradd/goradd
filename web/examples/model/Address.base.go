package model

// Code generated by goradd. DO NOT EDIT.

import (
	"context"
	"fmt"
	"github.com/goradd/goradd/web/examples/model/node"

	"github.com/goradd/goradd/pkg/orm/db"
	. "github.com/goradd/goradd/pkg/orm/op"
	"github.com/goradd/goradd/pkg/orm/query"

	//"./node"
	"bytes"
	"encoding/gob"
	"encoding/json"
)

// addressBase is a base structure to be embedded in a "subclass" and provides the ORM access to the database.
// Do not directly access the internal variables, but rather use the accessor functions, since this class maintains internal state
// related to the variables.

type addressBase struct {
	id        string
	idIsValid bool
	idIsDirty bool

	personID        string
	personIDIsNull  bool
	personIDIsValid bool
	personIDIsDirty bool
	oPerson         *Person

	street        string
	streetIsValid bool
	streetIsDirty bool

	city        string
	cityIsNull  bool
	cityIsValid bool
	cityIsDirty bool

	// Custom aliases, if specified
	_aliases map[string]interface{}

	// Indicates whether this is a new object, or one loaded from the database. Used by Save to know whether to Insert or Update
	_restored bool
}

const (
	AddressIDDefault       = ""
	AddressPersonIDDefault = ""
	AddressStreetDefault   = ""
	AddressCityDefault     = ""
)

const (
	AddressID       = `ID`
	AddressPersonID = `PersonID`
	AddressPerson   = `Person`
	AddressStreet   = `Street`
	AddressCity     = `City`
)

// Initialize or re-initialize a Address database object to default values.
func (o *addressBase) Initialize() {

	o.id = ""
	o.idIsValid = false
	o.idIsDirty = false

	o.personID = ""
	o.personIDIsNull = true
	o.personIDIsValid = true
	o.personIDIsDirty = true

	o.street = ""
	o.streetIsValid = false
	o.streetIsDirty = false

	o.city = ""
	o.cityIsNull = true
	o.cityIsValid = true
	o.cityIsDirty = true

	o._restored = false
}

func (o *addressBase) PrimaryKey() string {
	return o.id
}

// ID returns the loaded value of ID.
func (o *addressBase) ID() string {
	return fmt.Sprint(o.id)
}

// IDIsValid returns true if the value was loaded from the database or has been set.
func (o *addressBase) IDIsValid() bool {
	return o._restored && o.idIsValid
}

func (o *addressBase) PersonID() string {
	if o._restored && !o.personIDIsValid {
		panic("personID was not selected in the last query and so is not valid")
	}
	return o.personID
}

// PersonIDIsValid returns true if the value was loaded from the database or has been set.
func (o *addressBase) PersonIDIsValid() bool {
	return o.personIDIsValid
}

// PersonIDIsNull returns true if the related database value is null.
func (o *addressBase) PersonIDIsNull() bool {
	return o.personIDIsNull
}

// Person returns the current value of the loaded Person, and nil if its not loaded.
func (o *addressBase) Person() *Person {
	return o.oPerson
}

// LoadPerson returns the related Person. If it is not already loaded,
// it will attempt to load it first.
func (o *addressBase) LoadPerson(ctx context.Context) *Person {
	if !o.personIDIsValid {
		return nil
	}

	if o.oPerson == nil {
		// Load and cache
		o.oPerson = LoadPerson(ctx, o.PersonID())
	}
	return o.oPerson
}

func (o *addressBase) SetPersonID(i interface{}) {
	o.personIDIsValid = true
	if i == nil {
		if !o.personIDIsNull {
			o.personIDIsNull = true
			o.personIDIsDirty = true
			o.personID = ""
			o.oPerson = nil
		}
	} else {
		v := i.(string)
		if o.personIDIsNull ||
			!o._restored ||
			o.personID != v {

			o.personIDIsNull = false
			o.personID = v
			o.personIDIsDirty = true
			o.oPerson = nil
		}
	}
}

func (o *addressBase) SetPerson(v *Person) {
	o.personIDIsValid = true
	if v == nil {
		if !o.personIDIsNull || !o._restored {
			o.personIDIsNull = true
			o.personIDIsDirty = true
			o.personID = ""
			o.oPerson = nil
		}
	} else {
		o.oPerson = v
		if o.personIDIsNull || !o._restored || o.personID != v.PrimaryKey() {
			o.personIDIsNull = false
			o.personID = v.PrimaryKey()
			o.personIDIsDirty = true
		}
	}
}

func (o *addressBase) Street() string {
	if o._restored && !o.streetIsValid {
		panic("street was not selected in the last query and so is not valid")
	}
	return o.street
}

// StreetIsValid returns true if the value was loaded from the database or has been set.
func (o *addressBase) StreetIsValid() bool {
	return o.streetIsValid
}

// SetStreet sets the value of Street in the object, to be saved later using the Save() function.
func (o *addressBase) SetStreet(v string) {
	o.streetIsValid = true
	if o.street != v || !o._restored {
		o.street = v
		o.streetIsDirty = true
	}

}

func (o *addressBase) City() string {
	if o._restored && !o.cityIsValid {
		panic("city was not selected in the last query and so is not valid")
	}
	return o.city
}

// CityIsValid returns true if the value was loaded from the database or has been set.
func (o *addressBase) CityIsValid() bool {
	return o.cityIsValid
}

// CityIsNull returns true if the related database value is null.
func (o *addressBase) CityIsNull() bool {
	return o.cityIsNull
}

func (o *addressBase) SetCity(i interface{}) {
	o.cityIsValid = true
	if i == nil {
		if !o.cityIsNull {
			o.cityIsNull = true
			o.cityIsDirty = true
			o.city = ""
		}
	} else {
		v := i.(string)
		if o.cityIsNull ||
			!o._restored ||
			o.city != v {

			o.cityIsNull = false
			o.city = v
			o.cityIsDirty = true
		}
	}
}

// GetAlias returns the alias for the given key.
func (o *addressBase) GetAlias(key string) query.AliasValue {
	if a, ok := o._aliases[key]; ok {
		return query.NewAliasValue(a)
	} else {
		panic("Alias " + key + " not found.")
		return query.NewAliasValue([]byte{})
	}
}

// Load returns a Address from the database.
// joinOrSelectNodes lets you provide nodes for joining to other tables or selecting specific fields. Table nodes will
// be considered Join nodes, and column nodes will be Select nodes. See Join() and Select() for more info.
func LoadAddress(ctx context.Context, primaryKey string, joinOrSelectNodes ...query.NodeI) *Address {
	return queryAddresses(ctx).Where(Equal(node.Address().ID(), primaryKey)).joinOrSelect(joinOrSelectNodes...).Get(ctx)
}

// The AddressesBuilder uses the QueryBuilderI interface from the database to build a query.
// All query operations go through this query builder.
// End a query by calling either Load, Count, or Delete
type AddressesBuilder struct {
	base                query.QueryBuilderI
	hasConditionalJoins bool
}

func newAddressBuilder() *AddressesBuilder {
	b := &AddressesBuilder{
		base: db.GetDatabase("goradd").
			NewBuilder(),
	}
	return b.Join(node.Address())
}

// Load terminates the query builder, performs the query, and returns a slice of Address objects. If there are
// any errors, they are returned in the context object. If no results come back from the query, it will return
// an empty slice
func (b *AddressesBuilder) Load(ctx context.Context) (addressSlice []*Address) {
	results := b.base.Load(ctx)
	if results == nil {
		return
	}
	for _, item := range results {
		o := new(Address)
		o.load(item, !b.hasConditionalJoins, o, nil, "")
		addressSlice = append(addressSlice, o)
	}
	return addressSlice
}

// LoadI terminates the query builder, performs the query, and returns a slice of interfaces. If there are
// any errors, they are returned in the context object. If no results come back from the query, it will return
// an empty slice.
func (b *AddressesBuilder) LoadI(ctx context.Context) (addressSlice []interface{}) {
	results := b.base.Load(ctx)
	if results == nil {
		return
	}
	for _, item := range results {
		o := new(Address)
		o.load(item, !b.hasConditionalJoins, o, nil, "")
		addressSlice = append(addressSlice, o)
	}
	return addressSlice
}

// Get is a convenience method to return only the first item found in a query. It is equivalent to adding
// Limit(1,0) to the query, and then getting the first item from the returned slice.
// Limits with joins do not currently work, so don't try it if you have a join
// TODO: Change this to Load1 to be more descriptive and avoid confusion with other Getters
func (b *AddressesBuilder) Get(ctx context.Context) *Address {
	results := b.Limit(1, 0).Load(ctx)
	if results != nil && len(results) > 0 {
		obj := results[0]
		return obj
	} else {
		return nil
	}
}

// Expand expands an array type node so that it will produce individual rows instead of an array of items
func (b *AddressesBuilder) Expand(n query.NodeI) *AddressesBuilder {
	b.base.Expand(n)
	return b
}

// Join adds a node to the node tree so that its fields will appear in the query. Optionally add conditions to filter
// what gets included. The conditions will be AND'd with the basic condition matching the primary keys of the join.
func (b *AddressesBuilder) Join(n query.NodeI, conditions ...query.NodeI) *AddressesBuilder {
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
func (b *AddressesBuilder) Where(c query.NodeI) *AddressesBuilder {
	b.base.Condition(c)
	return b
}

// OrderBy specifies how the resulting data should be sorted.
func (b *AddressesBuilder) OrderBy(nodes ...query.NodeI) *AddressesBuilder {
	b.base.OrderBy(nodes...)
	return b
}

// Limit will return a subset of the data, limited to the offset and number of rows specified
func (b *AddressesBuilder) Limit(maxRowCount int, offset int) *AddressesBuilder {
	b.base.Limit(maxRowCount, offset)
	return b
}

// Select optimizes the query to only return the specified fields. Once you put a Select in your query, you must
// specify all the fields that you will eventually read out. Be careful when selecting fields in joined tables, as joined
// tables will also contain pointers back to the parent table, and so the parent node should have the same field selected
// as the child node if you are querying those fields.
func (b *AddressesBuilder) Select(nodes ...query.NodeI) *AddressesBuilder {
	b.base.Select(nodes...)
	return b
}

// Alias lets you add a node with a custom name. After the query, you can read out the data using GetAlias() on a
// returned object. Alias is useful for adding calculations or subqueries to the query.
func (b *AddressesBuilder) Alias(name string, n query.NodeI) *AddressesBuilder {
	b.base.Alias(name, n)
	return b
}

// Distinct removes duplicates from the results of the query. Adding a Select() may help you get to the data you want, although
// using Distinct with joined tables is often not effective, since we force joined tables to include primary keys in the query, and this
// often ruins the effect of Distinct.
func (b *AddressesBuilder) Distinct() *AddressesBuilder {
	b.base.Distinct()
	return b
}

// GroupBy controls how results are grouped when using aggregate functions in an Alias() call.
func (b *AddressesBuilder) GroupBy(nodes ...query.NodeI) *AddressesBuilder {
	b.base.GroupBy(nodes...)
	return b
}

// Having does additional filtering on the results of the query.
func (b *AddressesBuilder) Having(node query.NodeI) *AddressesBuilder {
	b.base.Having(node)
	return b
}

// Count terminates a query and returns just the number of items selected.
func (b *AddressesBuilder) Count(ctx context.Context, distinct bool, nodes ...query.NodeI) uint {
	return b.base.Count(ctx, distinct, nodes...)
}

// Delete uses the query builder to delete a group of records that match the criteria
func (b *AddressesBuilder) Delete(ctx context.Context) {
	b.base.Delete(ctx)
}

// Subquery uses the query builder to define a subquery within a larger query. You MUST include what
// you are selecting by adding Alias or Select functions on the subquery builder. Generally you would use
// this as a node to an Alias function on the surrounding query builder.
func (b *AddressesBuilder) Subquery() *query.SubqueryNode {
	return b.base.Subquery()
}

// joinOrSelect us a private helper function for the Load* functions
func (b *AddressesBuilder) joinOrSelect(nodes ...query.NodeI) *AddressesBuilder {
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

func CountAddressByID(ctx context.Context, id string) uint {
	return queryAddresses(ctx).Where(Equal(node.Address().ID(), id)).Count(ctx, false)
}

func CountAddressByPersonID(ctx context.Context, personID string) uint {
	return queryAddresses(ctx).Where(Equal(node.Address().PersonID(), personID)).Count(ctx, false)
}

func CountAddressByStreet(ctx context.Context, street string) uint {
	return queryAddresses(ctx).Where(Equal(node.Address().Street(), street)).Count(ctx, false)
}

func CountAddressByCity(ctx context.Context, city string) uint {
	return queryAddresses(ctx).Where(Equal(node.Address().City(), city)).Count(ctx, false)
}

// load is the private loader that transforms data coming from the database into a tree structure reflecting the relationships
// between the object chain requested by the user in the query.
// If linkParent is true we will have child relationships use a pointer back to the parent object. If false, it will create a separate object.
// Care must be taken in the query, as Select clauses might not be honored if the child object has fields selected which the parent object does not have.
// Also, if any joins are conditional, that might affect which child objects are included, so in this situation, linkParent should be false
func (o *addressBase) load(m map[string]interface{}, linkParent bool, objThis *Address, objParent interface{}, parentKey string) {
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

	if v, ok := m["person_id"]; ok {
		if v == nil {
			o.personID = ""
			o.personIDIsNull = true
			o.personIDIsValid = true
			o.personIDIsDirty = false
		} else if o.personID, ok = v.(string); ok {
			o.personIDIsNull = false
			o.personIDIsValid = true
			o.personIDIsDirty = false
		} else {
			panic("Wrong type found for person_id.")
		}
	} else {
		o.personIDIsValid = false
		o.personIDIsNull = true
		o.personID = ""
	}
	if linkParent && parentKey == "Person" {
		o.oPerson = objParent.(*Person)
		o.personIDIsValid = true
		o.personIDIsDirty = false
	} else if v, ok := m["Person"]; ok {
		if oPerson, ok2 := v.(map[string]interface{}); ok2 {
			o.oPerson = new(Person)
			o.oPerson.load(oPerson, linkParent, o.oPerson, objThis, "Addresses")
			o.personIDIsValid = true
			o.personIDIsDirty = false
		} else {
			panic("Wrong type found for oPerson object.")
		}
	} else {
		o.oPerson = nil
	}

	if v, ok := m["street"]; ok && v != nil {
		if o.street, ok = v.(string); ok {
			o.streetIsValid = true
			o.streetIsDirty = false
		} else {
			panic("Wrong type found for street.")
		}
	} else {
		o.streetIsValid = false
		o.street = ""
	}

	if v, ok := m["city"]; ok {
		if v == nil {
			o.city = ""
			o.cityIsNull = true
			o.cityIsValid = true
			o.cityIsDirty = false
		} else if o.city, ok = v.(string); ok {
			o.cityIsNull = false
			o.cityIsValid = true
			o.cityIsDirty = false
		} else {
			panic("Wrong type found for city.")
		}
	} else {
		o.cityIsValid = false
		o.cityIsNull = true
		o.city = ""
	}

	if v, ok := m["aliases_"]; ok {
		o._aliases = map[string]interface{}(v.(db.ValueMap))
	}
	o._restored = true
}

// Save will update or insert the object, depending on the state of the object.
// If it has any auto-generated ids, those will be updated.
func (o *addressBase) Save(ctx context.Context) {
	if o._restored {
		o.Update(ctx)
	} else {
		o.Insert(ctx)
	}
}

// Update will update the values in the database, saving any changed values.
func (o *addressBase) Update(ctx context.Context) {
	if !o._restored {
		panic("Cannot update a record that was not originally read from the database.")
	}
	m := o.getModifiedFields()
	if len(m) == 0 {
		return
	}
	d := db.GetDatabase("goradd")
	d.Update(ctx, "address", m, "id", fmt.Sprint(o.id))
	o.resetDirtyStatus()
}

// Insert forces the object to be inserted into the database. If the object was loaded from the database originally,
// this will create a duplicate in the database.
func (o *addressBase) Insert(ctx context.Context) {
	m := o.getModifiedFields()
	if len(m) == 0 {
		return
	}
	d := db.GetDatabase("goradd")
	id := d.Insert(ctx, "address", m)
	o.id = id
	o.resetDirtyStatus()
	o._restored = true
}

func (o *addressBase) getModifiedFields() (fields map[string]interface{}) {
	fields = map[string]interface{}{}
	if o.idIsDirty {
		fields["id"] = o.id
	}

	if o.personIDIsDirty {
		if o.personIDIsNull {
			fields["person_id"] = nil
		} else {
			fields["person_id"] = o.personID
		}
	}

	if o.streetIsDirty {
		fields["street"] = o.street
	}

	if o.cityIsDirty {
		if o.cityIsNull {
			fields["city"] = nil
		} else {
			fields["city"] = o.city
		}
	}

	return
}

// Delete deletes the associated record from the database.
func (o *addressBase) Delete(ctx context.Context) {
	if !o._restored {
		panic("Cannot delete a record that has no primary key value.")
	}
	d := db.GetDatabase("goradd")
	d.Delete(ctx, "address", "id", o.id)
}

// deleteAddress deletes the associated record from the database.
func deleteAddress(ctx context.Context, pk string) {
	d := db.GetDatabase("goradd")
	d.Delete(ctx, "address", "id", pk)
}

func (o *addressBase) resetDirtyStatus() {
	o.idIsDirty = false
	o.personIDIsDirty = false
	o.streetIsDirty = false
	o.cityIsDirty = false
}

func (o *addressBase) IsDirty() bool {
	return o.idIsDirty ||
		o.personIDIsDirty ||
		o.streetIsDirty ||
		o.cityIsDirty
}

// Get returns the value of a field in the object based on the field's name.
// It will also get related objects if they are loaded.
// Invalid fields and objects are returned as nil
func (o *addressBase) Get(key string) interface{} {

	switch key {
	case "ID":
		if !o.idIsValid {
			return nil
		}
		return o.id

	case "PersonID":
		if !o.personIDIsValid {
			return nil
		}
		return o.personID

	case "Person":
		return o.Person()

	case "Street":
		if !o.streetIsValid {
			return nil
		}
		return o.street

	case "City":
		if !o.cityIsValid {
			return nil
		}
		return o.city

	}
	return nil
}

// MarshalBinary serializes the object into a buffer that is deserializable using UnmarshalBinary.
// It should be used for transmitting database object over the wire, or for temporary storage. It does not send
// a version number, so if the data format changes, its up to you to invalidate the old stored objects.
// The framework uses this to serialize the object when it is stored in a control.
func (o *addressBase) MarshalBinary() (data []byte, err error) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)

	if err = encoder.Encode(o.id); err != nil {
		return
	}
	if err = encoder.Encode(o.idIsValid); err != nil {
		return
	}
	if err = encoder.Encode(o.idIsDirty); err != nil {
		return
	}

	if err = encoder.Encode(o.personID); err != nil {
		return
	}
	if err = encoder.Encode(o.personIDIsNull); err != nil {
		return
	}
	if err = encoder.Encode(o.personIDIsValid); err != nil {
		return
	}
	if err = encoder.Encode(o.personIDIsDirty); err != nil {
		return
	}

	if err = encoder.Encode(o.oPerson); err != nil {
		return
	}
	if err = encoder.Encode(o.street); err != nil {
		return
	}
	if err = encoder.Encode(o.streetIsValid); err != nil {
		return
	}
	if err = encoder.Encode(o.streetIsDirty); err != nil {
		return
	}

	if err = encoder.Encode(o.city); err != nil {
		return
	}
	if err = encoder.Encode(o.cityIsNull); err != nil {
		return
	}
	if err = encoder.Encode(o.cityIsValid); err != nil {
		return
	}
	if err = encoder.Encode(o.cityIsDirty); err != nil {
		return
	}

	if o._aliases == nil {
		if err = encoder.Encode(false); err != nil {
			return
		}
	} else {
		if err = encoder.Encode(true); err != nil {
			return
		}
		if err = encoder.Encode(o._aliases); err != nil {
			return
		}
	}

	if err = encoder.Encode(o._restored); err != nil {
		return
	}

	return
}

func (o *addressBase) UnmarshalBinary(data []byte) (err error) {

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

	if err = dec.Decode(&o.personID); err != nil {
		return
	}
	if err = dec.Decode(&o.personIDIsNull); err != nil {
		return
	}
	if err = dec.Decode(&o.personIDIsValid); err != nil {
		return
	}
	if err = dec.Decode(&o.personIDIsDirty); err != nil {
		return
	}

	if err = dec.Decode(&o.oPerson); err != nil {
		return
	}
	if err = dec.Decode(&o.street); err != nil {
		return
	}
	if err = dec.Decode(&o.streetIsValid); err != nil {
		return
	}
	if err = dec.Decode(&o.streetIsDirty); err != nil {
		return
	}

	if err = dec.Decode(&o.city); err != nil {
		return
	}
	if err = dec.Decode(&o.cityIsNull); err != nil {
		return
	}
	if err = dec.Decode(&o.cityIsValid); err != nil {
		return
	}
	if err = dec.Decode(&o.cityIsDirty); err != nil {
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

	return err
}

// MarshalJSON serializes the object into a JSON object.
// Only valid data will be serialized, meaning, you can control what gets serialized by using Select to
// select only the fields you want when you query for the object.
func (o *addressBase) MarshalJSON() (data []byte, err error) {
	v := make(map[string]interface{})

	if o.idIsValid {
		v["id"] = o.id
	}

	if o.personIDIsValid {
		if o.personIDIsNull {
			v["personID"] = nil
		} else {
			v["personID"] = o.personID
		}
	}

	if val := o.Person(); val != nil {
		v["person"] = val
	}
	if o.streetIsValid {
		v["street"] = o.street
	}

	if o.cityIsValid {
		if o.cityIsNull {
			v["city"] = nil
		} else {
			v["city"] = o.city
		}
	}

	return json.Marshal(v)
}
