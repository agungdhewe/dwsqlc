package dwsqlc

import "reflect"

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
