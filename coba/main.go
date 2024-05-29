package main

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"slices"
	"time"

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

type Data struct {
	Id         string  `field:"id"`
	Nama       string  `field:"nama"`
	IsDisabled bool    `field:"isdisabled"`
	Tanggal    string  `field:"tanggal"`
	Amount     float32 `field:"amount"`
	Jam        string  `field:"jam"`
	Timestamp  string  `field:"dt"`
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
		Id:         "satu",
		Nama:       "agung nugroho",
		Tanggal:    "2014-10-11",
		Amount:     5600000,
		Jam:        "11:12",
		Timestamp:  time.Now().Format("2006-01-02 15:04:05"),
		IsDisabled: false,
	}
	keys := []string{"Id"}

	fmt.Println(data)

	// coba extract dengan reflection
	val := reflect.ValueOf(data).Elem()
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		valueField := val.Field(i)

		name := typeField.Name
		field_index := 1 + typeField.Index[0]
		field_name := typeField.Tag.Get("field")
		field_value := valueField.Interface()

		var iskey string
		if slices.Contains(keys, name) {
			iskey = "KEY"
		} else {
			iskey = ""
		}

		fmt.Println(field_index, name, field_name, field_value, iskey)
	}
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
