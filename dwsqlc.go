package dwsqlc

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
)

type DbTypeVarchar any
type DbTypeInteger any
type DbTypeBool any
type DbTypeDate any
type DbTypeDatetime any
type DbTypeTime any
type DbTypeDecimal any

type DwSqlCommand struct {
	schema     string
	tablename  string
	model      interface{}
	fielddata  map[string]*FieldData
	fieldnames []string
	conn       *sql.DB
	tx         *sql.Tx
	ctx        context.Context
}

type Relation struct {
	Table  string
	Schema string
}

type FieldData struct {
	Name          string
	FieldName     string
	FieldValue    any
	FieldTypeName string
	DefaultValue  string
	Index         int
}

func New(relation interface{}, model interface{}) (sqlc *DwSqlCommand, err error) {
	var tablename string
	var schema string

	// cek ripe data dari relasi, harus berupa string, berarti langsung nama tabel
	// atau struct Relation, apabila akan menggunakan prefix beruma nama skema
	typename := reflect.TypeOf(relation).String()
	if typename == "string" {
		tablename = relation.(string)
		schema = ""
	} else if typename == "dwsqlc.Relation" {
		tablename = relation.(Relation).Table
		schema = relation.(Relation).Schema
	} else {
		return nil, fmt.Errorf("relation parameter value must be string, or struct Relation")

	}

	// ambil seluruh informasi field yang ada di model
	var fielddata map[string]*FieldData
	var fieldnames []string
	fielddata, fieldnames, err = parseFieldData(model)
	if err != nil {
		return nil, err
	}

	// siapkan hasil berupa DwSqlCommand
	sqlc = &DwSqlCommand{
		schema:     schema,
		tablename:  tablename,
		model:      model,
		fielddata:  fielddata,
		fieldnames: fieldnames,
		tx:         nil,
		conn:       nil,
	}

	return sqlc, nil
}

func (sqlc *DwSqlCommand) Connect(conn *sql.DB) {
	sqlc.conn = conn
}

func (sqlc *DwSqlCommand) BeginTransaction() (tx *sql.Tx, err error) {
	if sqlc.conn == nil {
		return nil, fmt.Errorf("connection is not set, you have to call Connect(*sql.Db) first")
	}

	sqlc.ctx = context.Background()
	tx, err = sqlc.conn.BeginTx(sqlc.ctx, nil)
	if err != nil {
		return nil, err
	}

	sqlc.tx = tx
	return tx, nil
}

func (sqlc *DwSqlCommand) Commit() {
	if sqlc.tx == nil {
		return
	}
	sqlc.tx.Commit()
}

func (sqlc *DwSqlCommand) Rollback() {
	if sqlc.tx == nil {
		return
	}
	sqlc.tx.Rollback()
}

func (sqlc *DwSqlCommand) GetModel() (model interface{}) {
	return sqlc.model
}

func (sqlc *DwSqlCommand) affectedFields(fieldnames ...string) (fields []string) {
	if len(fieldnames) == 0 {
		// jika fieldnames tidak diisi, berarti insert untuk semua field
		fields = sqlc.fieldnames
	} else {
		// jika fieldname diisi, berarti yang diinsert hanya field yang dipilih saja
		fields = fieldnames
	}
	return fields
}

func (sqlc *DwSqlCommand) GetTablename() (tablename string) {
	if sqlc.schema != "" {
		return fmt.Sprintf("%s.%s", sqlc.schema, sqlc.tablename)
	} else {
		return sqlc.tablename
	}
}
