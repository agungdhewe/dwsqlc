package dwsqlc

import (
	"fmt"
	"reflect"
	"strings"
)

// membuat query string insert
func (sqlc *DwSqlCommand) CreateInsertQuery(fieldnames ...string) (query *DwQuery, err error) {

	// ambil field yang akan diinsert
	fields := sqlc.affectedFields(fieldnames...)

	// loop data insert fields
	var i int
	var name string

	c := len(fields)
	ins_fieldnames := make([]string, c)
	ins_fieldvalues := make([]string, c)
	for i, name = range fields {
		fielddata, inmap := sqlc.fielddata[name]
		if !inmap {
			return nil, fmt.Errorf("%s tidak ada di struktur %s", name, reflect.TypeOf(sqlc.model).Name())
		}
		ins_fieldnames[i] = fielddata.FieldName
		ins_fieldvalues[i] = fmt.Sprintf("$%d", i+1)
	}

	tablename := sqlc.GetTablename()
	f := strings.Join(ins_fieldnames, ", ")
	v := strings.Join(ins_fieldvalues, ", ")
	sqltext := fmt.Sprintf("insert into %s\r\n(%s)\r\nvalues\r\n(%s)", tablename, f, v)

	query = &DwQuery{
		sql:    sqltext,
		fields: fields,
	}
	return query, nil
}

func (sqlc *DwSqlCommand) Insert(model interface{}) (res any, err error) {
	if sqlc.conn == nil {
		return nil, fmt.Errorf("connection is not set, you have to call Connect(*sql.Db) first")
	}

	query, err := sqlc.CreateInsertQuery("Id", "Art", "Mat")
	if err != nil {

	}
	params := sqlc.CreateParameter(query, model)

	return nil, nil
}
