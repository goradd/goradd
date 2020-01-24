package model

// Code generated by goradd. DO NOT EDIT.

import (
	"context"
	"fmt"

	"github.com/goradd/goradd/pkg/orm/broadcast"
	"github.com/goradd/goradd/pkg/orm/db"
	"github.com/goradd/goradd/pkg/orm/op"
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
	oForwardCascades              []*ForwardCascade          // Objects in the order they were queried
	mForwardCascades              map[string]*ForwardCascade // Objects by PK
	oForwardCascadesIsDirty       bool
	oForwardCascadeUnique         *ForwardCascadeUnique
	oForwardCascadeUniqueIsDirty  bool
	oForwardNulls                 []*ForwardNull          // Objects in the order they were queried
	mForwardNulls                 map[string]*ForwardNull // Objects by PK
	oForwardNullsIsDirty          bool
	oForwardNullUnique            *ForwardNullUnique
	oForwardNullUniqueIsDirty     bool
	oForwardRestricts             []*ForwardRestrict          // Objects in the order they were queried
	mForwardRestricts             map[string]*ForwardRestrict // Objects by PK
	oForwardRestrictsIsDirty      bool
	oForwardRestrictUnique        *ForwardRestrictUnique
	oForwardRestrictUniqueIsDirty bool

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
	Reverse_ID             = `ID`
	Reverse_Name           = `Name`
	ReverseForwardCascades = `ForwardCascades`

	ReverseForwardCascadeUnique = `ForwardCascadeUnique`
	ReverseForwardNulls         = `ForwardNulls`

	ReverseForwardNullUnique = `ForwardNullUnique`
	ReverseForwardRestricts  = `ForwardRestricts`

	ReverseForwardRestrictUnique = `ForwardRestrictUnique`
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
		panic("name was not selected in the last query and has not been set, and so is not valid")
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

// ForwardCascade returns a single ForwardCascade object by primary key, if one was loaded.
// Otherwise, it will return nil. It will not return ForwardCascade objects that are not saved.
func (o *reverseBase) ForwardCascade(pk string) *ForwardCascade {
	if o.mForwardCascades == nil {
		return nil
	}
	v, _ := o.mForwardCascades[pk]
	return v
}

// ForwardCascades returns a slice of ForwardCascade objects if loaded.
func (o *reverseBase) ForwardCascades() []*ForwardCascade {
	if o.oForwardCascades == nil {
		return nil
	}
	return o.oForwardCascades
}

// LoadForwardCascades loads a new slice of ForwardCascade objects and returns it.
func (o *reverseBase) LoadForwardCascades(ctx context.Context, conditions ...interface{}) []*ForwardCascade {
	qb := queryForwardCascades(ctx)
	cond := Equal(node.ForwardCascade().ReverseID(), o.PrimaryKey())
	if conditions != nil {
		conditions = append(conditions, cond)
		cond = And(conditions...)
	}

	o.oForwardCascades = qb.Where(cond).Load(ctx)
	return o.oForwardCascades
}

// SetForwardCascades associates the given objects with the Reverse.
// If it has items already associated with it that will not be associated after a save,
// the foreign keys for those will be set to null.
// If you did not use a join to query the items in the first place, used a conditional join,
// or joined with an expansion, be particularly careful, since you may be changing items
// that are not currently attached to this Reverse.
func (o *reverseBase) SetForwardCascades(objs []*ForwardCascade) {
	for _, obj := range o.oForwardCascades {
		if obj.IsDirty() {
			panic("You cannot overwrite items that have changed but have not been saved.")
		}
	}

	o.oForwardCascades = objs
	o.mForwardCascades = make(map[string]*ForwardCascade)
	for _, obj := range o.oForwardCascades {
		pk := obj.ID()
		if pk != "" {
			o.mForwardCascades[pk] = obj
		}
	}
	o.oForwardCascadesIsDirty = true
}

// ForwardCascadeUnique returns the connected ForwardCascadeUnique object, if one was loaded
// otherwise, it will return nil.
func (o *reverseBase) ForwardCascadeUnique() *ForwardCascadeUnique {
	if o.oForwardCascadeUnique == nil {
		return nil
	}
	return o.oForwardCascadeUnique
}

// LoadForwardCascadeUnique returns the connected ForwardCascadeUnique object, if one was loaded
// otherwise, it will return nil.
func (o *reverseBase) LoadForwardCascadeUnique(ctx context.Context) *ForwardCascadeUnique {
	if o.oForwardCascadeUnique == nil {
		o.oForwardCascadeUnique = LoadForwardCascadeUniqueByReverseID(ctx, o.ID())
	}
	return o.oForwardCascadeUnique
}

// SetForwardCascadeUnique associates the given object with the Reverse.
// If it has an item already associated with it,
// the foreign key for that item will be set to null.
// If you did not use a join to query the items in the first place, used a conditional join,
// or joined with an expansion, be particularly careful, since you may be changing an item
// that is not currently attached to this Reverse.
func (o *reverseBase) SetForwardCascadeUnique(obj *ForwardCascadeUnique) {
	if o.oForwardCascadeUnique != nil && o.oForwardCascadeUnique.IsDirty() {
		panic("The ForwardCascadeUnique has changed. You must save it first before changing to a different one.")
	}
	o.oForwardCascadeUnique = obj
	o.oForwardCascadeUniqueIsDirty = true
}

// ForwardNull returns a single ForwardNull object by primary key, if one was loaded.
// Otherwise, it will return nil. It will not return ForwardNull objects that are not saved.
func (o *reverseBase) ForwardNull(pk string) *ForwardNull {
	if o.mForwardNulls == nil {
		return nil
	}
	v, _ := o.mForwardNulls[pk]
	return v
}

// ForwardNulls returns a slice of ForwardNull objects if loaded.
func (o *reverseBase) ForwardNulls() []*ForwardNull {
	if o.oForwardNulls == nil {
		return nil
	}
	return o.oForwardNulls
}

// LoadForwardNulls loads a new slice of ForwardNull objects and returns it.
func (o *reverseBase) LoadForwardNulls(ctx context.Context, conditions ...interface{}) []*ForwardNull {
	qb := queryForwardNulls(ctx)
	cond := Equal(node.ForwardNull().ReverseID(), o.PrimaryKey())
	if conditions != nil {
		conditions = append(conditions, cond)
		cond = And(conditions...)
	}

	o.oForwardNulls = qb.Where(cond).Load(ctx)
	return o.oForwardNulls
}

// SetForwardNulls associates the given objects with the Reverse.
// If it has items already associated with it that will not be associated after a save,
// the foreign keys for those will be set to null.
// If you did not use a join to query the items in the first place, used a conditional join,
// or joined with an expansion, be particularly careful, since you may be changing items
// that are not currently attached to this Reverse.
func (o *reverseBase) SetForwardNulls(objs []*ForwardNull) {
	for _, obj := range o.oForwardNulls {
		if obj.IsDirty() {
			panic("You cannot overwrite items that have changed but have not been saved.")
		}
	}

	o.oForwardNulls = objs
	o.mForwardNulls = make(map[string]*ForwardNull)
	for _, obj := range o.oForwardNulls {
		pk := obj.ID()
		if pk != "" {
			o.mForwardNulls[pk] = obj
		}
	}
	o.oForwardNullsIsDirty = true
}

// ForwardNullUnique returns the connected ForwardNullUnique object, if one was loaded
// otherwise, it will return nil.
func (o *reverseBase) ForwardNullUnique() *ForwardNullUnique {
	if o.oForwardNullUnique == nil {
		return nil
	}
	return o.oForwardNullUnique
}

// LoadForwardNullUnique returns the connected ForwardNullUnique object, if one was loaded
// otherwise, it will return nil.
func (o *reverseBase) LoadForwardNullUnique(ctx context.Context) *ForwardNullUnique {
	if o.oForwardNullUnique == nil {
		o.oForwardNullUnique = LoadForwardNullUniqueByReverseID(ctx, o.ID())
	}
	return o.oForwardNullUnique
}

// SetForwardNullUnique associates the given object with the Reverse.
// If it has an item already associated with it,
// the foreign key for that item will be set to null.
// If you did not use a join to query the items in the first place, used a conditional join,
// or joined with an expansion, be particularly careful, since you may be changing an item
// that is not currently attached to this Reverse.
func (o *reverseBase) SetForwardNullUnique(obj *ForwardNullUnique) {
	if o.oForwardNullUnique != nil && o.oForwardNullUnique.IsDirty() {
		panic("The ForwardNullUnique has changed. You must save it first before changing to a different one.")
	}
	o.oForwardNullUnique = obj
	o.oForwardNullUniqueIsDirty = true
}

// ForwardRestrict returns a single ForwardRestrict object by primary key, if one was loaded.
// Otherwise, it will return nil. It will not return ForwardRestrict objects that are not saved.
func (o *reverseBase) ForwardRestrict(pk string) *ForwardRestrict {
	if o.mForwardRestricts == nil {
		return nil
	}
	v, _ := o.mForwardRestricts[pk]
	return v
}

// ForwardRestricts returns a slice of ForwardRestrict objects if loaded.
func (o *reverseBase) ForwardRestricts() []*ForwardRestrict {
	if o.oForwardRestricts == nil {
		return nil
	}
	return o.oForwardRestricts
}

// LoadForwardRestricts loads a new slice of ForwardRestrict objects and returns it.
func (o *reverseBase) LoadForwardRestricts(ctx context.Context, conditions ...interface{}) []*ForwardRestrict {
	qb := queryForwardRestricts(ctx)
	cond := Equal(node.ForwardRestrict().ReverseID(), o.PrimaryKey())
	if conditions != nil {
		conditions = append(conditions, cond)
		cond = And(conditions...)
	}

	o.oForwardRestricts = qb.Where(cond).Load(ctx)
	return o.oForwardRestricts
}

// SetForwardRestricts associates the given objects with the Reverse.
// WARNING! If it has items already associated with it that will not be associated after a save,
// those items will be DELETED since they cannot be null.
// If you did not use a join to query the items in the first place, used a conditional join,
// or joined with an expansion, be particularly careful, since you may be changing items
// that are not currently attached to this Reverse.
func (o *reverseBase) SetForwardRestricts(objs []*ForwardRestrict) {
	for _, obj := range o.oForwardRestricts {
		if obj.IsDirty() {
			panic("You cannot overwrite items that have changed but have not been saved.")
		}
	}

	o.oForwardRestricts = objs
	o.mForwardRestricts = make(map[string]*ForwardRestrict)
	for _, obj := range o.oForwardRestricts {
		pk := obj.ID()
		if pk != "" {
			o.mForwardRestricts[pk] = obj
		}
	}
	o.oForwardRestrictsIsDirty = true
}

// ForwardRestrictUnique returns the connected ForwardRestrictUnique object, if one was loaded
// otherwise, it will return nil.
func (o *reverseBase) ForwardRestrictUnique() *ForwardRestrictUnique {
	if o.oForwardRestrictUnique == nil {
		return nil
	}
	return o.oForwardRestrictUnique
}

// LoadForwardRestrictUnique returns the connected ForwardRestrictUnique object, if one was loaded
// otherwise, it will return nil.
func (o *reverseBase) LoadForwardRestrictUnique(ctx context.Context) *ForwardRestrictUnique {
	if o.oForwardRestrictUnique == nil {
		o.oForwardRestrictUnique = LoadForwardRestrictUniqueByReverseID(ctx, o.ID())
	}
	return o.oForwardRestrictUnique
}

// SetForwardRestrictUnique associates the given object with the Reverse.
// If it has an item already associated with it,
// the foreign key for that item will be set to null.
// If you did not use a join to query the items in the first place, used a conditional join,
// or joined with an expansion, be particularly careful, since you may be changing an item
// that is not currently attached to this Reverse.
func (o *reverseBase) SetForwardRestrictUnique(obj *ForwardRestrictUnique) {
	if o.oForwardRestrictUnique != nil && o.oForwardRestrictUnique.IsDirty() {
		panic("The ForwardRestrictUnique has changed. You must save it first before changing to a different one.")
	}
	o.oForwardRestrictUnique = obj
	o.oForwardRestrictUniqueIsDirty = true
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
		o.load(item, o, nil, "")
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
		o.load(item, o, nil, "")
		reverseSlice = append(reverseSlice, o)
	}
	return reverseSlice
}

