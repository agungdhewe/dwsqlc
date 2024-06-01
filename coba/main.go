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

	// TestInsert()
	// TestUpdate()
	TestDelete()

	// TestInsertWithTransaction()
	// TestUpdateWithTransaction()
}

func TestDelete() {
	conn, err := sql.Open("pgx", "host=localhost port=5432 dbname=fgodblocal user=fgta password=rahasia")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	log.Println("connection initialized.")

	// Siapkan data
	sqlc, err := dwsqlc.New("latihan.heinv", &Heinv{})
	if err != nil {
		panic(err.Error())
	}

	sqlc.Connect(conn)
	_, err = sqlc.Delete(&Heinv{
		Id: "TM345",
	})
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Data Deleted")
}

func TestInsert() {

	// Koneksi ke database
	conn, err := sql.Open("pgx", "host=localhost port=5432 dbname=fgodblocal user=fgta password=rahasia")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	log.Println("connection initialized.")

	// Siapkan data
	sqlc, err := dwsqlc.New("latihan.heinv", &Heinv{})
	if err != nil {
		panic(err.Error())
	}

	sqlc.Connect(conn)
	_, err = sqlc.Insert(&Heinv{
		Id:    "TM345",
		Art:   "55643",
		Descr: "ini test insert",
	})
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Data Inserted")
}

func TestUpdate() {
	// Koneksi ke database
	conn, err := sql.Open("pgx", "host=localhost port=5432 dbname=fgodblocal user=fgta password=rahasia")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connection initialized.")
	defer conn.Close()

	// Siapkan data
	sqlc, err := dwsqlc.New("latihan.heinv", &Heinv{})
	if err != nil {
		panic(err.Error())
	}

	sqlc.Connect(conn)
	keys := []string{"Id"}
	_, err = sqlc.Update(&Heinv{
		Id:    "TM345",
		Art:   "55643",
		Descr: "ini yang coba diupdate",
	}, keys)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Data Updated")
}

func TestUpdateWithTransaction() {
	// Koneksi ke database
	conn, err := sql.Open("pgx", "host=localhost port=5432 dbname=fgodblocal user=fgta password=rahasia")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connection initialized.")
	defer conn.Close()

	// Siapkan data
	sqlc, err := dwsqlc.New("latihan.heinv", &Heinv{})
	if err != nil {
		panic(err.Error())
	}

	// Siapkan Transaksi
	ctx := context.Background()
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		panic(err.Error())
	}

	sqlc.Connect(conn)
	sqlc.SetTransaction(tx)
	keys := []string{"Id"}
	_, err = sqlc.Update(&Heinv{
		Id:    "TM345",
		Art:   "55643",
		Descr: "ini yang coba diupdate",
	}, keys)
	if err != nil {
		panic(err.Error())
	}
	tx.Commit()
	fmt.Println("Data Updated")
}

func TestInsertWithTransaction() {

	// Koneksi ke database
	conn, err := sql.Open("pgx", "host=localhost port=5432 dbname=fgodblocal user=fgta password=rahasia")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	log.Println("connection initialized.")

	// Siapkan data
	sqlc, err := dwsqlc.New("latihan.heinv", &Heinv{})
	if err != nil {
		panic(err.Error())
	}

	// Siapkan Transaksi
	ctx := context.Background()
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		panic(err.Error())
	}

	sqlc.Connect(conn)
	sqlc.SetTransaction(tx)
	_, err = sqlc.Insert(&Heinv{
		Id:    "TM345",
		Art:   "55643",
		Descr: "ini test insert",
	})
	if err != nil {
		panic(err.Error())
	}
	tx.Commit()
	fmt.Println("Data Inserted")
}
