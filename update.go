package dwsqlc

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

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

func (sqlc *DwSqlCommand) Update(model interface{}, keys []string) (res any, err error) {
	if sqlc.conn == nil {
		return nil, fmt.Errorf("connection is not set, you have to call Connect(*sql.Db) first")
	}

	return nil, nil
}
