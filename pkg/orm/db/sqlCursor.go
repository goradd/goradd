package db

import (
	"database/sql"
	"github.com/goradd/goradd/pkg/orm/query"
	"log"
)

type sqlCursor struct {
	rows *sql.Rows
	columnTypes []query.GoColumnType
	columnNames []string
	builder *sqlBuilder
	columnReceivers []SqlReceiver
	columnValueReceivers []interface{}
}

func NewSqlCursor(rows *sql.Rows,
	columnTypes []query.GoColumnType,
	columnNames []string,
	builder *sqlBuilder,
	) *sqlCursor {
	var err error

	if columnNames == nil {
		columnNames, err = rows.Columns()
		if err != nil {
			log.Panic(err)
		}
	}

	cursor := sqlCursor{
		rows:                 rows,
		columnTypes:          columnTypes,
		columnNames:          columnNames,
		builder:			  builder,
		columnReceivers:      make([]SqlReceiver, len(columnTypes)),
		columnValueReceivers: make([]interface{}, len(columnTypes)),
	}

	for i := range cursor.columnReceivers {
		cursor.columnValueReceivers[i] = &(cursor.columnReceivers[i].R)
	}
	return &cursor
}

// Next returns the values of the next row in the result set.
//
// Returns nil if there are no more rows in the result set.
//
// If an error occurs, will panic with the error.
func (r sqlCursor) Next() map[string]interface{} {
	var err error

	if r.rows.Next() {
		if err = r.rows.Scan(r.columnValueReceivers...); err != nil {
			log.Panic(err)
		}

		values := make(map[string]interface{}, len(r.columnReceivers))
		for j, vr := range r.columnReceivers {
			values[r.columnNames[j]] = vr.Unpack(r.columnTypes[j])
		}
		if r.builder != nil {
			v2 := r.builder.unpackResult([]map[string]interface{}{values})
			return v2[0]
		} else {
			return values
		}
	} else {
		if err = r.rows.Err(); err != nil {
			log.Panic(err)
		}
		return nil
	}
}

// Close closes the cursor.
//
// Once you are done with the cursor, you MUST call Close, so its
// probably best to put a defer Close statement ahead of using Next.
func (r sqlCursor) Close() error {
	return r.rows.Close()
}

// sqlReceiveRows gets data from a sql result set and returns it as a slice of maps.
//
// Each column is mapped to its column name.
// If you provide columnNames, those will be used in the map. Otherwise it will get the column names out of the
// result set provided.
func sqlReceiveRows(rows *sql.Rows,
	columnTypes []query.GoColumnType,
	columnNames []string,
	builder *sqlBuilder,
	) []map[string]interface{} {

	var values []map[string]interface{}

	cursor := NewSqlCursor(rows, columnTypes, columnNames, nil)
	defer cursor.Close()
	for v := cursor.Next();v != nil;v = cursor.Next() {
		values = append(values, v)
	}
	if builder != nil {
		values = builder.unpackResult(values)
	}

	return values
}


// ReceiveRows gets data from a sql result set and returns it as a slice of maps. Each column is mapped to its column name.
// If you provide column names, those will be used in the map. Otherwise it will get the column names out of the
// result set provided
func sqlReceiveRows2(rows *sql.Rows, columnTypes []query.GoColumnType, columnNames []string) (values []map[string]interface{}) {
	var err error

	values = []map[string]interface{}{}

	columnReceivers := make([]SqlReceiver, len(columnTypes))
	columnValueReceivers := make([]interface{}, len(columnTypes))

	if columnNames == nil {
		columnNames, err = rows.Columns()
		if err != nil {
			log.Panic(err)
		}
	}

	for i, _ := range columnReceivers {
		columnValueReceivers[i] = &(columnReceivers[i].R)
	}

	for rows.Next() {
		err = rows.Scan(columnValueReceivers...)

		if err != nil {
			log.Panic(err)
		}

		v1 := make(map[string]interface{}, len(columnReceivers))
		for j, vr := range columnReceivers {
			v1[columnNames[j]] = vr.Unpack(columnTypes[j])
		}
		values = append(values, v1)

	}
	err = rows.Err()
	if err != nil {
		log.Panic(err)
	}
	return
}