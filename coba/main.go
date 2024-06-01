package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/agungdhewe/dwsqlc"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Heinv struct {
	Id    dwsqlc.DbTypeVarchar `field:"heinv_id" default:""`
	Art   dwsqlc.DbTypeVarchar `field:"art" default:""`
	Mat   dwsqlc.DbTypeVarchar `field:"mat" default:""`
	Name  dwsqlc.DbTypeVarchar `field:"nama" default:""`
	Descr dwsqlc.DbTypeVarchar `field:"descr" default:""`
}

func main() {
	// connect ke database
	conn, err := sql.Open("pgx", "host=localhost port=5432 dbname=fgodblocal user=fgta password=rahasia")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	log.Println("database connected.")

	// TestInsert(conn)
	//TestUpdate(conn)

	// siapkan sqlcommand
	sqlc, err := dwsqlc.New("latihan.heinv", &Heinv{})
	if err != nil {
		panic(err.Error())
	}

	ctx := context.Background()
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(tx)

	// siapkan query insert
	query, err := sqlc.CreateInsertQuery("Id", "Art", "Mat")
	if err != nil {
		panic(err.Error())
	}

	// siapkan SQL Statement dari query yang sudah dibuat
	var stmt *sql.Stmt

	stmt, err = conn.Prepare(query.Sql())
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	model := sqlc.GetModel().(*Heinv)
	model.Id = "19"
	model.Art = "345345"
	model.Mat = "44" //sql.NullBool{}
	params := sqlc.CreateParameter(query, model)

	_, err = stmt.Exec(params...)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("data inserted")
}

func TestUpdate(conn *sql.DB) {
	// siapkan table yang akan diupdate
	rel := dwsqlc.Relation{
		Table:  "heinv",
		Schema: "latihan",
	}

	// siapkan sqlcommand
	sqlc, err := dwsqlc.New(rel, &Heinv{})
	if err != nil {
		panic(err.Error())
	}

	// siapkan query update
	keys := []string{"Id"}
	query, err := sqlc.CreateUpdateQuery(keys, "Id", "Art", "Mat")
	if err != nil {
		panic(err.Error())
	}

	// siapkan SQL Statement dari query yang sudah dibuat
	var stmt *sql.Stmt
	stmt, err = conn.Prepare(query.Sql())
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	model := sqlc.GetModel().(*Heinv)
	model.Id = "17"
	model.Art = "diupdate"
	//model.Mat = "44" //sql.NullBool{}
	params := sqlc.CreateParameter(query, model)

	_, err = stmt.Exec(params...)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("data updated")
}

func TestInsert(conn *sql.DB) {

	// siapkan table yang akan diinsert
	rel := dwsqlc.Relation{
		Table:  "heinv",
		Schema: "latihan",
	}

	// siapkan sqlcommand
	sqlc, err := dwsqlc.New(rel, &Heinv{})
	if err != nil {
		panic(err.Error())
	}

	// siapkan query insert
	query, err := sqlc.CreateInsertQuery("Id", "Art", "Mat")
	if err != nil {
		panic(err.Error())
	}

	// siapkan SQL Statement dari query yang sudah dibuat
	var stmt *sql.Stmt
	stmt, err = conn.Prepare(query.Sql())
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	model := sqlc.GetModel().(*Heinv)
	model.Id = "17"
	model.Art = "345345"
	model.Mat = "44" //sql.NullBool{}
	params := sqlc.CreateParameter(query, model)

	_, err = stmt.Exec(params...)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("data inserted")

}

// func main() {

// 	rel := dwsqlc.Relation{
// 		Table:  "heinv",
// 		Schema: "latihan",
// 	}

// 	sqlc, err := dwsqlc.New(rel, &Heinv{})
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	query, err := sqlc.CreateInsertQuery("Id", "Art", "Mat")
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	fmt.Println(query.Sql())

// 	// untuk update
// 	/*
// 		keys := []string{"Id"}
// 		query, err := sqlc.CreateUpdateQuery(keys, "Id", "Art", "Mat")
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		fmt.Println(query.Sql())
// 	*/
// 	// var model *Heinv

// 	// // loop 1
// 	model := sqlc.GetModel().(*Heinv)
// 	model.Id = 123
// 	model.Art = "234"
// 	model.Mat = sql.NullBool{}
// 	params := sqlc.CreateParameter(query, model)
// 	fmt.Println(params...)

// 	// // loop 2
// 	// model = sqlc.GetModel().(*Heinv)
// 	// model.Id = "CC"
// 	// sqlc.CreateParameter(query, model)
// }
