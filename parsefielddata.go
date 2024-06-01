package dwsqlc

import (
	"reflect"
)

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
