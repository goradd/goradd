package db

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	. "github.com/goradd/goradd/pkg/orm/query"
	"io/ioutil"
	"net/http"
	"path"
)

//
type Rest struct {
	url         string
	description *Database
}

// NewRest returns a new Rest goradd database object that you can add to the datastore.
func NewRest(dbKey string, url string) *Rest {

	d := Rest{
		url: url,
	}

	return &d
}

// NewBuilder returns a new query builder to build a query that will be processed by the database.
func (r *Rest) NewBuilder() QueryBuilderI {
	return NewRestBuilder(r)
}

// Describe returns the database description object. Rest databases are not describable at this point.
// Maybe someday.
func (r *Rest) Describe() *Database {
	return nil
}

// Update sets a record that already exists in the database to the given data, updating only the fields given.
func (r *Rest) Update(ctx context.Context, table string, fields map[string]interface{}, pkName string, pkValue string) {
	var response *http.Response
	var body = make(map[string]interface{})

	body["t"] = table
	body["f"] = fields
	body["n"] = pkName
	body["v"] = pkValue

	j, err := json.Marshal(body)
	response, err = http.Post(path.Join(r.url, "upd"), "application/json", bytes.NewBuffer(j))
	if err == nil {
		data, _ := ioutil.ReadAll(response.Body)
		if data != nil && len(data) != 0 {
			err = fmt.Errorf(string(data))
		}
	}
	return

}

// Insert inserts the given data as a new record in the database.
func (r *Rest) Insert(ctx context.Context, table string, fields map[string]interface{}) (id string, err error) {
	var response *http.Response
	var body = make(map[string]interface{})

	body["t"] = table
	body["f"] = fields

	j, err := json.Marshal(body)
	response, err = http.Post(path.Join(r.url, "ins"), "application/json", bytes.NewBuffer(j))
	if err == nil {
		data, _ := ioutil.ReadAll(response.Body)
		var result map[string]string
		err = json.Unmarshal(data, &result)
		if result == nil {
			if v, ok := result["err"]; ok {
				err = fmt.Errorf(v)
			} else if v, ok := result["id"]; ok {
				id = v
			} else {
				err = fmt.Errorf("missing return id")
			}
		} else {
			err = fmt.Errorf("missing return result")
		}
	}
	return
}

// Delete deletes the indicated record from the database.
func (r *Rest) Delete(ctx context.Context, table string, pkName string, pkValue string) (err error) {
	var response *http.Response
	var body = make(map[string]string)

	body["t"] = table
	body["n"] = pkName // Not sure this is needed
	body["v"] = pkValue

	j, err := json.Marshal(body)
	response, err = http.Post(path.Join(r.url, "del"), "application/json", bytes.NewBuffer(j))
	if err == nil {
		data, _ := ioutil.ReadAll(response.Body)
		if data != nil && len(data) != 0 {
			err = fmt.Errorf(string(data))
		}
	}
	return
}

func (r *Rest) Get(ctx context.Context, jsonRequest []byte) (result []map[string]interface{}, err error) {
	var response *http.Response
	response, err = http.Post(path.Join(r.url, "get"), "application/json", bytes.NewBuffer(jsonRequest))
	if err == nil {
		data, _ := ioutil.ReadAll(response.Body)
		err = json.Unmarshal(data, result)
	}
	return
}

func (r *Rest) Begin(ctx context.Context) (txid TransactionID) {
	// TODO: Implement transactions in coordination with the rest server, where the rest server holds
	// on to the transaction id, and rolls back after a set timeout.
	return TransactionID(0)
}

// Commit commits the transaction, and if an error occurs, will panic with the error.
func (r *Rest) Commit(ctx context.Context, txid TransactionID) {

}

// Rollback will rollback the transaction if the transaction is still pointing to the given txid. This gives the effect
// that if you call Rollback on a transaction that has already been committed, no Rollback will happen. This makes it easier
// to implement a transaction management scheme, because you simply always defer a Rollback after a Begin. Pass the txid
// that you got from the Begin to the Rollback
func (r *Rest) Rollback(ctx context.Context, txid TransactionID) {

}
