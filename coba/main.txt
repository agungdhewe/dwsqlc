package main

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"slices"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	// data := &Data{}
	// keys := []string{"id"}
	// dwsqlc.New("toko", data, keys)

	// cmd := dwsqlc.CreateSqlInsert()
	// sql := cmd.GetSQL()
	// param := cmd.GetParam()

	// cmd.Ready()

	if false {
		cobaQueryDasar()
	}

	cobaComposeDwSqlC()
}

const DbNull = "DB_VALUE_NULL"
const DbTrue = "DB_VALUE_TRUE"
const DbFalse = "DB_VALUE_FALSE"
const DbNow = "DB_VALUE_NOW"

type SqlVarchar string
type SqlBool string
type SqlDate string
type SqlTime string
type SqlDateTime string
type SqlDecimal string

type Data struct {
	Id         SqlVarchar  `field:"id"`
	Nama       SqlVarchar  `field:"nama"`
	IsDisabled SqlBool     `field:"isdisabled"`
	Tanggal    SqlDate     `field:"tanggal"`
	Amount     SqlDecimal  `field:"amount"`
	Jam        SqlTime     `field:"jam"`
	Timestamp  SqlDateTime `field:"dt"`
}

type FieldInfo struct {
	Index int
	Name  string
	Value any
	IsKey bool
	Type  string
}

func cobaComposeDwSqlC() {
	// koneksi database
	conn, err := sql.Open("pgx", "host=localhost port=5432 dbname=tfidblocal user=tfi password=rahasia")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	log.Println("database connected.")

	data := &Data{
		Id:   "123",
		Nama: "agung nugroho",
		// Tanggal: "2014-10-11",
		// Amount:     5600000,
		// Jam:        "11:12",
		// Timestamp:  time.Now().Format("2006-01-02 15:04:05"),
		IsDisabled: DbTrue,
	}

	keys := []string{"Id"}

	// coba extract dengan reflection
	val := reflect.ValueOf(data).Elem()

	// loop data di struct
	comp := map[string]*FieldInfo{}
	n := val.NumField()
	for i := 0; i < n; i++ {
		typeField := val.Type().Field(i)
		valueField := val.Field(i)

		name := typeField.Name
		field_index := 1 + typeField.Index[0]
		field_name := typeField.Tag.Get("field")
		field_value := fmt.Sprintf("%v", valueField.Interface())
		field_type := typeField.Type.Name()

		var iskey bool
		if slices.Contains(keys, name) {
			iskey = true
		} else {
			iskey = false
		}

		if field_value == "" {
			continue
		}

		comp[field_name] = &FieldInfo{
			Index: field_index,
			Name:  field_name,
			Value: field_value,
			IsKey: iskey,
			Type:  field_type,
		}
	}

	// siapkan query
	n = len(comp)

	fmt.Println("===================")
	fmt.Println(n)

}

func cobaQueryDasar() {
	var err error

	fmt.Println("Query Select")

	// koneksi database
	conn, err := sql.Open("pgx", "host=localhost port=5432 dbname=tfidblocal user=tfi password=rahasia")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	log.Println("database connected.")

	// test insert data
	query := `
		insert into temp.tmp_heinv 
		(heinv_id, heinv_art, heinv_col, heinv_mat, heinv_name)
		values
		($1, $2, $3, $4, $5)
	`
	var stmt *sql.Stmt
	stmt, err = conn.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// single insert
	args := []any{"id2", "art", "col", "mat", "name"}
	_, err = stmt.Exec(args...)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("data inserted")

	/*
		// coba multiple insert
		data := [][]any{
			[]any{"5", "art1", "col1", "mat1", "name1"},
			[]any{"6", "art2", "col2", "mat2", "name2"},
			[]any{"7", "art3", "col3", "mat3", "name3"},
			[]any{"8", "art4", "col4", "mat4", "name4"},
		}

		for _, rowargs := range data {
			_, err = stmt.Exec(rowargs...)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("data inserted")
		}
	*/

}
