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

// New creates a new instance of DwSqlCommand with the given relation and model.
//
// Parameters:
// - relation: The relation parameter can be either a string representing the table name or a struct of type Relation.
// - model: The model parameter represents the data model to be used with the DwSqlCommand.
//
// Returns:
// - sqlc: A pointer to the newly created DwSqlCommand instance.
// - err: An error if the relation parameter is neither a string nor a struct of type Relation.
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

// parseFieldData parses the field data of a given model using reflection.
//
// It takes a model interface{} as a parameter and returns a map[string]*FieldData,
// a []string, and an error. The map[string]*FieldData contains the field data
// of the model, where the key is the field name and the value is a pointer to a
// FieldData struct. The []string contains the field names of the model. The error
// is returned if there is an error during the parsing process.
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

// Connect sets the connection for DwSqlCommand.
//
// Parameters:
// - conn: The database connection to be set.
func (sqlc *DwSqlCommand) Connect(conn *sql.DB) {
	sqlc.conn = conn
}

// SetTransaction sets the transaction for the DwSqlCommand.
//
// Parameters:
// - tx: The transaction to be set.
func (sqlc *DwSqlCommand) SetTransaction(tx *sql.Tx) {
	sqlc.tx = tx
}

// GetModel retrieves the data model from DwSqlCommand.
//
// No parameters.
// Returns the model interface{}.
func (sqlc *DwSqlCommand) GetModel() (model interface{}) {
	return sqlc.model
}

// affectedFields returns the list of fields affected by the DwSqlCommand.
//
// It takes an optional parameter `fieldnames` which is a variadic string slice.
// If `fieldnames` is empty, it returns all the fields in `sqlc.fieldnames`.
// Otherwise, it returns the specified `fieldnames`.
//
// Parameters:
// - fieldnames: A variadic string slice representing the names of the fields to be affected.
//
// Returns:
// - fields: A string slice representing the affected fields.
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

// GetTablename returns the fully qualified table name if a schema is defined,
// otherwise it returns only the table name.
//
// Parameters:
// - None
//
// Returns:
// - tablename (string): The fully qualified table name or the table name.
func (sqlc *DwSqlCommand) GetTablename() (tablename string) {
	if sqlc.schema != "" {
		return fmt.Sprintf("%s.%s", sqlc.schema, sqlc.tablename)
	} else {
		return sqlc.tablename
	}
}

// CreateInsertQuery generates an SQL insert query based on the provided field names.
//
// Parameters:
// - fieldnames: A variadic string slice representing the names of the fields to be inserted.
//
// Returns:
// - query: A pointer to DwQuery representing the generated insert query.
// - err: An error indicating any issues that occurred during query generation.
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

// CreateUpdateQuery generates an SQL update query based on the provided keys and field names.
//
// Parameters:
// - keys: a slice of strings representing the keys to be used for the update.
// - fieldnames: variadic parameter of strings representing the field names to be updated.
//
// Returns:
// - query: a pointer to a DwQuery struct representing the generated SQL update query.
// - err: an error if any occurred during the generation of the query.
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

// CreateDeleteQuery generates an SQL delete query based on the provided keys and field names.
//
// Parameters:
// - keys: a slice of strings representing the keys to be used for the delete.
// - fieldnames: variadic parameter of strings representing the field names that will be used in the query.
//
// Returns:
// - query: a pointer to a DwQuery struct representing the generated SQL delete query.
// - err: an error if any occurred during the generation of the query.
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

// CreateParameter creates parameters for the SQL query based on the model's fields.
//
// Parameters:
// - query: The DwQuery containing the fields for the SQL query.
// - model: The model interface{} from which the values are extracted.
//
// Returns:
// - params: A slice of any containing the extracted values as parameters.
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

// getAffectedNames returns a slice of strings containing the names of fields in the given model that are not nil.
//
// Parameters:
// - model: An interface{} representing the model from which to extract the field names.
//
// Returns:
// - aff_names: A slice of strings containing the names of non-nil fields in the model.
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

// ExecuteQuery executes a SQL query using the provided DwQuery and parameters.
//
// Parameters:
// - query: A pointer to a DwQuery object representing the SQL query to be executed.
// - params: A variadic parameter of any type representing the parameters to be used in the query.
//
// Returns:
// - res: A sql.Result object representing the result of the query execution.
// - err: An error object if any error occurred during the query execution.
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

// Insert inserts data into the database using the provided model.
//
// Parameters:
// - model: An interface{} representing the data model to be inserted.
//
// Returns:
// - res: A sql.Result object representing the result of the insert operation.
// - err: An error object if any error occurred during the insert operation.
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

// Update updates data in the database using the provided model and keys.
//
// Parameters:
// - model: An interface{} representing the data model to be updated.
// - keys: A string slice representing the keys to identify the records to be updated.
//
// Returns:
// - res: A sql.Result object representing the result of the update operation.
// - err: An error object if any error occurred during the update operation.
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

// Delete deletes data from the database using the provided model.
//
// Parameters:
// - model: An interface{} representing the data model to be deleted.
//
// Returns:
// - res: A sql.Result object representing the result of the delete operation.
// - err: An error object if any error occurred during the delete operation.
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
