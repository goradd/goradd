package sql

import (
	"context"
	"github.com/goradd/goradd/pkg/any"
)

// Associate is a helper function for the sql database implementations.
// It sets up the many-many association pointing from the given table and column to another table and column.
// table is the name of the association table.
// column is the name of the column in the association table that contains the pk for the record we are associating.
// pk is the value of the primary key.
// relatedTable is the table the association is pointing to.
// relatedColumn is the column in the association table that points to the relatedTable's pk.
// relatedPks are the new primary keys in the relatedTable we are associating.
func Associate(ctx context.Context,
	db DbI,
	table string,
	column string,
	pk interface{},
	relatedColumn string,
	relatedPks interface{}) { //relatedPks must be a slice of items

	// TODO: Could optimize by separating out what gets deleted, what gets added, and what stays the same.

	// TODO: Make this part of a transaction
	// First delete all previous associations
	var sql = "DELETE FROM " + db.QuoteIdentifier(table) + " WHERE " +
		db.QuoteIdentifier(column) + "=" + db.FormatArgument(1)
	_, e := db.Exec(ctx, sql, pk)
	if e != nil {
		panic(e.Error())
	}
	if relatedPks == nil {
		return
	}

	// Add new associations
	for _, relatedPk := range any.InterfaceSlice(relatedPks) {
		sql = "INSERT INTO " + db.QuoteIdentifier(table) + "(" +
			db.QuoteIdentifier(column) + "," + db.QuoteIdentifier(relatedColumn) +
			") VALUES (" + db.FormatArgument(1) + "," + db.FormatArgument(2) + ")"
		_, e = db.Exec(ctx, sql, pk, relatedPk)
		if e != nil {
			panic(e.Error())
		}
	}
}
