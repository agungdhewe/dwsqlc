package dwsqlc

import (
	"database/sql"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

type DbTypeVarchar any
type DbTypeInteger any
type DbTypeBool any
type DbTypeDate any
type DbTypeDatetime any
type DbTypeTime any
type DbTypeDecimal any

type DbFieldNotSet any

type DwSqlCommand struct {
	schema     string
	tablename  string
	model      interface{}
	fielddata  map[string]*FieldData
	fieldnames []string
	conn       *sql.DB
	tx         *sql.Tx
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

	// notset := nil.(DbFieldNotSet)

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

func parseFieldData(model interface{}) (fielddata map[string]*FieldData, fieldnames []string, err error) {
	fielddata = map[string]*FieldData{}

	// coba extract dengan reflection
	val := reflect.ValueOf(model).Elem()
	n := val.NumField()
	fieldnames = make([]string, n)
	for i := 0; i < n; i++ {
		typeField := val.Type().Field(i)
		// valueField := val.Field(i)

		name := typeField.Name
		index := 1 + typeField.Index[0]
		field_name := typeField.Tag.Get("field")
		field_typename := typeField.Type.Name()
		defaultvalue := typeField.Tag.Get("default")

		f := &FieldData{
			Index:         index,
			Name:          name,
			FieldName:     field_name,
			FieldTypeName: field_typename,
			DefaultValue:  defaultvalue,
		}

		fielddata[name] = f
		fieldnames[i] = name
	}

	return fielddata, fieldnames, err
}

func (sqlc *DwSqlCommand) Connect(conn *sql.DB) {
	sqlc.conn = conn
}

func (sqlc *DwSqlCommand) SetTransaction(tx *sql.Tx) {
	sqlc.tx = tx
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

func (sqlc *DwSqlCommand) CreateUpdateQuery(keys []string, fieldnames ...string) (query *DwQuery, err error) {

	// ambil field yang akan dipakai sebagai key
	upd_keynames := make([]string, len(keys))
	for i, keyname := range keys {
		_, inmap := sqlc.fielddata[keyname]
		if !inmap {
			return nil, fmt.Errorf("key %s tidak ada di struktur %s", keyname, reflect.TypeOf(sqlc.model).Name())
		}
		upd_keynames[i] = keyname
	}

	// ambil field yang akan diupdate
	aff_fields := append(sqlc.affectedFields(fieldnames...), upd_keynames...)

	// loop data yang akan diupdate
	for _, name := range aff_fields {
		_, inmap := sqlc.fielddata[name]
		if !inmap {
			return nil, fmt.Errorf("field %s tidak ada di struktur %s", name, reflect.TypeOf(sqlc.model).Name())
		}
	}
	slices.Sort(aff_fields)

	fields := slices.Compact(aff_fields)

	c := len(fields)
	tmp_pairs := make([]string, c)
	tmp_keys := make([]string, c)

	p := 0
	k := 0
	for i, name := range fields {
		fielddata := sqlc.fielddata[name]
		if slices.Contains(upd_keynames, name) {
			tmp_keys[k] = fmt.Sprintf("%s=$%d", fielddata.FieldName, i+1)
			k++
		} else {
			tmp_pairs[p] = fmt.Sprintf("%s=$%d", fielddata.FieldName, i+1)
			p++
		}

	}

	upd_pairs := make([]string, p)
	for i := 0; i < p; i++ {
		upd_pairs[i] = tmp_pairs[i]
	}

	upd_keys := make([]string, k)
	for i := 0; i < k; i++ {
		upd_keys[i] = tmp_keys[i]
	}

	tablename := sqlc.GetTablename()
	fp := strings.Join(upd_pairs, ",\r\n")
	fk := strings.Join(upd_keys, " AND ")
	sqltext := fmt.Sprintf("update %s\r\nset\r\n%s\r\nwhere\r\n%s", tablename, fp, fk)

	query = &DwQuery{
		sql:    sqltext,
		fields: fields,
	}
	return query, nil
}

func (sqlc *DwSqlCommand) CreateDeleteQuery(keys []string, fieldnames ...string) (query *DwQuery, err error) {

	// ambil field yang akan dipakai sebagai key
	del_keynames := make([]string, len(keys))
	for i, keyname := range keys {
		_, inmap := sqlc.fielddata[keyname]
		if !inmap {
			return nil, fmt.Errorf("key %s tidak ada di struktur %s", keyname, reflect.TypeOf(sqlc.model).Name())
		}
		del_keynames[i] = keyname
	}

	// ambil field yang akan diupdate
	aff_fields := append(sqlc.affectedFields(fieldnames...), del_keynames...)

	// loop data yang akan diupdate
	for _, name := range aff_fields {
		_, inmap := sqlc.fielddata[name]
		if !inmap {
			return nil, fmt.Errorf("field %s tidak ada di struktur %s", name, reflect.TypeOf(sqlc.model).Name())
		}
	}
	slices.Sort(aff_fields)

	fields := slices.Compact(aff_fields)

	c := len(fields)
	tmp_keys := make([]string, c)

	k := 0
	for i, name := range fields {
		fielddata := sqlc.fielddata[name]
		if slices.Contains(del_keynames, name) {
			tmp_keys[k] = fmt.Sprintf("%s=$%d", fielddata.FieldName, i+1)
			k++
		}
	}

	del_keys := make([]string, k)
	for i := 0; i < k; i++ {
		del_keys[i] = tmp_keys[i]
	}

	tablename := sqlc.GetTablename()
	fk := strings.Join(del_keys, " AND ")
	sqltext := fmt.Sprintf("delete from %s where %s", tablename, fk)

	query = &DwQuery{
		sql:    sqltext,
		fields: fields,
	}
	return query, nil
}

func (sqlc *DwSqlCommand) CreateParameter(query *DwQuery, model interface{}) (params []any) {
	n := len(query.fields)
	params = make([]any, n)
	val := reflect.ValueOf(model).Elem()
	for i, name := range query.fields {
		value := val.FieldByName(name).Interface()
		params[i] = value
	}
	return params
}

func getAffectedNames(model interface{}) (aff_names []string) {
	val := reflect.ValueOf(model).Elem()
	n := val.NumField()
	tmp_names := make([]string, n)
	u := 0
	for i := 0; i < n; i++ {
		typeField := val.Type().Field(i)
		valueField := val.Field(i)
		currValue := valueField.Interface()
		if currValue != nil {
			tmp_names[i] = typeField.Name
			u++
		}
	}

	// ambil nama field yang akan diinsert ke slice
	// jumlah field yang akan diambil adalah nilai u terakhir
	aff_names = make([]string, u)
	j := 0
	for i := 0; i < n; i++ {
		name := tmp_names[i]
		if name != "" {
			aff_names[j] = name
			j++
		}

	}

	return aff_names
}

func (sqlc *DwSqlCommand) ExecuteQuery(query *DwQuery, params []any) (res sql.Result, err error) {
	var stmt *sql.Stmt
	if sqlc.tx != nil {
		stmt, err = sqlc.tx.Prepare(query.Sql())
		if err != nil {
			return nil, err
		}
	} else {
		stmt, err = sqlc.conn.Prepare(query.Sql())
		if err != nil {
			return nil, err
		}
	}
	defer stmt.Close()

	res, err = stmt.Exec(params...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (sqlc *DwSqlCommand) Insert(model interface{}) (res sql.Result, err error) {
	// tandai nama-nama field yang akan diinsert
	aff_names := getAffectedNames(model)

	// siapkan query
	query, err := sqlc.CreateInsertQuery(aff_names...)
	if err != nil {
		return nil, err
	}
	// siapkan parameter
	params := sqlc.CreateParameter(query, model)

	// eksekusi
	res, err = sqlc.ExecuteQuery(query, params)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (sqlc *DwSqlCommand) Update(model interface{}, keys []string) (res sql.Result, err error) {
	// tandai nama-nama field yang akan diinsert
	aff_names := getAffectedNames(model)

	// siapkan query
	query, err := sqlc.CreateUpdateQuery(keys, aff_names...)
	if err != nil {
		return nil, err
	}
	// siapkan parameter
	params := sqlc.CreateParameter(query, model)

	// eksekusi
	res, err = sqlc.ExecuteQuery(query, params)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (sqlc *DwSqlCommand) Delete(model interface{}) (res sql.Result, err error) {
	// tandai nama-nama field yang akan diinsert
	aff_names := getAffectedNames(model)

	// siapkan query
	query, err := sqlc.CreateDeleteQuery(aff_names, aff_names...)
	if err != nil {
		return nil, err
	}
	// siapkan parameter
	params := sqlc.CreateParameter(query, model)

	// eksekusi
	res, err = sqlc.ExecuteQuery(query, params)
	if err != nil {
		return nil, err
	}

	return res, nil
}
