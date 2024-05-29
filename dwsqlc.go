package dwsqlc

import (
	"fmt"
	"log"
	"reflect"
)

const (
	DbValueNull  = "_DB_VALUE_NULL_"
	DbValueTrue  = "_DB_VALUE_TRUE_"
	DbValueFalse = "_DB_VALUE_FALSE_"
	DbValueNow   = "_DB_VALUE_NOW_"
)

type DbTypeVarchar string
type DbTypeInteger string
type DbTypeBool string
type DbTypeDate string
type DbTypeDatetime string
type DbTypeTime string
type DbTypeDecimal string

type DwSqlCommand struct {
	schema     string
	tablename  string
	modeltype  reflect.Type
	fielddata  map[string]*FieldData
	fieldnames []string
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

func New(relation interface{}, modeltype reflect.Type) (sqlc *DwSqlCommand, err error) {
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

	// ambil seluruh informasi field yang ada di modeltype
	var fielddata map[string]*FieldData
	var fieldnames []string
	fielddata, fieldnames, err = parseFieldData(modeltype)
	if err != nil {
		return nil, err
	}

	// siapkan hasil berupa DwSqlCommand
	sqlc = &DwSqlCommand{
		schema:     schema,
		tablename:  tablename,
		modeltype:  modeltype,
		fielddata:  fielddata,
		fieldnames: fieldnames,
	}

	return sqlc, nil
}

func (sqlc *DwSqlCommand) GetModel() (model interface{}) {
	value := reflect.New(sqlc.modeltype)
	model = value.Interface()
	return model
}

func (sqlc *DwSqlCommand) CreateInsertQuery(fieldnames ...string) (query *DwQuery) {
	log.Println("buat querynya")

	if len(fieldnames) == 0 {
		// semua field
		log.Println("select *")
	} else {
		// field yang dipilih saja
		log.Println("select", fieldnames)

	}

	// log.Println(len(fieldname))

	query = &DwQuery{
		sql:    "select bla bla",
		fields: fieldnames,
	}
	return query
}

func (sqlc *DwSqlCommand) CreateUpdateQuery(keys []string, fieldnames ...string) (query *DwQuery) {
	log.Println(len(fieldnames))
	query = &DwQuery{
		sql:    "update blu blu blu",
		fields: fieldnames,
	}
	return query
}

func (sqlc *DwSqlCommand) CreateParameter(query *DwQuery, model interface{}) (params []string) {
	return params
}