// Get is a convenience method to return only the first item found in a query.
// The entire query is performed, so you should generally use this only if you know
// you are selecting on one or very few items.
func (b *ReversesBuilder) Get(ctx context.Context) *Reverse {
	results := b.Load(ctx)
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
// Care must be taken in the query, as Select clauses might not be honored if the child object has fields selected which the parent object does not have.
func (o *reverseBase) load(m map[string]interface{}, objThis *Reverse, objParent interface{}, parentKey string) {
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

	if v, ok := m["ForwardCascades"]; ok {
		switch oForwardCascades := v.(type) {
		case []db.ValueMap:
			o.oForwardCascades = make([]*ForwardCascade, 0, len(oForwardCascades))
			o.mForwardCascades = make(map[string]*ForwardCascade, len(oForwardCascades))
			for _, v2 := range oForwardCascades {
				obj := new(ForwardCascade)
				obj.load(v2, obj, objThis, "Reverse")
				o.oForwardCascades = append(o.oForwardCascades, obj)
				o.mForwardCascades[obj.PrimaryKey()] = obj
				o.oForwardCascadesIsDirty = false
			}
		case db.ValueMap: // single expansion
			obj := new(ForwardCascade)
			obj.load(oForwardCascades, obj, objThis, "Reverse")
			o.oForwardCascades = []*ForwardCascade{obj}
			o.oForwardCascadesIsDirty = false
		default:
			panic("Wrong type found for oForwardCascades object.")
		}
	} else {
		o.oForwardCascades = nil
		o.oForwardCascadesIsDirty = false
	}

	if v, ok := m["ForwardCascadeUnique"]; ok {
		if oForwardCascadeUnique, ok2 := v.(db.ValueMap); ok2 {
			o.oForwardCascadeUnique = new(ForwardCascadeUnique)
			o.oForwardCascadeUnique.load(oForwardCascadeUnique, o.oForwardCascadeUnique, objThis, "Reverse")
			o.oForwardCascadeUniqueIsDirty = false
		} else {
			panic("Wrong type found for oForwardCascadeUnique object.")
		}
	} else {
		o.oForwardCascadeUnique = nil
		o.oForwardCascadeUniqueIsDirty = false
	}

	if v, ok := m["ForwardNulls"]; ok {
		switch oForwardNulls := v.(type) {
		case []db.ValueMap:
			o.oForwardNulls = make([]*ForwardNull, 0, len(oForwardNulls))
			o.mForwardNulls = make(map[string]*ForwardNull, len(oForwardNulls))
			for _, v2 := range oForwardNulls {
				obj := new(ForwardNull)
				obj.load(v2, obj, objThis, "Reverse")
				o.oForwardNulls = append(o.oForwardNulls, obj)
				o.mForwardNulls[obj.PrimaryKey()] = obj
				o.oForwardNullsIsDirty = false
			}
		case db.ValueMap: // single expansion
			obj := new(ForwardNull)
			obj.load(oForwardNulls, obj, objThis, "Reverse")
			o.oForwardNulls = []*ForwardNull{obj}
			o.oForwardNullsIsDirty = false
		default:
			panic("Wrong type found for oForwardNulls object.")
		}
	} else {
		o.oForwardNulls = nil
		o.oForwardNullsIsDirty = false
	}

	if v, ok := m["ForwardNullUnique"]; ok {
		if oForwardNullUnique, ok2 := v.(db.ValueMap); ok2 {
			o.oForwardNullUnique = new(ForwardNullUnique)
			o.oForwardNullUnique.load(oForwardNullUnique, o.oForwardNullUnique, objThis, "Reverse")
			o.oForwardNullUniqueIsDirty = false
		} else {
			panic("Wrong type found for oForwardNullUnique object.")
		}
	} else {
		o.oForwardNullUnique = nil
		o.oForwardNullUniqueIsDirty = false
	}

	if v, ok := m["ForwardRestricts"]; ok {
		switch oForwardRestricts := v.(type) {
		case []db.ValueMap:
			o.oForwardRestricts = make([]*ForwardRestrict, 0, len(oForwardRestricts))
			o.mForwardRestricts = make(map[string]*ForwardRestrict, len(oForwardRestricts))
			for _, v2 := range oForwardRestricts {
				obj := new(ForwardRestrict)
				obj.load(v2, obj, objThis, "Reverse")
				o.oForwardRestricts = append(o.oForwardRestricts, obj)
				o.mForwardRestricts[obj.PrimaryKey()] = obj
				o.oForwardRestrictsIsDirty = false
			}
		case db.ValueMap: // single expansion
			obj := new(ForwardRestrict)
			obj.load(oForwardRestricts, obj, objThis, "Reverse")
			o.oForwardRestricts = []*ForwardRestrict{obj}
			o.oForwardRestrictsIsDirty = false
		default:
			panic("Wrong type found for oForwardRestricts object.")
		}
	} else {
		o.oForwardRestricts = nil
		o.oForwardRestrictsIsDirty = false
	}

	if v, ok := m["ForwardRestrictUnique"]; ok {
		if oForwardRestrictUnique, ok2 := v.(db.ValueMap); ok2 {
			o.oForwardRestrictUnique = new(ForwardRestrictUnique)
			o.oForwardRestrictUnique.load(oForwardRestrictUnique, o.oForwardRestrictUnique, objThis, "Reverse")
			o.oForwardRestrictUniqueIsDirty = false
		} else {
			panic("Wrong type found for oForwardRestrictUnique object.")
		}
	} else {
		o.oForwardRestrictUnique = nil
		o.oForwardRestrictUniqueIsDirty = false
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
		o.update(ctx)
	} else {
		o.insert(ctx)
	}
}

// update will update the values in the database, saving any changed values.
func (o *reverseBase) update(ctx context.Context) {
	var modifiedFields map[string]interface{}
	d := Database()
	db.ExecuteTransaction(ctx, d, func() {

		if !o._restored {
			panic("Cannot update a record that was not originally read from the database.")
		}

		modifiedFields = o.getModifiedFields()
		if len(modifiedFields) != 0 {
			d.Update(ctx, "reverse", modifiedFields, "id", fmt.Sprint(o.id))
		}

		if o.oForwardCascadesIsDirty {

			// Since the other side of the relationship cannot be null, the objects to be detached must be deleted
			// We take care to only delete objects that are not being reattached
			objs := QueryForwardCascades(ctx).
				Where(op.Equal(node.ForwardCascade().ReverseID(), o.PrimaryKey())).
				Load(ctx)
			// TODO: select only the required fields
			for _, obj := range objs {
				if _, ok := o.mForwardCascades[obj.PrimaryKey()]; !ok {
					// The old object is not in the group of new objects
					obj.Delete(ctx)
				}
			}
			for _, obj := range o.oForwardCascades {
				obj.SetReverseID(o.PrimaryKey())
				obj.Save(ctx)
			}
		} else {
			for _, obj := range o.oForwardCascades {
				obj.Save(ctx)
			}
		}
		if o.oForwardCascadeUniqueIsDirty {

			// Since the other side of the relationship cannot be null, the object to be detached must be deleted
			obj := QueryForwardCascadeUniques(ctx).
				Where(op.Equal(node.ForwardCascadeUnique().ReverseID(), o.PrimaryKey())).
				Get(ctx)
			if obj != nil && obj.PrimaryKey() != o.oForwardCascadeUnique.PrimaryKey() {
				obj.Delete(ctx)
			}
			o.oForwardCascadeUnique.SetReverseID(o.PrimaryKey())
			o.oForwardCascadeUnique.Save(ctx)
		} else {
			if o.oForwardCascadeUnique != nil {
				o.oForwardCascadeUnique.Save(ctx)
			}
		}
		if o.oForwardNullsIsDirty {
			objs := QueryForwardNulls(ctx).
				Where(op.Equal(node.ForwardNull().ReverseID(), o.PrimaryKey())).
				Load(ctx)
			// TODO:select only the required fields
			for _, obj := range objs {
				if _, ok := o.mForwardNulls[obj.PrimaryKey()]; !ok {
					// The old object is not in the group of new objects
					obj.SetReverseID(nil)
					obj.Save(ctx)
				}
			}
			for _, obj := range o.oForwardNulls {
				obj.SetReverseID(o.PrimaryKey())
				obj.Save(ctx)
			}

		} else {
			for _, obj := range o.oForwardNulls {
				obj.Save(ctx)
			}
		}
		if o.oForwardNullUniqueIsDirty {
			obj := QueryForwardNullUniques(ctx).
				Where(op.Equal(node.ForwardNullUnique().ReverseID(), o.PrimaryKey())).
				Get(ctx)
			if obj != nil && obj.PrimaryKey() != o.oForwardNullUnique.PrimaryKey() {
				obj.SetReverseID(nil)
				obj.Save(ctx)
			}
			o.oForwardNullUnique.SetReverseID(o.PrimaryKey())
			o.oForwardNullUnique.Save(ctx)
		} else {
			if o.oForwardNullUnique != nil {
				o.oForwardNullUnique.Save(ctx)
			}
		}
		if o.oForwardRestrictsIsDirty {

			// Since the other side of the relationship cannot be null, the objects to be detached must be deleted
			// We take care to only delete objects that are not being reattached
			objs := QueryForwardRestricts(ctx).
				Where(op.Equal(node.ForwardRestrict().ReverseID(), o.PrimaryKey())).
				Load(ctx)
			// TODO: select only the required fields
			for _, obj := range objs {
				if _, ok := o.mForwardRestricts[obj.PrimaryKey()]; !ok {
					// The old object is not in the group of new objects
					obj.Delete(ctx)
				}
			}
			for _, obj := range o.oForwardRestricts {
				obj.SetReverseID(o.PrimaryKey())
				obj.Save(ctx)
			}
		} else {
			for _, obj := range o.oForwardRestricts {
				obj.Save(ctx)
			}
		}
		if o.oForwardRestrictUniqueIsDirty {
			obj := QueryForwardRestrictUniques(ctx).
				Where(op.Equal(node.ForwardRestrictUnique().ReverseID(), o.PrimaryKey())).
				Get(ctx)
			if obj != nil && obj.PrimaryKey() != o.oForwardRestrictUnique.PrimaryKey() {
				obj.SetReverseID(nil)
				obj.Save(ctx)
			}
			o.oForwardRestrictUnique.SetReverseID(o.PrimaryKey())
			o.oForwardRestrictUnique.Save(ctx)
		} else {
			if o.oForwardRestrictUnique != nil {
				o.oForwardRestrictUnique.Save(ctx)
			}
		}

	}) // transaction
	o.resetDirtyStatus()
	if len(modifiedFields) != 0 {
		broadcast.Update(ctx, "goraddUnit", "reverse", fmt.Sprint(o.id), stringmap.SortedKeys(modifiedFields)...)
	}
}

// insert will insert the item into the database. Related items will be saved.
func (o *reverseBase) insert(ctx context.Context) {
	d := Database()
	db.ExecuteTransaction(ctx, d, func() {

		if !o.nameIsValid {
			panic("a value for Name is required, and there is no default value. Call SetName() before inserting the record.")
		}
		m := o.getValidFields()

		id := d.Insert(ctx, "reverse", m)
		o.id = id

		if o.oForwardCascades != nil {
			o.mForwardCascades = make(map[string]*ForwardCascade)
			for _, obj := range o.oForwardCascades {
				obj.SetReverseID(id)
				obj.Save(ctx)
				o.mForwardCascades[obj.PrimaryKey()] = obj
			}
		}

		if o.oForwardCascadeUnique != nil {
			o.oForwardCascadeUnique.SetReverseID(id)
			o.oForwardCascadeUnique.Save(ctx)
		}

		if o.oForwardNulls != nil {
			o.mForwardNulls = make(map[string]*ForwardNull)
			for _, obj := range o.oForwardNulls {
				obj.SetReverseID(id)
				obj.Save(ctx)
				o.mForwardNulls[obj.PrimaryKey()] = obj
			}
		}

		if o.oForwardNullUnique != nil {
			o.oForwardNullUnique.SetReverseID(id)
			o.oForwardNullUnique.Save(ctx)
		}

		if o.oForwardRestricts != nil {
			o.mForwardRestricts = make(map[string]*ForwardRestrict)
			for _, obj := range o.oForwardRestricts {
				obj.SetReverseID(id)
				obj.Save(ctx)
				o.mForwardRestricts[obj.PrimaryKey()] = obj
			}
		}

		if o.oForwardRestrictUnique != nil {
			o.oForwardRestrictUnique.SetReverseID(id)
			o.oForwardRestrictUnique.Save(ctx)
		}

	}) // transaction
	o.resetDirtyStatus()
	o._restored = true
	broadcast.Insert(ctx, "goraddUnit", "reverse", fmt.Sprint(o.id))
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

func (o *reverseBase) getValidFields() (fields map[string]interface{}) {
	fields = map[string]interface{}{}
	if o.nameIsValid {
		fields["name"] = o.name
	}
	return
}

// Delete deletes the associated record from the database.
func (o *reverseBase) Delete(ctx context.Context) {
	if !o._restored {
		panic("Cannot delete a record that has no primary key value.")
	}
	d := Database()
	db.ExecuteTransaction(ctx, d, func() {
		{
			objs := QueryForwardCascades(ctx).
				Where(op.Equal(node.ForwardCascade().ReverseID(), o.PrimaryKey())).
				Select(node.ForwardCascade().PrimaryKeyNode()).
				Load(ctx)
			for _, obj := range objs {
				obj.Delete(ctx)
			}
			o.oForwardCascades = nil
		}
		{
			obj := QueryForwardCascadeUniques(ctx).
				Where(op.Equal(node.ForwardCascadeUnique().ReverseID(), o.PrimaryKey())).
				Select(node.ForwardCascadeUnique().PrimaryKeyNode()).
				Get(ctx)
			if obj != nil {
				obj.Delete(ctx)
			}
			o.oForwardCascadeUnique = nil
		}
		{
			objs := QueryForwardNulls(ctx).
				Where(op.Equal(node.ForwardNull().ReverseID(), o.PrimaryKey())).
				Select(node.ForwardNull().PrimaryKeyNode()).
				Load(ctx)
			for _, obj := range objs {
				obj.SetReverseID(nil)
				obj.Save(ctx)
			}
			o.oForwardNulls = nil
		}
		{
			obj := QueryForwardNullUniques(ctx).
				Where(op.Equal(node.ForwardNullUnique().ReverseID(), o.PrimaryKey())).
				Select(node.ForwardNullUnique().PrimaryKeyNode()).
				Get(ctx)
			if obj != nil {
				obj.SetReverseID(nil)
				obj.Save(ctx)
			}
			o.oForwardNullUnique = nil
		}
		{
			c := QueryForwardRestricts(ctx).
				Where(op.Equal(node.ForwardRestrict().ReverseID(), o.PrimaryKey())).
				Count(ctx, false)
			if c > 0 {
				panic("Cannot delete a record that has restricted foreign keys pointing to it.")
			}
		}
		{
			c := QueryForwardRestrictUniques(ctx).
				Where(op.Equal(node.ForwardRestrictUnique().ReverseID(), o.PrimaryKey())).
				Count(ctx, false)
			if c > 0 {
				panic("Cannot delete a record that has a restricted foreign key pointing to it.")
			}
		}

		d.Delete(ctx, "reverse", "id", o.id)
	})
	broadcast.Delete(ctx, "goraddUnit", "reverse", fmt.Sprint(o.id))
}

// deleteReverse deletes the associated record from the database.
func deleteReverse(ctx context.Context, pk string) {
	if obj := LoadReverse(ctx, pk, node.Reverse().PrimaryKeyNode()); obj != nil {
		obj.Delete(ctx)
	}
}

func (o *reverseBase) resetDirtyStatus() {
	o.idIsDirty = false
	o.nameIsDirty = false
	o.oForwardCascadesIsDirty = false
	o.oForwardCascadeUniqueIsDirty = false
	o.oForwardNullsIsDirty = false
	o.oForwardNullUniqueIsDirty = false
	o.oForwardRestrictsIsDirty = false
	o.oForwardRestrictUniqueIsDirty = false

}

func (o *reverseBase) IsDirty() bool {
	return o.idIsDirty ||
		o.nameIsDirty ||
		o.oForwardCascadesIsDirty ||
		o.oForwardCascadeUniqueIsDirty ||
		o.oForwardNullsIsDirty ||
		o.oForwardNullUniqueIsDirty ||
		o.oForwardRestrictsIsDirty ||
		o.oForwardRestrictUniqueIsDirty
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

	case "ForwardCascades":
		return o.ForwardCascades()

	case "ForwardCascadeUnique":
		return o.ForwardCascadeUnique()

	case "ForwardNulls":
		return o.ForwardNulls()

	case "ForwardNullUnique":
		return o.ForwardNullUnique()

	case "ForwardRestricts":
		return o.ForwardRestricts()

	case "ForwardRestrictUnique":
		return o.ForwardRestrictUnique()

	}
	return nil
}

// MarshalBinary serializes the object into a buffer that is deserializable using UnmarshalBinary.
// It should be used for transmitting database objects over the wire, or for temporary storage. It does not send
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

	if err := encoder.Encode(o.oForwardCascades); err != nil {
		return nil, err
	}

	if err := encoder.Encode(o.oForwardCascadeUnique); err != nil {
		return nil, err
	}
	if err := encoder.Encode(o.oForwardNulls); err != nil {
		return nil, err
	}

	if err := encoder.Encode(o.oForwardNullUnique); err != nil {
		return nil, err
	}
	if err := encoder.Encode(o.oForwardRestricts); err != nil {
		return nil, err
	}

	if err := encoder.Encode(o.oForwardRestrictUnique); err != nil {
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

	if err = dec.Decode(&o.oForwardCascades); err != nil {
		return
	}
	if len(o.oForwardCascades) > 0 {
		o.mForwardCascades = make(map[string]*ForwardCascade)
		for _, p := range o.oForwardCascades {
			o.mForwardCascades[p.PrimaryKey()] = p
		}
	}
	if err = dec.Decode(&o.oForwardCascadeUnique); err != nil {
		return
	}
	if err = dec.Decode(&o.oForwardNulls); err != nil {
		return
	}
	if len(o.oForwardNulls) > 0 {
		o.mForwardNulls = make(map[string]*ForwardNull)
		for _, p := range o.oForwardNulls {
			o.mForwardNulls[p.PrimaryKey()] = p
		}
	}
	if err = dec.Decode(&o.oForwardNullUnique); err != nil {
		return
	}
	if err = dec.Decode(&o.oForwardRestricts); err != nil {
		return
	}
	if len(o.oForwardRestricts) > 0 {
		o.mForwardRestricts = make(map[string]*ForwardRestrict)
		for _, p := range o.oForwardRestricts {
			o.mForwardRestricts[p.PrimaryKey()] = p
		}
	}
	if err = dec.Decode(&o.oForwardRestrictUnique); err != nil {
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

	if val := o.ForwardCascades(); val != nil {
		v["reverse"] = val
	}

	if val := o.ForwardCascadeUnique(); val != nil {
		v["reverse"] = val
	}

	if val := o.ForwardNulls(); val != nil {
		v["reverse"] = val
	}

	if val := o.ForwardNullUnique(); val != nil {
		v["reverse"] = val
	}

	if val := o.ForwardRestricts(); val != nil {
		v["reverse"] = val
	}

	if val := o.ForwardRestrictUnique(); val != nil {
		v["reverse"] = val
	}

	return json.Marshal(v)
}