package sql

import (
	"database/sql"
	"github.com/goradd/goradd/pkg/orm/query"
	"log"
)

type sqlCursor struct {
	rows *sql.Rows
	columnTypes []query.GoColumnType
	columnNames []string
	builder *Builder
	columnReceivers []SqlReceiver
	columnValueReceivers []interface{}
}

func NewSqlCursor(rows *sql.Rows,
	columnTypes []query.GoColumnType,
	columnNames []string,
	builder *Builder,
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
